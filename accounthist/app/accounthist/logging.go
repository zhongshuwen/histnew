package tokenhist

import (
	"github.com/streamingfast/logging"
	"go.uber.org/zap"
)

var zlog *zap.Logger

func init() {
	logging.Register("github.com/zhongshuwen/histnew/accounthist/app/accounthist", &zlog)
}
