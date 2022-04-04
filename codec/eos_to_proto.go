// Copyright 2019 dfuse Platform Inc.
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

package codec

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sort"
	"time"

	pbcodec "github.com/zhongshuwen/histnew/pb/dfuse/zswhq/codec/v1"
	"github.com/zhongshuwen/zswchain-go"
	"github.com/zhongshuwen/zswchain-go/ecc"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"
)

func checksumsToBytesSlices(in []zsw.Checksum256) [][]byte {
	out := make([][]byte, len(in))
	for i, s := range in {
		out[i] = s
	}
	return out
}

func hexBytesToBytesSlices(in []zsw.HexBytes) [][]byte {
	out := make([][]byte, len(in))
	for i, s := range in {
		out[i] = s
	}
	return out
}

func bytesSlicesToHexBytes(in [][]byte) []zsw.HexBytes {
	out := make([]zsw.HexBytes, len(in))
	for i, s := range in {
		out[i] = s
	}
	return out
}

func BlockHeaderToEOS(in *pbcodec.BlockHeader) *zsw.BlockHeader {
	stamp, _ := ptypes.Timestamp(in.Timestamp)
	prev, _ := hex.DecodeString(in.Previous)
	out := &zsw.BlockHeader{
		Timestamp:        zsw.BlockTimestamp{Time: stamp},
		Producer:         zsw.AccountName(in.Producer),
		Confirmed:        uint16(in.Confirmed),
		Previous:         prev,
		TransactionMRoot: in.TransactionMroot,
		ActionMRoot:      in.ActionMroot,
		ScheduleVersion:  in.ScheduleVersion,
		HeaderExtensions: ExtensionsToEOS(in.HeaderExtensions),
	}

	if in.NewProducersV1 != nil {
		out.NewProducersV1 = ProducerScheduleToEOS(in.NewProducersV1)
	}

	return out
}

func BlockSigningAuthorityToEOS(in *pbcodec.BlockSigningAuthority) *zsw.BlockSigningAuthority {
	switch v := in.Variant.(type) {
	case *pbcodec.BlockSigningAuthority_V0:
		return &zsw.BlockSigningAuthority{
			BaseVariant: zsw.BaseVariant{
				TypeID: zsw.BlockSigningAuthorityVariant.TypeID("block_signing_authority_v0"),
				Impl: zsw.BlockSigningAuthorityV0{
					Threshold: v.V0.Threshold,
					Keys:      KeyWeightsPToEOS(v.V0.Keys),
				},
			},
		}
	default:
		panic(fmt.Errorf("unknown block signing authority variant %t", in.Variant))
	}
}

func ProducerScheduleToEOS(in *pbcodec.ProducerSchedule) *zsw.ProducerSchedule {
	return &zsw.ProducerSchedule{
		Version:   in.Version,
		Producers: ProducerKeysToEOS(in.Producers),
	}
}

func ProducerAuthorityScheduleToEOS(in *pbcodec.ProducerAuthoritySchedule) *zsw.ProducerAuthoritySchedule {
	return &zsw.ProducerAuthoritySchedule{
		Version:   in.Version,
		Producers: ProducerAuthoritiesToEOS(in.Producers),
	}
}

func ProducerKeysToEOS(in []*pbcodec.ProducerKey) (out []zsw.ProducerKey) {
	out = make([]zsw.ProducerKey, len(in))
	for i, producer := range in {
		// panic on error instead?
		key, _ := ecc.NewPublicKey(producer.BlockSigningKey)

		out[i] = zsw.ProducerKey{
			AccountName:     zsw.AccountName(producer.AccountName),
			BlockSigningKey: key,
		}
	}
	return
}

func PublicKeysToEOS(in []string) (out []*ecc.PublicKey) {
	if len(in) <= 0 {
		return nil
	}
	out = make([]*ecc.PublicKey, len(in))
	for i, inkey := range in {
		// panic on error instead?
		key, _ := ecc.NewPublicKey(inkey)

		out[i] = &key
	}
	return
}

func ExtensionsToEOS(in []*pbcodec.Extension) (out []*zsw.Extension) {
	if len(in) <= 0 {
		return nil
	}

	out = make([]*zsw.Extension, len(in))
	for i, extension := range in {
		out[i] = &zsw.Extension{
			Type: uint16(extension.Type),
			Data: extension.Data,
		}
	}
	return
}

func ProducerAuthoritiesToEOS(producerAuthorities []*pbcodec.ProducerAuthority) (out []*zsw.ProducerAuthority) {
	if len(producerAuthorities) <= 0 {
		return nil
	}

	out = make([]*zsw.ProducerAuthority, len(producerAuthorities))
	for i, authority := range producerAuthorities {
		out[i] = &zsw.ProducerAuthority{
			AccountName:           zsw.AccountName(authority.AccountName),
			BlockSigningAuthority: BlockSigningAuthorityToEOS(authority.BlockSigningAuthority),
		}
	}
	return
}

func TransactionReceiptHeaderToEOS(in *pbcodec.TransactionReceiptHeader) *zsw.TransactionReceiptHeader {
	return &zsw.TransactionReceiptHeader{
		Status:               TransactionStatusToEOS(in.Status),
		CPUUsageMicroSeconds: in.CpuUsageMicroSeconds,
		NetUsageWords:        zsw.Varuint32(in.NetUsageWords),
	}
}

func SignaturesToEOS(in []string) []ecc.Signature {
	out := make([]ecc.Signature, len(in))
	for i, signature := range in {
		sig, err := ecc.NewSignature(signature)
		if err != nil {
			panic(fmt.Sprintf("failed to read signature %q: %s", signature, err))
		}

		out[i] = sig
	}
	return out
}

func TransactionTraceToEOS(in *pbcodec.TransactionTrace) (out *zsw.TransactionTrace) {
	out = &zsw.TransactionTrace{
		ID:              ChecksumToEOS(in.Id),
		BlockNum:        uint32(in.BlockNum),
		BlockTime:       TimestampToBlockTimestamp(in.BlockTime),
		ProducerBlockID: ChecksumToEOS(in.ProducerBlockId),
		Elapsed:         zsw.Int64(in.Elapsed),
		NetUsage:        zsw.Uint64(in.NetUsage),
		Scheduled:       in.Scheduled,
		ActionTraces:    ActionTracesToEOS(in.ActionTraces),
		Except:          ExceptionToEOS(in.Exception),
		ErrorCode:       ErrorCodeToEOS(in.ErrorCode),
	}

	if in.FailedDtrxTrace != nil {
		out.FailedDtrxTrace = TransactionTraceToEOS(in.FailedDtrxTrace)
	}
	if in.Receipt != nil {
		out.Receipt = TransactionReceiptHeaderToEOS(in.Receipt)
	}

	return out
}

func AuthoritiesToEOS(authority *pbcodec.Authority) zsw.Authority {
	return zsw.Authority{
		Threshold: authority.Threshold,
		Keys:      KeyWeightsToEOS(authority.Keys),
		Accounts:  PermissionLevelWeightsToEOS(authority.Accounts),
		Waits:     WaitWeightsToEOS(authority.Waits),
	}
}

func WaitWeightsToEOS(waits []*pbcodec.WaitWeight) (out []zsw.WaitWeight) {
	if len(waits) <= 0 {
		return nil
	}

	out = make([]zsw.WaitWeight, len(waits))
	for i, o := range waits {
		out[i] = zsw.WaitWeight{
			WaitSec: o.WaitSec,
			Weight:  uint16(o.Weight),
		}
	}
	return out
}

func PermissionLevelWeightsToEOS(weights []*pbcodec.PermissionLevelWeight) (out []zsw.PermissionLevelWeight) {
	if len(weights) == 0 {
		return []zsw.PermissionLevelWeight{}
	}

	out = make([]zsw.PermissionLevelWeight, len(weights))
	for i, o := range weights {
		out[i] = zsw.PermissionLevelWeight{
			Permission: PermissionLevelToEOS(o.Permission),
			Weight:     uint16(o.Weight),
		}
	}
	return
}

func PermissionLevelToEOS(perm *pbcodec.PermissionLevel) zsw.PermissionLevel {
	return zsw.PermissionLevel{
		Actor:      zsw.AccountName(perm.Actor),
		Permission: zsw.PermissionName(perm.Permission),
	}
}

func KeyWeightsToEOS(keys []*pbcodec.KeyWeight) (out []zsw.KeyWeight) {
	if len(keys) <= 0 {
		return nil
	}

	out = make([]zsw.KeyWeight, len(keys))
	for i, o := range keys {
		out[i] = zsw.KeyWeight{
			PublicKey: ecc.MustNewPublicKey(o.PublicKey),
			Weight:    uint16(o.Weight),
		}
	}
	return

}

func KeyWeightsPToEOS(keys []*pbcodec.KeyWeight) (out []*zsw.KeyWeight) {
	if len(keys) <= 0 {
		return nil
	}

	out = make([]*zsw.KeyWeight, len(keys))
	for i, o := range keys {
		out[i] = &zsw.KeyWeight{
			PublicKey: ecc.MustNewPublicKey(o.PublicKey),
			Weight:    uint16(o.Weight),
		}
	}
	return

}

func TransactionToEOS(trx *pbcodec.Transaction) *zsw.Transaction {
	var contextFreeActions []*zsw.Action
	if len(trx.ContextFreeActions) > 0 {
		contextFreeActions = make([]*zsw.Action, len(trx.ContextFreeActions))
		for i, act := range trx.ContextFreeActions {
			contextFreeActions[i] = ActionToEOS(act)
		}
	}

	var actions []*zsw.Action
	if len(trx.Actions) > 0 {
		actions = make([]*zsw.Action, len(trx.Actions))
		for i, act := range trx.Actions {
			actions[i] = ActionToEOS(act)
		}
	}

	return &zsw.Transaction{
		TransactionHeader:  *(TransactionHeaderToEOS(trx.Header)),
		ContextFreeActions: contextFreeActions,
		Actions:            actions,
		Extensions:         ExtensionsToEOS(trx.Extensions),
	}
}

func TransactionHeaderToEOS(trx *pbcodec.TransactionHeader) *zsw.TransactionHeader {
	out := &zsw.TransactionHeader{
		Expiration:       TimestampToJSONTime(trx.Expiration),
		RefBlockNum:      uint16(trx.RefBlockNum),
		RefBlockPrefix:   uint32(trx.RefBlockPrefix),
		MaxNetUsageWords: zsw.Varuint32(trx.MaxNetUsageWords),
		MaxCPUUsageMS:    uint8(trx.MaxCpuUsageMs),
		DelaySec:         zsw.Varuint32(trx.DelaySec),
	}

	return out
}

func SignedTransactionToEOS(trx *pbcodec.SignedTransaction) *zsw.SignedTransaction {
	return &zsw.SignedTransaction{
		Transaction:     TransactionToEOS(trx.Transaction),
		Signatures:      SignaturesToEOS(trx.Signatures),
		ContextFreeData: bytesSlicesToHexBytes(trx.ContextFreeData),
	}
}

func ActionTracesToEOS(actionTraces []*pbcodec.ActionTrace) (out []zsw.ActionTrace) {
	if len(actionTraces) <= 0 {
		return nil
	}

	out = make([]zsw.ActionTrace, len(actionTraces))
	for i, actionTrace := range actionTraces {
		out[i] = ActionTraceToEOS(actionTrace)
	}

	sort.Slice(out, func(i, j int) bool { return out[i].ActionOrdinal < out[j].ActionOrdinal })

	return
}

func AuthSequenceListToEOS(in []*pbcodec.AuthSequence) (out []zsw.TransactionTraceAuthSequence) {
	if len(in) == 0 {
		return []zsw.TransactionTraceAuthSequence{}
	}

	out = make([]zsw.TransactionTraceAuthSequence, len(in))
	for i, seq := range in {
		out[i] = AuthSequenceToEOS(seq)
	}

	return
}

func AuthSequenceToEOS(in *pbcodec.AuthSequence) zsw.TransactionTraceAuthSequence {
	return zsw.TransactionTraceAuthSequence{
		Account:  zsw.AccountName(in.AccountName),
		Sequence: zsw.Uint64(in.Sequence),
	}
}

func ErrorCodeToEOS(in uint64) *zsw.Uint64 {
	if in != 0 {
		val := zsw.Uint64(in)
		return &val
	}
	return nil
}

func ActionTraceToEOS(in *pbcodec.ActionTrace) (out zsw.ActionTrace) {
	out = zsw.ActionTrace{
		Receiver:             zsw.AccountName(in.Receiver),
		Action:               ActionToEOS(in.Action),
		Elapsed:              zsw.Int64(in.Elapsed),
		Console:              zsw.SafeString(in.Console),
		TransactionID:        ChecksumToEOS(in.TransactionId),
		ContextFree:          in.ContextFree,
		ProducerBlockID:      ChecksumToEOS(in.ProducerBlockId),
		BlockNum:             uint32(in.BlockNum),
		BlockTime:            TimestampToBlockTimestamp(in.BlockTime),
		AccountRAMDeltas:     AccountRAMDeltasToEOS(in.AccountRamDeltas),
		Except:               ExceptionToEOS(in.Exception),
		ActionOrdinal:        zsw.Varuint32(in.ActionOrdinal),
		CreatorActionOrdinal: zsw.Varuint32(in.CreatorActionOrdinal),
		ErrorCode:            ErrorCodeToEOS(in.ErrorCode),
	}
	out.ClosestUnnotifiedAncestorActionOrdinal = zsw.Varuint32(in.ClosestUnnotifiedAncestorActionOrdinal) // freaking long line, stay away from me

	if in.Receipt != nil {
		receipt := in.Receipt

		out.Receipt = &zsw.ActionTraceReceipt{
			Receiver:        zsw.AccountName(receipt.Receiver),
			ActionDigest:    ChecksumToEOS(receipt.Digest),
			GlobalSequence:  zsw.Uint64(receipt.GlobalSequence),
			AuthSequence:    AuthSequenceListToEOS(receipt.AuthSequence),
			ReceiveSequence: zsw.Uint64(receipt.RecvSequence),
			CodeSequence:    zsw.Varuint32(receipt.CodeSequence),
			ABISequence:     zsw.Varuint32(receipt.AbiSequence),
		}
	}

	return
}

func ChecksumToEOS(in string) zsw.Checksum256 {
	out, err := hex.DecodeString(in)
	if err != nil {
		panic(fmt.Sprintf("failed decoding checksum %q: %s", in, err))
	}

	return zsw.Checksum256(out)
}

func ActionToEOS(action *pbcodec.Action) (out *zsw.Action) {
	d := zsw.ActionData{}
	d.SetToServer(false) // rather, what we expect FROM `nodeos` servers

	d.HexData = zsw.HexBytes(action.RawData)
	if len(action.JsonData) != 0 {
		err := json.Unmarshal([]byte(action.JsonData), &d.Data)
		if err != nil {
			panic(fmt.Sprintf("unmarshaling action json data %q: %s", action.JsonData, err))
		}
	}

	out = &zsw.Action{
		Account:       zsw.AccountName(action.Account),
		Name:          zsw.ActionName(action.Name),
		Authorization: AuthorizationToEOS(action.Authorization),
		ActionData:    d,
	}

	return out
}

func AuthorizationToEOS(authorization []*pbcodec.PermissionLevel) (out []zsw.PermissionLevel) {
	if len(authorization) == 0 {
		return []zsw.PermissionLevel{}
	}

	out = make([]zsw.PermissionLevel, len(authorization))
	for i, permission := range authorization {
		out[i] = PermissionLevelToEOS(permission)
	}
	return
}

func AccountRAMDeltasToEOS(deltas []*pbcodec.AccountRAMDelta) (out []*zsw.AccountRAMDelta) {
	if len(deltas) == 0 {
		return []*zsw.AccountRAMDelta{}
	}

	out = make([]*zsw.AccountRAMDelta, len(deltas))
	for i, delta := range deltas {
		out[i] = &zsw.AccountRAMDelta{
			Account: zsw.AccountName(delta.Account),
			Delta:   zsw.Int64(delta.Delta),
		}
	}
	return
}

func ExceptionToEOS(in *pbcodec.Exception) *zsw.Except {
	if in == nil {
		return nil
	}
	out := &zsw.Except{
		Code:    zsw.Int64(in.Code),
		Name:    in.Name,
		Message: in.Message,
	}

	if len(in.Stack) > 0 {
		out.Stack = make([]*zsw.ExceptLogMessage, len(in.Stack))
		for i, el := range in.Stack {
			msg := &zsw.ExceptLogMessage{
				Format: el.Format,
			}

			ctx := LogContextToEOS(el.Context)
			if ctx != nil {
				msg.Context = *ctx
			}

			if len(el.Data) > 0 {
				msg.Data = json.RawMessage(el.Data)
			}

			out.Stack[i] = msg
		}
	}

	return out
}

func LogContextToEOS(in *pbcodec.Exception_LogContext) *zsw.ExceptLogContext {
	if in == nil {
		return nil
	}

	var exceptLevel zsw.ExceptLogLevel
	exceptLevel.FromString(in.Level)

	return &zsw.ExceptLogContext{
		Level:      exceptLevel,
		File:       in.File,
		Line:       uint64(in.Line),
		Method:     in.Method,
		Hostname:   in.Hostname,
		ThreadName: in.ThreadName,
		Timestamp:  TimestampToJSONTime(in.Timestamp),
		Context:    LogContextToEOS(in.Context),
	}
}

func TimestampToJSONTime(in *timestamp.Timestamp) zsw.JSONTime {
	out, _ := ptypes.Timestamp(in)
	return zsw.JSONTime{Time: out}
}

func TimestampToBlockTimestamp(in *timestamp.Timestamp) zsw.BlockTimestamp {
	out, _ := ptypes.Timestamp(in)
	return zsw.BlockTimestamp{Time: out}
}

func TransactionStatusToEOS(in pbcodec.TransactionStatus) zsw.TransactionStatus {
	switch in {
	case pbcodec.TransactionStatus_TRANSACTIONSTATUS_EXECUTED:
		return zsw.TransactionStatusExecuted
	case pbcodec.TransactionStatus_TRANSACTIONSTATUS_SOFTFAIL:
		return zsw.TransactionStatusSoftFail
	case pbcodec.TransactionStatus_TRANSACTIONSTATUS_HARDFAIL:
		return zsw.TransactionStatusHardFail
	case pbcodec.TransactionStatus_TRANSACTIONSTATUS_DELAYED:
		return zsw.TransactionStatusDelayed
	case pbcodec.TransactionStatus_TRANSACTIONSTATUS_EXPIRED:
		return zsw.TransactionStatusExpired
	default:
		return zsw.TransactionStatusUnknown
	}
}

func ExtractEOSSignedTransactionFromReceipt(trxReceipt *pbcodec.TransactionReceipt) (*zsw.SignedTransaction, error) {
	eosPackedTx, err := pbcodecPackedTransactionToEOS(trxReceipt.PackedTransaction)
	if err != nil {
		return nil, fmt.Errorf("pbcodec.PackedTransaction to EOS conversion failed: %s", err)
	}

	signedTransaction, err := eosPackedTx.UnpackBare()
	if err != nil {
		return nil, fmt.Errorf("unable to unpack packed transaction: %s", err)
	}

	return signedTransaction, nil
}

func mustProtoTimestamp(in time.Time) *timestamp.Timestamp {
	out, err := ptypes.TimestampProto(in)
	if err != nil {
		panic(fmt.Sprintf("invalid timestamp conversion %q: %s", in, err))
	}
	return out
}

func pbcodecPackedTransactionToEOS(packedTrx *pbcodec.PackedTransaction) (*zsw.PackedTransaction, error) {
	signatures := make([]ecc.Signature, len(packedTrx.Signatures))
	for i, signature := range packedTrx.Signatures {
		eccSignature, err := ecc.NewSignature(signature)
		if err != nil {
			return nil, err
		}

		signatures[i] = eccSignature
	}

	return &zsw.PackedTransaction{
		Signatures:            signatures,
		Compression:           zsw.CompressionType(packedTrx.Compression),
		PackedContextFreeData: packedTrx.PackedContextFreeData,
		PackedTransaction:     packedTrx.PackedTransaction,
	}, nil
}
