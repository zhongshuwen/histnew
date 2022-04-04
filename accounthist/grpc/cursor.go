package grpc

import pbaccounthist "github.com/zhongshuwen/histnew/pb/dfuse/zswhq/accounthist/v1"

const CursorMagicValue = 4374

func ActionKeyToCursor(key []byte, shardNo byte, seqNum uint64) *pbaccounthist.Cursor {
	return &pbaccounthist.Cursor{
		Version:        0,
		Magic:          CursorMagicValue,
		Key:            key,
		ShardNum:       uint32(shardNo),
		SequenceNumber: seqNum,
	}
}
