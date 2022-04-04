package zswhq

import (
	"github.com/zhongshuwen/zswchain-go"
	"github.com/zhongshuwen/zswchain-go/ecc"
)

// BlockState
//
// File hierarchy:
//  - https://github.com/EOSIO/eos/blob/v2.1.0/libraries/chain/include/zswhq/chain/block_header_state.hpp#L57
//  - https://github.com/EOSIO/eos/blob/v2.1.0/libraries/chain/include/zswhq/chain/block_header_state.hpp#L126
//  - https://github.com/EOSIO/eos/blob/v2.1.0/libraries/chain/include/zswhq/chain/block_state.hpp#L10
type BlockState struct {
	// From 'struct block_header_state_common'
	BlockNum                         uint32                         `json:"block_num"`
	DPoSProposedIrreversibleBlockNum uint32                         `json:"dpos_proposed_irreversible_blocknum"`
	DPoSIrreversibleBlockNum         uint32                         `json:"dpos_irreversible_blocknum"`
	ActiveSchedule                   *zsw.ProducerAuthoritySchedule `json:"active_schedule"`
	BlockrootMerkle                  *zsw.MerkleRoot                `json:"blockroot_merkle,omitempty"`
	ProducerToLastProduced           []zsw.PairAccountNameBlockNum  `json:"producer_to_last_produced,omitempty"`
	ProducerToLastImpliedIRB         []zsw.PairAccountNameBlockNum  `json:"producer_to_last_implied_irb,omitempty"`
	ValidBlockSigningAuthorityV2     *zsw.BlockSigningAuthority     `json:"valid_block_signing_authority,omitempty"`
	ConfirmCount                     []uint8                        `json:"confirm_count,omitempty"`

	// From 'struct block_header_state'
	BlockID                   zsw.Checksum256                   `json:"id"`
	Header                    *zsw.SignedBlockHeader            `json:"header,omitempty"`
	PendingSchedule           *zsw.PendingSchedule              `json:"pending_schedule"`
	ActivatedProtocolFeatures *zsw.ProtocolFeatureActivationSet `json:"activated_protocol_features,omitempty" eos:"optional"`
	AdditionalSignatures      []ecc.Signature                   `json:"additional_signatures"`

	// From 'struct block_state'
	// Type changed in v2.1.x
	SignedBlock *SignedBlock `json:"block,omitempty" eos:"optional"`
	Validated   bool         `json:"validated"`
}

// BlockState
//
// File hierarchy:
//  - https://github.com/EOSIO/eos/blob/v2.1.0/libraries/chain/include/zswhq/chain/block.hpp#L135
type SignedBlock struct {
	zsw.SignedBlockHeader
	// Added in v2.1.x
	PruneState uint8 `json:"prune_state"`
	// Type changed in v2.1.x
	Transactions    []*TransactionReceipt `json:"transactions"`
	BlockExtensions []*zsw.Extension      `json:"block_extensions"`
}

// TransactionTrace
//
// File hierarchy:
//  - https://github.com/EOSIO/eos/blob/v2.1.0/libraries/chain/include/zswhq/chain/trace.hpp#L51
type TransactionTrace struct {
	ID              zsw.Checksum256               `json:"id"`
	BlockNum        uint32                        `json:"block_num"`
	BlockTime       zsw.BlockTimestamp            `json:"block_time"`
	ProducerBlockID zsw.Checksum256               `json:"producer_block_id" eos:"optional"`
	Receipt         *zsw.TransactionReceiptHeader `json:"receipt,omitempty" eos:"optional"`
	Elapsed         zsw.Int64                     `json:"elapsed"`
	NetUsage        zsw.Uint64                    `json:"net_usage"`
	Scheduled       bool                          `json:"scheduled"`
	ActionTraces    []*ActionTrace                `json:"action_traces"`
	AccountRamDelta *AccountDelta                 `json:"account_ram_delta" eos:"optional"`
	FailedDtrxTrace *TransactionTrace             `json:"failed_dtrx_trace,omitempty" eos:"optional"`
	Except          *zsw.Except                   `json:"except,omitempty" eos:"optional"`
	ErrorCode       *zsw.Uint64                   `json:"error_code,omitempty" eos:"optional"`
}

// TransactionTrace
//
// File hierarchy:
//  - https://github.com/EOSIO/eos/blob/v2.1.0/libraries/chain/include/zswhq/chain/trace.hpp#L22
type ActionTrace struct {
	ActionOrdinal                          zsw.Varuint32           `json:"action_ordinal"`
	CreatorActionOrdinal                   zsw.Varuint32           `json:"creator_action_ordinal"`
	ClosestUnnotifiedAncestorActionOrdinal zsw.Varuint32           `json:"closest_unnotified_ancestor_action_ordinal"`
	Receipt                                *zsw.ActionTraceReceipt `json:"receipt,omitempty" eos:"optional"`
	Receiver                               zsw.AccountName         `json:"receiver"`
	Action                                 *zsw.Action             `json:"act"`
	ContextFree                            bool                    `json:"context_free"`
	ElapsedUs                              zsw.Int64               `json:"elapsed"`
	Console                                zsw.SafeString          `json:"console"`
	TransactionID                          zsw.Checksum256         `json:"trx_id"`
	BlockNum                               uint32                  `json:"block_num"`
	BlockTime                              zsw.BlockTimestamp      `json:"block_time"`
	ProducerBlockID                        zsw.Checksum256         `json:"producer_block_id" eos:"optional"`
	AccountRAMDeltas                       []AccountDelta          `json:"account_ram_deltas"`
	// Added in 2.1.x
	AccountDiskDeltas []AccountDelta `json:"account_disk_deltas"`
	Except            *zsw.Except    `json:"except,omitempty" eos:"optional"`
	ErrorCode         *zsw.Uint64    `json:"error_code,omitempty" eos:"optional"`
	// Added in 2.1.x
	ReturnValue zsw.HexBytes `json:"return_value"`
}

type AccountDelta struct {
	Account zsw.AccountName `json:"account"`
	Delta   zsw.Int64       `json:"delta"`
}

type TransactionReceipt struct {
	zsw.TransactionReceiptHeader
	Transaction Transaction `json:"trx"`
}

var TransactionVariant = zsw.NewVariantDefinition([]zsw.VariantType{
	{Name: "transaction_id", Type: zsw.Checksum256{}},
	{Name: "packed_transaction", Type: (*PackedTransaction)(nil)},
})

type Transaction struct {
	zsw.BaseVariant
}

func (r *Transaction) UnmarshalBinary(decoder *zsw.Decoder) error {
	return r.BaseVariant.UnmarshalBinaryVariant(decoder, TransactionVariant)
}

type PackedTransaction struct {
	Compression       zsw.CompressionType `json:"compression"`
	PrunableData      *PrunableData       `json:"prunable_data"`
	PackedTransaction zsw.HexBytes        `json:"packed_trx"`
}

var PrunableDataVariant = zsw.NewVariantDefinition([]zsw.VariantType{
	{Name: "full_legacy", Type: (*PackedTransactionPrunableFullLegacy)(nil)},
	{Name: "none", Type: (*PackedTransactionPrunableNone)(nil)},
	{Name: "partial", Type: (*PackedTransactionPrunablePartial)(nil)},
	{Name: "full", Type: (*PackedTransactionPrunableFull)(nil)},
})

type PackedTransactionPrunableNone struct {
	Digest zsw.Checksum256 `json:"digest"`
}

type PackedTransactionPrunablePartial struct {
	Signatures          []ecc.Signature `json:"signatures"`
	ContextFreeSegments []*Segment      `json:"context_free_segments"`
}

type PackedTransactionPrunableFull struct {
	Signatures          []ecc.Signature `json:"signatures"`
	ContextFreeSegments []zsw.HexBytes  `json:"context_free_segments"`
}

type PackedTransactionPrunableFullLegacy struct {
	Signatures            []ecc.Signature `json:"signatures"`
	PackedContextFreeData zsw.HexBytes    `json:"packed_context_free_data"`
}

type PrunableData struct {
	zsw.BaseVariant
}

func (r *PrunableData) UnmarshalBinary(decoder *zsw.Decoder) error {
	return r.BaseVariant.UnmarshalBinaryVariant(decoder, PrunableDataVariant)
}

var SegmentVariant = zsw.NewVariantDefinition([]zsw.VariantType{
	{Name: "digest", Type: zsw.Checksum256{}},
	{Name: "bytes", Type: zsw.HexBytes{}},
})

type Segment struct {
	zsw.BaseVariant
}

func (r *Segment) UnmarshalBinary(decoder *zsw.Decoder) error {
	return r.BaseVariant.UnmarshalBinaryVariant(decoder, SegmentVariant)
}
