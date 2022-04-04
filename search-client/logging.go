package searchclient

import (
	"os"

	"github.com/streamingfast/logging"
	"go.uber.org/zap"
)

var traceEnabled = false
var zlog = zap.NewNop()

func init() {
	logging.Register("github.com/zhongshuwen/histnew/search-client", &zlog)

	if os.Getenv("TRACE") == "true" {
		traceEnabled = true
	}
}
