package search

import (
	"github.com/streamingfast/logging"
	"go.uber.org/zap"
)

var traceEnabled = logging.IsTraceEnabled("search", "github.com/zhongshuwen/histnew/search")
var zlog *zap.Logger

func init() {
	logging.Register("github.com/zhongshuwen/histnew/search", &zlog)
}
