package accounthist

import (
	"fmt"

	"github.com/zhongshuwen/zswchain-go"

	"github.com/zhongshuwen/histnew/accounthist/keyer"
)

type AccountFacet uint64

func (a AccountFacet) String() string {
	return fmt.Sprintf("account (%s)", zsw.NameToString(uint64(a)))
}

func (a AccountFacet) Account() uint64 {
	return uint64(a)
}

func (a AccountFacet) Row(shard byte, seqData uint64) RowKey {
	return keyer.EncodeAccountKey(uint64(a), shard, seqData)
}

func (a AccountFacet) Bytes() []byte {
	return keyer.EncodeAccountPrefixKey(uint64(a))
}
