package tokenmeta

import (
	"encoding/json"
	"testing"

	pbtokenmeta "github.com/zhongshuwen/histnew/pb/dfuse/zswhq/tokenmeta/v1"
	"github.com/zhongshuwen/zswchain-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAbi_getAccountBalanceFromDBRow(t *testing.T) {
	tests := []struct {
		name        string
		scope       string
		contract    zsw.AccountName
		symbol      *zsw.Symbol
		dbRow       json.RawMessage
		expectValue *pbtokenmeta.AccountBalance
		expectError bool
	}{
		{
			name:     "simple raw message",
			scope:    "eoscanadacom",
			contract: "zswhq.token",
			symbol: &zsw.Symbol{
				Precision: 4,
				Symbol:    "ZSWCC",
			},
			dbRow: []byte("{\"balance\":\"0.0064 EOS\"}"),
			expectValue: &pbtokenmeta.AccountBalance{
				TokenContract: "zswhq.token",
				Account:       "eoscanadacom",
				Amount:        64,
				Precision:     4,
				Symbol:        "ZSWCC",
			},
		},
		{
			name:     "invalid raw message",
			scope:    "eoscanadacom",
			contract: "zswhq.token",
			symbol: &zsw.Symbol{
				Precision: 4,
				Symbol:    "ZSWCC",
			},
			dbRow:       []byte("{\"bbbbalances\":\"0.0064 EOS\"}"),
			expectError: true,
		},
		{
			name:     "invalid balance in raw message",
			scope:    "eoscanadacom",
			contract: "zswhq.token",
			symbol: &zsw.Symbol{
				Precision: 4,
				Symbol:    "ZSWCC",
			},
			dbRow:       []byte("{\"balance\":\"0.0064EOS\"}"),
			expectError: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			value, err := getAccountBalanceFromDBRow(test.contract, test.symbol, test.scope, test.dbRow)

			if test.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, test.expectValue, value)
			}

		})
	}
}

func TestAbi_getTokenFromDBRow(t *testing.T) {
	tests := []struct {
		name        string
		contract    zsw.AccountName
		symbol      *zsw.Symbol
		dbRow       json.RawMessage
		expectValue *pbtokenmeta.Token
		expectError bool
	}{
		{
			name:     "updating an existing token",
			contract: "zswhq.token",
			symbol: &zsw.Symbol{
				Precision: 4,
				Symbol:    "ZSWCC",
			},
			dbRow: []byte("{\"supply\":\"1074146295.0178 EOS\",\"max_supply\":\"10000000000.0000 EOS\",\"issuer\":\"zswhq\"}"),
			expectValue: &pbtokenmeta.Token{
				Contract:      "zswhq.token",
				Symbol:        "ZSWCC",
				Precision:     4,
				Issuer:        "zswhq",
				MaximumSupply: 100000000000000,
				TotalSupply:   10741462950178,
				Holders:       0,
			},
		},
		{
			name:     "creating a new token",
			contract: "eoscanadaca",
			dbRow:    []byte("{\"supply\":\"52000.00000 WAL\",\"max_supply\":\"100000.00000 WAL\",\"issuer\":\"eoscanadacom\"}"),
			expectValue: &pbtokenmeta.Token{
				Contract:      "eoscanadaca",
				Symbol:        "WAL",
				Precision:     5,
				Issuer:        "eoscanadacom",
				MaximumSupply: 10000000000,
				TotalSupply:   5200000000,
				Holders:       0,
			},
		},
		{
			name:     "invalid raw message",
			contract: "zswhq.token",
			symbol: &zsw.Symbol{
				Precision: 4,
				Symbol:    "ZSWCC",
			},
			dbRow:       []byte("{\"balance\":\"234\"}"),
			expectError: true,
		},
		{
			name:     "invalid balance in raw message",
			contract: "zswhq.token",
			symbol: &zsw.Symbol{
				Precision: 4,
				Symbol:    "ZSWCC",
			},
			dbRow:       []byte("{\"supply\":\"1074146295.0178EOS\",\"max_supply\":\"10000000000.0000 EOS\",\"issuer\":\"zswhq\"}"),
			expectError: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			value, err := getTokenFromDBRow(test.contract, test.symbol, test.dbRow)

			if test.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, test.expectValue, value)
			}

		})
	}
}
