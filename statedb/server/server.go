// Copyright 2020 dfuse Platform Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package server

import (
	crand "crypto/rand"

	"context"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/rand"
	"net"
	"net/http"
	"os"
	"strconv"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/teris-io/shortid"

	stackdriverPropagation "contrib.go.opencensus.io/exporter/stackdriver/propagation"
	"github.com/francoispqt/gojay"
	"github.com/gorilla/mux"
	"github.com/streamingfast/derr"
	"github.com/streamingfast/dtracing"
	"github.com/streamingfast/fluxdb"
	"github.com/streamingfast/logging"
	"go.opencensus.io/plugin/ochttp"
	"go.opencensus.io/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var parallelReadRequestCount = 64

var defaultFormat stackdriverPropagation.HTTPFormat = stackdriverPropagation.HTTPFormat{}
var shortIDGenerator *shortid.Shortid
var traceIDGenerator *defaultIDGenerator

type defaultIDGenerator struct {
	sync.Mutex

	// Please keep these as the first fields
	// so that these 8 byte fields will be aligned on addresses
	// divisible by 8, on both 32-bit and 64-bit machines when
	// performing atomic increments and accesses.
	// See:
	// * https://github.com/census-instrumentation/opencensus-go/issues/587
	// * https://github.com/census-instrumentation/opencensus-go/issues/865
	// * https://golang.org/pkg/sync/atomic/#pkg-note-BUG
	nextSpanID uint64
	spanIDInc  uint64

	traceIDAdd  [2]uint64
	traceIDRand *rand.Rand
}

// NewSpanID returns a non-zero span ID from a randomly-chosen sequence.
func (gen *defaultIDGenerator) NewSpanID() [8]byte {
	var id uint64
	for id == 0 {
		id = atomic.AddUint64(&gen.nextSpanID, gen.spanIDInc)
	}
	var sid [8]byte
	binary.LittleEndian.PutUint64(sid[:], id)
	return sid
}

// NewTraceID returns a non-zero trace ID from a randomly-chosen sequence.
// mu should be held while this function is called.
func (gen *defaultIDGenerator) NewTraceID() trace.TraceID {
	var tid [16]byte
	// Construct the trace ID from two outputs of traceIDRand, with a constant
	// added to each half for additional entropy.
	gen.Lock()
	binary.LittleEndian.PutUint64(tid[0:8], gen.traceIDRand.Uint64()+gen.traceIDAdd[0])
	binary.LittleEndian.PutUint64(tid[8:16], gen.traceIDRand.Uint64()+gen.traceIDAdd[1])
	gen.Unlock()
	return tid
}

func init() {
	// A new generator using the default alphabet set
	shortIDGenerator = shortid.MustNew(1, shortid.DefaultABC, uint64(time.Now().UnixNano()))

	traceIDGenerator = &defaultIDGenerator{}
	// initialize traceID and spanID generators.
	var rngSeed int64
	for _, p := range []interface{}{
		&rngSeed, &traceIDGenerator.traceIDAdd, &traceIDGenerator.nextSpanID, &traceIDGenerator.spanIDInc,
	} {
		binary.Read(crand.Reader, binary.LittleEndian, p)
	}
	traceIDGenerator.traceIDRand = rand.New(rand.NewSource(rngSeed))
	traceIDGenerator.spanIDInc |= 1

}

type EOSServer struct {
	httpServer *http.Server
	db         *fluxdb.FluxDB
	addr       string
	mux        *mux.Router
}

func New(addr string, db *fluxdb.FluxDB) *EOSServer {
	router := mux.NewRouter()
	srv := &EOSServer{
		addr: addr,
		mux:  router,
		db:   db,
	}

	metricsRouter := router.PathPrefix("/").Subrouter()
	coreRouter := router.PathPrefix("/").Subrouter()

	// Metrics & health endpoints
	metricsRouter.HandleFunc("/ping", srv.pingHandler)
	metricsRouter.HandleFunc("/healthz", srv.healthzHandler)

	// Core endpoints
	coreRouter.Use(openCensusMiddleware)
	coreRouter.Use(LoggingMiddleware)
	coreRouter.Use(trackingMiddleware)

	coreRouter.Methods("GET").Path("/v0/state/abi").HandlerFunc(srv.getABIHandler)
	coreRouter.Methods("POST").Path("/v0/state/abi/bin_to_json").HandlerFunc(srv.decodeABIHandler)
	coreRouter.Methods("GET", "POST").Path("/v0/state/key_accounts").HandlerFunc(srv.listKeyAccountsHandler)
	coreRouter.Methods("GET").Path("/v0/state/permission_links").HandlerFunc(srv.listLinkedPermissionsHandler)
	coreRouter.Methods("GET").Path("/v0/state/table").HandlerFunc(srv.listTableRowsHandler)

	coreRouter.Methods("GET").Path("/v0/state/table/row").HandlerFunc(srv.getTableRowHandler)
	coreRouter.Methods("GET").Path("/v0/state/table_scopes").HandlerFunc(srv.listTableScopesHandler)
	coreRouter.Methods("GET", "POST").Path("/v0/state/tables/accounts").HandlerFunc(srv.listTablesRowsForAccountsHandler)
	coreRouter.Methods("GET", "POST").Path("/v0/state/tables/scopes").HandlerFunc(srv.listTablesRowsForScopesHandler)

	db.OnTerminating(func(_ error) {
		zlog.Info("gracefully shutting down http server, draining connections")
		if srv.httpServer != nil {
			zlog.Info("allowing fluxdb server to gracefully shuts down without interrupting any active connections")
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			srv.httpServer.Shutdown(ctx)
		}
	})

	return srv
}

func (srv *EOSServer) Handler() http.Handler {
	return srv.mux
}

func (srv *EOSServer) Serve() {
	zlog.Info("listening & serving HTTP content", zap.String("http_listen_addr", srv.addr))
	errorLogger, err := zap.NewStdLogAt(zlog, zap.ErrorLevel)
	if err != nil {
		srv.db.Shutdown(fmt.Errorf("unable to create error logger: %w", err))
		return
	}

	srv.httpServer = &http.Server{
		Addr:     srv.addr,
		Handler:  srv.Handler(),
		ErrorLog: errorLogger,
	}

	err = srv.httpServer.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		srv.db.Shutdown(fmt.Errorf("failed listening http %q: %w", srv.addr, err))
	}
}

func (srv *EOSServer) pingHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("pong\n"))
}

func (srv *EOSServer) healthzHandler(w http.ResponseWriter, r *http.Request) {
	if !derr.IsShuttingDown() && srv.db.IsReady() {
		w.Write([]byte("ready\n"))
	} else {
		http.Error(w, "not ready\n", http.StatusServiceUnavailable)
	}
}

func openCensusMiddleware(next http.Handler) http.Handler {
	return &ochttp.Handler{
		Handler:     next,
		Propagation: &stackdriverPropagation.HTTPFormat{},
	}
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return NewAddTraceIDAwareLoggerMiddleware(next, zlog, stackdriverPropagation.HTTPFormat{})

}

func NewAddTraceIDAwareLoggerMiddleware(next http.Handler, rootLogger *zap.Logger, propagation stackdriverPropagation.HTTPFormat) *addTraceIDMiddleware {
	if rootLogger == nil {
		panic("root logger must not be nil")
	}

	return &addTraceIDMiddleware{
		next:        next,
		rootLogger:  rootLogger,
		propagation: propagation,
	}
}

type addTraceIDMiddleware struct {
	// Handler is the handler used to handle the incoming request.
	next http.Handler

	// Propagation defines how traces are propagated. If unspecified,
	// Stackdriver propagation will be used.
	propagation stackdriverPropagation.HTTPFormat

	// Actual root logger to instrument with request information
	rootLogger *zap.Logger
}

func (h *addTraceIDMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	rootLogger := *h.rootLogger
	spanContext, ok := extractSpanContext(r, h.propagation)

	var logger *zap.Logger
	if !ok {
		// Not found in the header, check from the context directly than
		span := trace.FromContext(r.Context())
		if span == nil {
			traceIDField := zap.Stringer("trace_id", traceID(traceIDGenerator.NewTraceID()))
			logger = rootLogger.With(traceIDField)
		} else {
			spanContext := span.SpanContext()
			traceID := hex.EncodeToString(spanContext.TraceID[:])
			logger = rootLogger.With(zap.String("trace_id", traceID))
		}
	} else {
		traceID := hex.EncodeToString(spanContext.TraceID[:])
		logger = rootLogger.With(zap.String("trace_id", traceID))
	}

	ctx := logging.WithLogger(r.Context(), logger)
	h.next.ServeHTTP(w, r.WithContext(ctx))
}

func extractSpanContext(r *http.Request, propagation stackdriverPropagation.HTTPFormat) (trace.SpanContext, bool) {

	return propagation.SpanContextFromRequest(r)
}

type traceID [16]byte

func (t traceID) String() string {
	return hex.EncodeToString(t[:])
}
func trackingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		zlogger := logging.Logger(ctx, zlog)
		zlogger.Debug("handling HTTP request",
			zap.String("method", r.Method),
			zap.Any("host", r.Host),
			zap.Any("url", r.URL),
			zap.Any("headers", r.Header),
		)

		span := trace.FromContext(ctx)
		if span == nil {
			zlogger.Panic("trace is not present in request but should have been")
		}

		spanContext := span.SpanContext()
		traceID := spanContext.TraceID.String()

		w.Header().Set("X-Trace-ID", traceID)

		next.ServeHTTP(w, r)
	})
}

func writeError(ctx context.Context, w http.ResponseWriter, err error) {
	derr.WriteError(ctx, w, "unable to fullfil request", err)
}

func streamResponse(ctx context.Context, w http.ResponseWriter, response interface{}) {
	ctx, span := dtracing.StartSpan(ctx, "streaming JSON response", "type", fmt.Sprintf("%T", response))
	defer span.End()

	zlogger := logging.Logger(ctx, zlog)
	zlogger.Debug("streaming response")

	w.Header().Add("Content-Type", "application/json")
	if err := gojay.NewEncoder(w).Encode(response); err != nil {
		level := zapcore.ErrorLevel
		if derr.Find(err, isClientSideNetworkError) != nil {
			level = zapcore.DebugLevel
		}

		zlogger.Check(level, "an error occurred while streaming response").Write(zap.Error(err))
	}
}

func writeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) {
	ctx, span := dtracing.StartSpan(ctx, "writing JSON response", "type", fmt.Sprintf("%T", response))
	defer span.End()

	zlogger := logging.Logger(ctx, zlog)
	zlogger.Debug("writing response")

	w.Header().Set("Content-type", "application/json")

	if err := json.NewEncoder(w).Encode(response); err != nil {
		level := zapcore.ErrorLevel
		if isClientSideNetworkError(err) {
			level = zapcore.DebugLevel
		}

		zlogger.Check(level, "an error occurred while writing response").Write(zap.Error(err))
	}
}

func extractReadRequestCommon(r *http.Request) *readRequestCommon {
	blockNum64, _ := strconv.ParseInt(r.FormValue("block_num"), 10, 64)
	offset, _ := intInput(r.FormValue("offset"), 0)
	limit, _ := intInput(r.FormValue("limit"), 0)

	return &readRequestCommon{
		BlockNum:     uint64(blockNum64),
		Offset:       offset,
		Limit:        limit,
		Key:          r.FormValue("key"),
		KeyType:      r.FormValue("key_type"),
		ToJSON:       boolInput(r.FormValue("json")),
		WithABI:      boolInput(r.FormValue("with_abi")),
		WithBlockNum: boolInput(r.FormValue("with_block_num")),
	}
}

func isClientSideNetworkError(err error) bool {
	netErr, isNetErr := err.(*net.OpError)
	if !isNetErr {
		return false
	}

	syscallErr, isSyscallErr := netErr.Err.(*os.SyscallError)
	if !isSyscallErr {
		return false
	}

	return syscallErr.Err == syscall.ECONNRESET || syscallErr.Err == syscall.EPIPE
}

func boolInput(in string) bool {
	return in == "true" || in == "1"
}

func intInput(in string, defaultValue int) (int, error) {
	if in == "" {
		return defaultValue, nil
	}

	return strconv.Atoi(in)
}
