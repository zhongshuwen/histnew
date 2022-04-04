package cache

import (
	"github.com/streamingfast/logging"
	"go.uber.org/zap"
)

var zlog = zap.NewNop()

func init() {
	logging.Register("github.com/zhongshuwen/histnew/tokenmeta/cache", &zlog)
}
