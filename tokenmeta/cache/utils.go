package cache

import (
	pbtokenmeta "github.com/zhongshuwen/histnew/pb/dfuse/zswhq/tokenmeta/v1"
	"github.com/zhongshuwen/zswchain-go"
)

func ProtoEOSAccountBalanceToOwnedAsset(bal *pbtokenmeta.AccountBalance) *OwnedAsset {
	return &OwnedAsset{
		Owner: zsw.AccountName(bal.Account),
		Asset: &zsw.ExtendedAsset{
			Asset: zsw.Asset{
				Amount: zsw.Int64(bal.Amount),
				Symbol: zsw.Symbol{
					Precision: uint8(bal.Precision),
					Symbol:    bal.Symbol,
				},
			},
			Contract: zsw.AccountName(bal.TokenContract),
		},
	}
}

func AssetToProtoAccountBalance(asset *OwnedAsset) *pbtokenmeta.AccountBalance {
	return &pbtokenmeta.AccountBalance{
		TokenContract: string(asset.Asset.Contract),
		Account:       string(asset.Owner),
		Amount:        uint64(asset.Asset.Asset.Amount),
		Precision:     uint32(asset.Asset.Asset.Precision),
		Symbol:        asset.Asset.Asset.Symbol.Symbol,
	}
}

func lessValueToBool(value int, order SortingOrder) bool {
	if order == ASC {
		if value > 0 {
			return false
		}

		if value < 0 {
			return true
		}
	}
	if value < 0 {
		return false
	}

	if value > 0 {
		return true
	}

	return false
}
