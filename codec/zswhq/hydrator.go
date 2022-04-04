package zswhq

import pbcodec "github.com/zhongshuwen/histnew/pb/dfuse/zswhq/codec/v1"

type Hydrator interface {
	// HydrateBlock decodes the received Deep Mind AcceptedBlock data structure against the
	// correct struct for this version of EOSIO supported by this hydrator.
	HydrateBlock(block *pbcodec.Block, input []byte) error

	// DecodeTransactionTrace decodes the received Deep Mind AppliedTransaction data structure against the
	// correct struct for this version of EOSIO supported by this hydrator.
	DecodeTransactionTrace(input []byte, opts ...ConversionOption) (*pbcodec.TransactionTrace, error)
}
