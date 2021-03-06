package cache

import (
	"testing"

	"github.com/zhongshuwen/zswchain-go"
	"github.com/stretchr/testify/assert"
)

func Test_SortOwnedAssetSymbolAlpha(t *testing.T) {
	tests := []struct {
		name         string
		assets       []*OwnedAsset
		order        SortingOrder
		expectAssets []*OwnedAsset
	}{
		{
			name:  "ascending sort",
			order: ASC,
			assets: []*OwnedAsset{
				{Owner: "eoscanadadad", Asset: &zsw.ExtendedAsset{Asset: generateTestAsset(100, "ZSWCC")}},
				{Owner: "oiloiloiloil", Asset: &zsw.ExtendedAsset{Asset: generateTestAsset(3, "EIDOS")}},
				{Owner: "johndoeoneco", Asset: &zsw.ExtendedAsset{Asset: generateTestAsset(53, "ZSWCC")}},
				{Owner: "eoscanadacom", Asset: &zsw.ExtendedAsset{Asset: generateTestAsset(123, "WALL")}},
				{Owner: "cbocbocbocbo", Asset: &zsw.ExtendedAsset{Asset: generateTestAsset(293, "EIDOS")}},
				{Owner: "kolkolkolkol", Asset: &zsw.ExtendedAsset{Asset: generateTestAsset(23, "WALL")}},
			},
			expectAssets: []*OwnedAsset{
				{Owner: "cbocbocbocbo", Asset: &zsw.ExtendedAsset{Asset: generateTestAsset(293, "EIDOS")}},
				{Owner: "oiloiloiloil", Asset: &zsw.ExtendedAsset{Asset: generateTestAsset(3, "EIDOS")}},
				{Owner: "eoscanadadad", Asset: &zsw.ExtendedAsset{Asset: generateTestAsset(100, "ZSWCC")}},
				{Owner: "johndoeoneco", Asset: &zsw.ExtendedAsset{Asset: generateTestAsset(53, "ZSWCC")}},
				{Owner: "eoscanadacom", Asset: &zsw.ExtendedAsset{Asset: generateTestAsset(123, "WALL")}},
				{Owner: "kolkolkolkol", Asset: &zsw.ExtendedAsset{Asset: generateTestAsset(23, "WALL")}},
			},
		},
		{
			name:  "descending sort",
			order: DESC,
			assets: []*OwnedAsset{
				{Owner: "eoscanadadad", Asset: &zsw.ExtendedAsset{Asset: generateTestAsset(100, "ZSWCC")}},
				{Owner: "oiloiloiloil", Asset: &zsw.ExtendedAsset{Asset: generateTestAsset(3, "EIDOS")}},
				{Owner: "johndoeoneco", Asset: &zsw.ExtendedAsset{Asset: generateTestAsset(53, "ZSWCC")}},
				{Owner: "eoscanadacom", Asset: &zsw.ExtendedAsset{Asset: generateTestAsset(123, "WALL")}},
				{Owner: "cbocbocbocbo", Asset: &zsw.ExtendedAsset{Asset: generateTestAsset(293, "EIDOS")}},
				{Owner: "kolkolkolkol", Asset: &zsw.ExtendedAsset{Asset: generateTestAsset(23, "WALL")}},
			},
			expectAssets: []*OwnedAsset{
				{Owner: "eoscanadacom", Asset: &zsw.ExtendedAsset{Asset: generateTestAsset(123, "WALL")}},
				{Owner: "kolkolkolkol", Asset: &zsw.ExtendedAsset{Asset: generateTestAsset(23, "WALL")}},
				{Owner: "eoscanadadad", Asset: &zsw.ExtendedAsset{Asset: generateTestAsset(100, "ZSWCC")}},
				{Owner: "johndoeoneco", Asset: &zsw.ExtendedAsset{Asset: generateTestAsset(53, "ZSWCC")}},
				{Owner: "cbocbocbocbo", Asset: &zsw.ExtendedAsset{Asset: generateTestAsset(293, "EIDOS")}},
				{Owner: "oiloiloiloil", Asset: &zsw.ExtendedAsset{Asset: generateTestAsset(3, "EIDOS")}},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.expectAssets, SortOwnedAssetBySymbolAlpha(test.assets, test.order))
		})
	}
}

func Test_SortOwnedAssetTokenAmount(t *testing.T) {
	tests := []struct {
		name         string
		assets       []*OwnedAsset
		order        SortingOrder
		expectAssets []*OwnedAsset
	}{
		{
			name:  "ascending sort",
			order: ASC,
			assets: []*OwnedAsset{
				{Owner: "eoscanadadad", Asset: &zsw.ExtendedAsset{Asset: generateTestAsset(100, "ZSWCC")}},
				{Owner: "oiloiloiloil", Asset: &zsw.ExtendedAsset{Asset: generateTestAsset(3, "EIDOS")}},
				{Owner: "johndoeoneco", Asset: &zsw.ExtendedAsset{Asset: generateTestAsset(53, "ZSWCC")}},
				{Owner: "eoscanadacom", Asset: &zsw.ExtendedAsset{Asset: generateTestAsset(123, "WALL")}},
				{Owner: "cbocbocbocbo", Asset: &zsw.ExtendedAsset{Asset: generateTestAsset(293, "EIDOS")}},
				{Owner: "kolkolkolkol", Asset: &zsw.ExtendedAsset{Asset: generateTestAsset(23, "WALL")}},
			},
			expectAssets: []*OwnedAsset{
				{Owner: "oiloiloiloil", Asset: &zsw.ExtendedAsset{Asset: generateTestAsset(3, "EIDOS")}},
				{Owner: "kolkolkolkol", Asset: &zsw.ExtendedAsset{Asset: generateTestAsset(23, "WALL")}},
				{Owner: "johndoeoneco", Asset: &zsw.ExtendedAsset{Asset: generateTestAsset(53, "ZSWCC")}},
				{Owner: "eoscanadadad", Asset: &zsw.ExtendedAsset{Asset: generateTestAsset(100, "ZSWCC")}},
				{Owner: "eoscanadacom", Asset: &zsw.ExtendedAsset{Asset: generateTestAsset(123, "WALL")}},
				{Owner: "cbocbocbocbo", Asset: &zsw.ExtendedAsset{Asset: generateTestAsset(293, "EIDOS")}},
			},
		},
		{
			name:  "descending sort",
			order: DESC,
			assets: []*OwnedAsset{
				{Owner: "eoscanadadad", Asset: &zsw.ExtendedAsset{Asset: generateTestAsset(100, "ZSWCC")}},
				{Owner: "oiloiloiloil", Asset: &zsw.ExtendedAsset{Asset: generateTestAsset(3, "EIDOS")}},
				{Owner: "johndoeoneco", Asset: &zsw.ExtendedAsset{Asset: generateTestAsset(53, "ZSWCC")}},
				{Owner: "eoscanadacom", Asset: &zsw.ExtendedAsset{Asset: generateTestAsset(123, "WALL")}},
				{Owner: "cbocbocbocbo", Asset: &zsw.ExtendedAsset{Asset: generateTestAsset(293, "EIDOS")}},
				{Owner: "kolkolkolkol", Asset: &zsw.ExtendedAsset{Asset: generateTestAsset(23, "WALL")}},
			},
			expectAssets: []*OwnedAsset{
				{Owner: "cbocbocbocbo", Asset: &zsw.ExtendedAsset{Asset: generateTestAsset(293, "EIDOS")}},
				{Owner: "eoscanadacom", Asset: &zsw.ExtendedAsset{Asset: generateTestAsset(123, "WALL")}},
				{Owner: "eoscanadadad", Asset: &zsw.ExtendedAsset{Asset: generateTestAsset(100, "ZSWCC")}},
				{Owner: "johndoeoneco", Asset: &zsw.ExtendedAsset{Asset: generateTestAsset(53, "ZSWCC")}},
				{Owner: "kolkolkolkol", Asset: &zsw.ExtendedAsset{Asset: generateTestAsset(23, "WALL")}},
				{Owner: "oiloiloiloil", Asset: &zsw.ExtendedAsset{Asset: generateTestAsset(3, "EIDOS")}},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.expectAssets, SortOwnedAssetByTokenAmount(test.assets, test.order))
		})
	}
}

func Test_SortOwnedAssetAccountAlpha(t *testing.T) {
	tests := []struct {
		name         string
		assets       []*OwnedAsset
		order        SortingOrder
		expectAssets []*OwnedAsset
	}{
		{
			name:  "ascending sort",
			order: ASC,
			assets: []*OwnedAsset{
				{Owner: "eoscanadadad", Asset: &zsw.ExtendedAsset{Asset: generateTestAsset(100, "ZSWCC")}},
				{Owner: "oiloiloiloil", Asset: &zsw.ExtendedAsset{Asset: generateTestAsset(3, "EIDOS")}},
				{Owner: "johndoeoneco", Asset: &zsw.ExtendedAsset{Asset: generateTestAsset(53, "ZSWCC")}},
				{Owner: "eoscanadacom", Asset: &zsw.ExtendedAsset{Asset: generateTestAsset(123, "WALL")}},
				{Owner: "cbocbocbocbo", Asset: &zsw.ExtendedAsset{Asset: generateTestAsset(293, "EIDOS")}},
				{Owner: "kolkolkolkol", Asset: &zsw.ExtendedAsset{Asset: generateTestAsset(23, "WALL")}},
			},
			expectAssets: []*OwnedAsset{
				{Owner: "cbocbocbocbo", Asset: &zsw.ExtendedAsset{Asset: generateTestAsset(293, "EIDOS")}},
				{Owner: "eoscanadacom", Asset: &zsw.ExtendedAsset{Asset: generateTestAsset(123, "WALL")}},
				{Owner: "eoscanadadad", Asset: &zsw.ExtendedAsset{Asset: generateTestAsset(100, "ZSWCC")}},
				{Owner: "johndoeoneco", Asset: &zsw.ExtendedAsset{Asset: generateTestAsset(53, "ZSWCC")}},
				{Owner: "kolkolkolkol", Asset: &zsw.ExtendedAsset{Asset: generateTestAsset(23, "WALL")}},
				{Owner: "oiloiloiloil", Asset: &zsw.ExtendedAsset{Asset: generateTestAsset(3, "EIDOS")}},
			},
		},
		{
			name:  "descending sort",
			order: DESC,
			assets: []*OwnedAsset{
				{Owner: "eoscanadadad", Asset: &zsw.ExtendedAsset{Asset: generateTestAsset(100, "ZSWCC")}},
				{Owner: "oiloiloiloil", Asset: &zsw.ExtendedAsset{Asset: generateTestAsset(3, "EIDOS")}},
				{Owner: "johndoeoneco", Asset: &zsw.ExtendedAsset{Asset: generateTestAsset(53, "ZSWCC")}},
				{Owner: "eoscanadacom", Asset: &zsw.ExtendedAsset{Asset: generateTestAsset(123, "WALL")}},
				{Owner: "cbocbocbocbo", Asset: &zsw.ExtendedAsset{Asset: generateTestAsset(293, "EIDOS")}},
				{Owner: "kolkolkolkol", Asset: &zsw.ExtendedAsset{Asset: generateTestAsset(23, "WALL")}},
			},
			expectAssets: []*OwnedAsset{
				{Owner: "oiloiloiloil", Asset: &zsw.ExtendedAsset{Asset: generateTestAsset(3, "EIDOS")}},
				{Owner: "kolkolkolkol", Asset: &zsw.ExtendedAsset{Asset: generateTestAsset(23, "WALL")}},
				{Owner: "johndoeoneco", Asset: &zsw.ExtendedAsset{Asset: generateTestAsset(53, "ZSWCC")}},
				{Owner: "eoscanadadad", Asset: &zsw.ExtendedAsset{Asset: generateTestAsset(100, "ZSWCC")}},
				{Owner: "eoscanadacom", Asset: &zsw.ExtendedAsset{Asset: generateTestAsset(123, "WALL")}},
				{Owner: "cbocbocbocbo", Asset: &zsw.ExtendedAsset{Asset: generateTestAsset(293, "EIDOS")}},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.expectAssets, SortOwnedAssetByAccountAlpha(test.assets, test.order))
		})
	}
}

func Test_SortOwnedAssetTokenMarketValue(t *testing.T) {
	//TODO: implement me
}

func Test_compareOwnedAssetBySymbol(t *testing.T) {
	tests := []struct {
		name        string
		a           *OwnedAsset
		b           *OwnedAsset
		order       SortingOrder
		expectValue bool
	}{
		{
			name:        "ascending a before b",
			order:       ASC,
			a:           &OwnedAsset{Owner: "eoscanadadad", Asset: &zsw.ExtendedAsset{Contract: "zswhq.token", Asset: generateTestAsset(100, "ZSWCC")}},
			b:           &OwnedAsset{Owner: "eoscanadadad", Asset: &zsw.ExtendedAsset{Contract: "b1tokenmeta", Asset: generateTestAsset(100, "WAX")}},
			expectValue: true,
		},
		{
			name:        "ascending a after b",
			order:       ASC,
			a:           &OwnedAsset{Owner: "eoscanadadad", Asset: &zsw.ExtendedAsset{Contract: "zswhq.token", Asset: generateTestAsset(100, "WAX")}},
			b:           &OwnedAsset{Owner: "eoscanadadad", Asset: &zsw.ExtendedAsset{Contract: "b1tokenmeta", Asset: generateTestAsset(100, "ZSWCC")}},
			expectValue: false,
		},
		{
			name:        "ascending a equal b, should sort by contract",
			order:       ASC,
			a:           &OwnedAsset{Owner: "eoscanadadad", Asset: &zsw.ExtendedAsset{Contract: "zswhq.token", Asset: generateTestAsset(100, "ZSWCC")}},
			b:           &OwnedAsset{Owner: "eoscanadadad", Asset: &zsw.ExtendedAsset{Contract: "b1tokenmeta", Asset: generateTestAsset(100, "ZSWCC")}},
			expectValue: false,
		},
		{
			name:        "ascending a equal b, should sort by owner if contract is equal",
			order:       ASC,
			a:           &OwnedAsset{Owner: "eosaaaaaaaaa", Asset: &zsw.ExtendedAsset{Contract: "zswhq.token", Asset: generateTestAsset(100, "ZSWCC")}},
			b:           &OwnedAsset{Owner: "eoscanadadad", Asset: &zsw.ExtendedAsset{Contract: "zswhq.token", Asset: generateTestAsset(100, "ZSWCC")}},
			expectValue: true,
		},
		{
			name:        "descending a before b",
			order:       DESC,
			a:           &OwnedAsset{Owner: "eoscanadadad", Asset: &zsw.ExtendedAsset{Contract: "zswhq.token", Asset: generateTestAsset(100, "WAX")}},
			b:           &OwnedAsset{Owner: "eoscanadadad", Asset: &zsw.ExtendedAsset{Contract: "b1tokenmeta", Asset: generateTestAsset(100, "ZSWCC")}},
			expectValue: true,
		},
		{
			name:        "descending a after b",
			order:       DESC,
			a:           &OwnedAsset{Owner: "eoscanadadad", Asset: &zsw.ExtendedAsset{Contract: "zswhq.token", Asset: generateTestAsset(100, "ZSWCC")}},
			b:           &OwnedAsset{Owner: "eoscanadadad", Asset: &zsw.ExtendedAsset{Contract: "b1tokenmeta", Asset: generateTestAsset(100, "WAX")}},
			expectValue: false,
		},
		{
			name:        "descending a equal b, should sort by contract",
			order:       DESC,
			a:           &OwnedAsset{Owner: "eoscanadadad", Asset: &zsw.ExtendedAsset{Contract: "b1tokenmeta", Asset: generateTestAsset(100, "ZSWCC")}},
			b:           &OwnedAsset{Owner: "eoscanadadad", Asset: &zsw.ExtendedAsset{Contract: "zswhq.token", Asset: generateTestAsset(100, "ZSWCC")}},
			expectValue: true,
		},
		{
			name:        "descending a equal b, should sort by owner if contract is equal",
			order:       DESC,
			a:           &OwnedAsset{Owner: "eoscanadadad", Asset: &zsw.ExtendedAsset{Contract: "zswhq.token", Asset: generateTestAsset(100, "ZSWCC")}},
			b:           &OwnedAsset{Owner: "eosaaaaaaaaa", Asset: &zsw.ExtendedAsset{Contract: "zswhq.token", Asset: generateTestAsset(100, "ZSWCC")}},
			expectValue: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.expectValue, compareOwnedAssetBySymbol(test.a, test.b, test.order))
		})
	}
}

func Test_compareOwnedAssetByAccount(t *testing.T) {
	tests := []struct {
		name        string
		a           *OwnedAsset
		b           *OwnedAsset
		order       SortingOrder
		expectValue bool
	}{
		{
			name:        "ascending a before b",
			order:       ASC,
			a:           &OwnedAsset{Owner: "aaaaaaaaaaaa", Asset: &zsw.ExtendedAsset{Contract: "zswhq.token", Asset: generateTestAsset(100, "ZSWCC")}},
			b:           &OwnedAsset{Owner: "bbbbbbbbbbbb", Asset: &zsw.ExtendedAsset{Contract: "zswhq.token", Asset: generateTestAsset(23, "ZSWCC")}},
			expectValue: true,
		},
		{
			name:        "ascending a after b",
			order:       ASC,
			a:           &OwnedAsset{Owner: "bbbbbbbbbbbb", Asset: &zsw.ExtendedAsset{Contract: "zswhq.token", Asset: generateTestAsset(100, "ZSWCC")}},
			b:           &OwnedAsset{Owner: "aaaaaaaaaaaa", Asset: &zsw.ExtendedAsset{Contract: "zswhq.token", Asset: generateTestAsset(23, "ZSWCC")}},
			expectValue: false,
		},
		{
			name:        "ascending a equal b, should sort by symbol",
			order:       ASC,
			a:           &OwnedAsset{Owner: "aaaaaaaaaaaa", Asset: &zsw.ExtendedAsset{Contract: "zswhq.token", Asset: generateTestAsset(100, "ZSWCC")}},
			b:           &OwnedAsset{Owner: "aaaaaaaaaaaa", Asset: &zsw.ExtendedAsset{Contract: "zswhq.token", Asset: generateTestAsset(23, "WAX")}},
			expectValue: true,
		},
		{
			name:        "ascending a equal b, should sort by contract if symbol is equal",
			order:       ASC,
			a:           &OwnedAsset{Owner: "aaaaaaaaaaaa", Asset: &zsw.ExtendedAsset{Contract: "zswhq.token", Asset: generateTestAsset(100, "ZSWCC")}},
			b:           &OwnedAsset{Owner: "aaaaaaaaaaaa", Asset: &zsw.ExtendedAsset{Contract: "b1tokenmeta", Asset: generateTestAsset(23, "ZSWCC")}},
			expectValue: false,
		},
		{
			name:        "descending a before b",
			order:       DESC,
			a:           &OwnedAsset{Owner: "bbbbbbbbbbbb", Asset: &zsw.ExtendedAsset{Contract: "zswhq.token", Asset: generateTestAsset(100, "ZSWCC")}},
			b:           &OwnedAsset{Owner: "aaaaaaaaaaaa", Asset: &zsw.ExtendedAsset{Contract: "zswhq.token", Asset: generateTestAsset(23, "ZSWCC")}},
			expectValue: true,
		},
		{
			name:        "descending a after b",
			order:       DESC,
			a:           &OwnedAsset{Owner: "aaaaaaaaaaaa", Asset: &zsw.ExtendedAsset{Contract: "zswhq.token", Asset: generateTestAsset(100, "ZSWCC")}},
			b:           &OwnedAsset{Owner: "bbbbbbbbbbbb", Asset: &zsw.ExtendedAsset{Contract: "zswhq.token", Asset: generateTestAsset(23, "ZSWCC")}},
			expectValue: false,
		},
		{
			name:        "descending a equal b, should sort by symbol",
			order:       DESC,
			a:           &OwnedAsset{Owner: "bbbbbbbbbbbb", Asset: &zsw.ExtendedAsset{Contract: "zswhq.token", Asset: generateTestAsset(100, "WAX")}},
			b:           &OwnedAsset{Owner: "bbbbbbbbbbbb", Asset: &zsw.ExtendedAsset{Contract: "zswhq.token", Asset: generateTestAsset(23, "ZSWCC")}},
			expectValue: false,
		},
		{
			name:        "descending a equal b, should sort by owner if contract is equal",
			order:       DESC,
			a:           &OwnedAsset{Owner: "bbbbbbbbbbbb", Asset: &zsw.ExtendedAsset{Contract: "zswhq.token", Asset: generateTestAsset(100, "ZSWCC")}},
			b:           &OwnedAsset{Owner: "bbbbbbbbbbbb", Asset: &zsw.ExtendedAsset{Contract: "b1tokenmeta", Asset: generateTestAsset(23, "ZSWCC")}},
			expectValue: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.expectValue, compareOwnedAssetByAccount(test.a, test.b, test.order))
		})
	}
}

func Test_compareOwnedAssetByTokenAmount(t *testing.T) {
	tests := []struct {
		name        string
		a           *OwnedAsset
		b           *OwnedAsset
		order       SortingOrder
		expectValue bool
	}{
		{
			name:        "ascending a before b",
			order:       ASC,
			a:           &OwnedAsset{Owner: "aaaaaaaaaaaa", Asset: &zsw.ExtendedAsset{Contract: "zswhq.token", Asset: generateTestAsset(32, "ZSWCC")}},
			b:           &OwnedAsset{Owner: "bbbbbbbbbbbb", Asset: &zsw.ExtendedAsset{Contract: "zswhq.token", Asset: generateTestAsset(84, "ZSWCC")}},
			expectValue: true,
		},
		{
			name:        "ascending a after b",
			order:       ASC,
			a:           &OwnedAsset{Owner: "bbbbbbbbbbbb", Asset: &zsw.ExtendedAsset{Contract: "zswhq.token", Asset: generateTestAsset(100, "ZSWCC")}},
			b:           &OwnedAsset{Owner: "aaaaaaaaaaaa", Asset: &zsw.ExtendedAsset{Contract: "zswhq.token", Asset: generateTestAsset(32, "ZSWCC")}},
			expectValue: false,
		},
		{
			name:        "ascending a equal b, should sort by accounting",
			order:       ASC,
			a:           &OwnedAsset{Owner: "aaaaaaaaaaaa", Asset: &zsw.ExtendedAsset{Contract: "zswhq.token", Asset: generateTestAsset(23, "ZSWCC")}},
			b:           &OwnedAsset{Owner: "bbbbbbbbbbbb", Asset: &zsw.ExtendedAsset{Contract: "zswhq.token", Asset: generateTestAsset(23, "ZSWCC")}},
			expectValue: true,
		},
		{
			name:        "ascending a equal b, should sort by contract if symbol is equal",
			order:       ASC,
			a:           &OwnedAsset{Owner: "aaaaaaaaaaaa", Asset: &zsw.ExtendedAsset{Contract: "zswhq.token", Asset: generateTestAsset(23, "ZSWCC")}},
			b:           &OwnedAsset{Owner: "aaaaaaaaaaaa", Asset: &zsw.ExtendedAsset{Contract: "b1tokenmeta", Asset: generateTestAsset(23, "ZSWCC")}},
			expectValue: false,
		},
		{
			name:        "descending a before b",
			order:       DESC,
			a:           &OwnedAsset{Owner: "bbbbbbbbbbbb", Asset: &zsw.ExtendedAsset{Contract: "zswhq.token", Asset: generateTestAsset(100, "ZSWCC")}},
			b:           &OwnedAsset{Owner: "aaaaaaaaaaaa", Asset: &zsw.ExtendedAsset{Contract: "zswhq.token", Asset: generateTestAsset(23, "ZSWCC")}},
			expectValue: true,
		},
		{
			name:        "descending a after b",
			order:       DESC,
			a:           &OwnedAsset{Owner: "aaaaaaaaaaaa", Asset: &zsw.ExtendedAsset{Contract: "zswhq.token", Asset: generateTestAsset(23, "ZSWCC")}},
			b:           &OwnedAsset{Owner: "bbbbbbbbbbbb", Asset: &zsw.ExtendedAsset{Contract: "zswhq.token", Asset: generateTestAsset(100, "ZSWCC")}},
			expectValue: false,
		},
		{
			name:        "descending a equal b, should sort by symbol",
			order:       DESC,
			a:           &OwnedAsset{Owner: "aaaaaaaaaaaa", Asset: &zsw.ExtendedAsset{Contract: "zswhq.token", Asset: generateTestAsset(100, "ZSWCC")}},
			b:           &OwnedAsset{Owner: "bbbbbbbbbbbb", Asset: &zsw.ExtendedAsset{Contract: "zswhq.token", Asset: generateTestAsset(100, "ZSWCC")}},
			expectValue: true,
		},
		{
			name:        "descending a equal b, should sort by owner if contract is equal",
			order:       DESC,
			a:           &OwnedAsset{Owner: "bbbbbbbbbbbb", Asset: &zsw.ExtendedAsset{Contract: "zswhq.token", Asset: generateTestAsset(23, "ZSWCC")}},
			b:           &OwnedAsset{Owner: "bbbbbbbbbbbb", Asset: &zsw.ExtendedAsset{Contract: "b1tokenmeta", Asset: generateTestAsset(23, "ZSWCC")}},
			expectValue: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.expectValue, compareOwnedAssetByTokenAmount(test.a, test.b, test.order))
		})
	}
}
