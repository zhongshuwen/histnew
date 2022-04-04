package blockmeta

import (
	"context"

	"github.com/zhongshuwen/zswchain-go"
	"github.com/streamingfast/blockmeta"
)

func init() {
	blockmeta.GetBlockNumFromID = blockNumFromID
}

func blockNumFromID(ctx context.Context, id string) (uint64, error) {
	return uint64(zsw.BlockNum(id)), nil
}
