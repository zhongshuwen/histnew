// Copyright 2020 dfuse Platform Inc.
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

package tests

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"testing"

	_ "github.com/zhongshuwen/histnew/codec"
	_ "github.com/streamingfast/kvdb/store/badger"

	"github.com/streamingfast/bstream"
	ct "github.com/zhongshuwen/histnew/codec/testing"
	pbcodec "github.com/zhongshuwen/histnew/pb/dfuse/zswhq/codec/v1"
	"github.com/streamingfast/logging"
	"github.com/gavv/httpexpect/v2"
	"github.com/streamingfast/fluxdb/store"
	fluxdbKV "github.com/streamingfast/fluxdb/store/kv"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func init() {
	if os.Getenv("DEBUG") != "" {
		logger, _ := zap.NewDevelopment()
		logging.Override(logger)
	}
}

func TestAll(t *testing.T) {
	runAll(t, getKVTestFactory(t))
}

func getKVTestFactory(t *testing.T) func() (store.KVStore, StoreCleanupFunc) {
	return func() (store.KVStore, StoreCleanupFunc) {
		tmp, err := ioutil.TempDir("", "badger")
		require.NoError(t, err)
		kvStore, err := fluxdbKV.NewStore(fmt.Sprintf("badger://%s/test.db?createTables=true", tmp))
		require.NoError(t, err)

		closer := func() {
			kvStore.Close()
			os.RemoveAll(tmp)
		}

		return kvStore, closer
	}
}

func runAll(t *testing.T, storeFactory StoreFactory) {
	all := map[string][]e2eTester{
		"abi": {
			testStateABIHex,
		},
		"table": {
			testStateTableSingleRowHeadHex,
			testStateTableSingleRowHeadJSON,
			testStateTableSingleRowHistoricalJSON,
			testStateTableMultiRowsHeadJSON,
			testStateTableMultiRowsHistoricalJSON,
		},
		"table_scope": {
			testStateTableScopesHeadJSON,
			testStateTableScopesHistoricalJSON,
		},
		"tables_for_accounts": {
			testStateTablesForAccountsHeadJSON,
			testStateTablesForAccountsHistoricalJSON,
		},
		"table_row": {
			testStateTableRowHeadJSON,
		},
	}

	for group, tests := range all {
		for _, test := range tests {
			t.Run(group+"/"+getFunctionName(test), func(t *testing.T) {
				e2eTest(t, storeFactory, test)
			})
		}
	}
}

func testStateABIHex(ctx context.Context, t *testing.T, feedSourceWithBlocks blocksFeeder, e *httpexpect.Expect) {
	feedSourceWithBlocks(tableBlocks(t)...)

	response := okQueryStateABI(e, "zswhq.test", "")

	jsonValueEqual(t, "abi", `{"abi": "0e656f73696f3a3a6162692f312e3000010576616c7565000102746f0675696e743634000100000000008139bd0000000576616c756500000000", "account": "zswhq.test", "block_num": 5}`, response.Path("$"))
}

func testStateTableSingleRowHeadHex(ctx context.Context, t *testing.T, feedSourceWithBlocks blocksFeeder, e *httpexpect.Expect) {
	feedSourceWithBlocks(tableBlocks(t)...)

	response := okQueryStateTable(e, "zswhq.token/accounts/zswhq1", "")

	assertHeadBlockInfo(response, "00000006aa", "00000005aa")
	jsonValueEqual(t, "table-rows", `[{"key":"eos","payer":"zswhq1","hex":"a08601000000000004454f5300000000"}]`, response.Path("$.rows"))
}

func testStateTableSingleRowHeadJSON(ctx context.Context, t *testing.T, feedSourceWithBlocks blocksFeeder, e *httpexpect.Expect) {
	feedSourceWithBlocks(tableBlocks(t)...)

	response := okQueryStateTable(e, "zswhq.token/accounts/zswhq1", "json=true")

	assertHeadBlockInfo(response, "00000006aa", "00000005aa")
	jsonValueEqual(t, "table-rows", `[{"key":"eos","payer":"zswhq1","json":{"balance":"10.0000 EOS"}}]`, response.Path("$.rows"))
}

func testStateTableSingleRowHistoricalJSON(ctx context.Context, t *testing.T, feedSourceWithBlocks blocksFeeder, e *httpexpect.Expect) {
	feedSourceWithBlocks(tableBlocks(t)...)

	response := okQueryStateTable(e, "zswhq.token/accounts/zswhq1", "json=true&block_num=4")

	assertIrrBlockInfo(response, "00000005aa")
	jsonValueEqual(t, "table-rows", `[{"key":"eos","payer":"zswhq1","json":{"balance":"10.0000 EOS"}}]`, response.Path("$.rows"))
}

func testStateTableMultiRowsHeadJSON(ctx context.Context, t *testing.T, feedSourceWithBlocks blocksFeeder, e *httpexpect.Expect) {
	feedSourceWithBlocks(tableBlocks(t)...)

	response := okQueryStateTable(e, "zswhq.test/rows2/s", "json=true")

	assertHeadBlockInfo(response, "00000006aa", "00000005aa")
	jsonValueEqual(t, "table-rows", `[
		{"key":"b","payer":"s","json":{"to":20}},
		{"key":"c","payer":"s","json":{"to":3}},
		{"key":"d","payer":"s","json":{"to":4}},
		{"key":"e","payer":"s","json":{"to":5}},
		{"key":"f","payer":"s","json":{"to":6}}
	]`, response.Path("$.rows"))
}

func testStateTableMultiRowsHistoricalJSON(ctx context.Context, t *testing.T, feedSourceWithBlocks blocksFeeder, e *httpexpect.Expect) {
	feedSourceWithBlocks(tableBlocks(t)...)

	response := okQueryStateTable(e, "zswhq.test/rows/s", "json=true&block_num=3")

	assertIrrBlockInfo(response, "00000005aa")
	jsonValueEqual(t, "table-rows", `[
		{"key":"a","payer":"s","json":{"from":"a"}},
		{"key":"b","payer":"s","json":{"from":"b2"}},
		{"key":"c","payer":"s","json":{"from":"c"}}
	]`, response.Path("$.rows"))
}

func testStateTableScopesHeadJSON(ctx context.Context, t *testing.T, feedSourceWithBlocks blocksFeeder, e *httpexpect.Expect) {
	feedSourceWithBlocks(tableBlocks(t)...)

	response := okQueryStateTableScopes(e, "zswhq.token/accounts", "")

	response.ValueEqual("block_num", 6)
	jsonValueEqual(t, "scopes", `["zswhq1", "zswhq2"]`, response.Path("$.scopes"))
}

func testStateTableScopesHistoricalJSON(ctx context.Context, t *testing.T, feedSourceWithBlocks blocksFeeder, e *httpexpect.Expect) {
	feedSourceWithBlocks(tableBlocks(t)...)

	response := okQueryStateTableScopes(e, "zswhq.token/accounts", "block_num=3")

	response.ValueEqual("block_num", 3)
	jsonValueEqual(t, "scopes", `["zswhq1", "zswhq2", "zswhq3"]`, response.Path("$.scopes"))
}

func testStateTablesForScopesHeadJSON(ctx context.Context, t *testing.T, feedSourceWithBlocks blocksFeeder, e *httpexpect.Expect) {
	feedSourceWithBlocks(tableBlocks(t)...)

	response := okQueryStateTablesForScopes(e, "zswhq.token/accounts/zswhq1|zswhq2|zswhq3", "json=true")

	// That is not the correct behavior, there should be only `zswhq1` & `zswhq3` in the tests
	assertHeadBlockInfo(response, "00000006aa", "00000005aa")
	jsonValueEqual(t, "tables", `[
		{ "account": "zswhq.token","scope": "zswhq1", "rows": [{ "key": "eos", "payer": "zswhq1", "json": {"balance":"10.0000 EOS"}}]},
		{ "account": "zswhq.token","scope": "zswhq2", "rows": [{ "key": "eos", "payer": "zswhq2", "json": {"balance":"22.0000 EOS"}}]},
		{ "account": "zswhq.token","scope": "zswhq3", "rows": []}
	]`, response.Path("$.tables"))
}

func testStateTablesForScopesHistoricalJSON(ctx context.Context, t *testing.T, feedSourceWithBlocks blocksFeeder, e *httpexpect.Expect) {
	feedSourceWithBlocks(tableBlocks(t)...)

	response := okQueryStateTablesForScopes(e, "zswhq.token/accounts/zswhq1|zswhq2|zswhq3", "block_num=3&json=true")

	assertIrrBlockInfo(response, "00000005aa")
	jsonValueEqual(t, "tables", `[
		{ "account": "zswhq.token","scope": "zswhq1", "rows": [{ "key": "eos", "payer": "zswhq1", "json": {"balance":"1.0000 EOS"}}]},
		{ "account": "zswhq.token","scope": "zswhq2", "rows": [{ "key": "eos", "payer": "zswhq2", "json": {"balance":"20.0000 EOS"}}]},
		{ "account": "zswhq.token","scope": "zswhq3", "rows": [{ "key": "eos", "payer": "zswhq3", "json": {"balance":"3.0000 EOS"}}]}
	]`, response.Path("$.tables"))
}

func testStateTablesForAccountsHeadJSON(ctx context.Context, t *testing.T, feedSourceWithBlocks blocksFeeder, e *httpexpect.Expect) {
	feedSourceWithBlocks(tableBlocks(t)...)

	response := okQueryStateTablesForAccounts(e, "zswhq.token|zswhq.nekot/accounts/zswhq1", "json=true")

	assertHeadBlockInfo(response, "00000006aa", "00000005aa")
	jsonValueEqual(t, "tables", `[
		{"account":"zswhq.nekot","scope":"zswhq1","rows":[{"key":"eos","payer":"zswhq1","json":{"balance":"1.0000 SOE"}}]},
		{"account":"zswhq.token","scope":"zswhq1","rows":[{"key":"eos","payer":"zswhq1","json":{"balance":"10.0000 EOS"}}]}
	]`, response.Path("$.tables"))
}

func testStateTablesForAccountsHistoricalJSON(ctx context.Context, t *testing.T, feedSourceWithBlocks blocksFeeder, e *httpexpect.Expect) {
	feedSourceWithBlocks(tableBlocks(t)...)

	response := okQueryStateTablesForAccounts(e, "zswhq.token|zswhq.nekot/accounts/zswhq1", "block_num=4&json=true")

	assertIrrBlockInfo(response, "00000005aa")
	jsonValueEqual(t, "tables", `[
		{"account":"zswhq.nekot","scope":"zswhq1","rows":[{"key":"eos","payer":"zswhq1","json":{"balance":"1.0000 SOE"}}]},
		{"account":"zswhq.token","scope":"zswhq1","rows":[{"key":"eos","payer":"zswhq1","json":{"balance":"10.0000 EOS"}}]}
	]`, response.Path("$.tables"))
}

func testStateTableRowHeadJSON(ctx context.Context, t *testing.T, feedSourceWithBlocks blocksFeeder, e *httpexpect.Expect) {
	feedSourceWithBlocks(tableBlocks(t)...)

	response := okQueryStateTableRow(e, "zswhq.nekot/accounts/zswhq5/SOE", "json=true&key_type=symbol_code")

	assertHeadBlockInfo(response, "00000006aa", "00000005aa")
	jsonValueEqual(t, "row", `{"key":"SOE","payer":"zswhq5","json":{"balance":"5.0000 SOE"}}`, response.Path("$.row"))
}

func tableBlocks(t *testing.T) []*pbcodec.Block {
	zswhqTokenABI1 := readABI(t, "zswhq.token.1.abi.json")
	zswhqTestABI1 := readABI(t, "zswhq.test.1.abi.json")
	zswhqTestABI2 := readABI(t, "zswhq.test.2.abi.json")
	zswhqNekotABI1 := readABI(t, "zswhq.nekot.1.abi.json")

	return []*pbcodec.Block{
		// Block #2 | Sets ABI on `zswhq.token` (v1) and `zswhq.test` (v1)
		ct.Block(t, "00000002aa",
			ct.TrxTrace(t, ct.ActionTraceSetABI(t, "zswhq.token", zswhqTokenABI1)),
			ct.TrxTrace(t, ct.ActionTraceSetABI(t, "zswhq.test", zswhqTestABI1)),
		),

		// Block #3
		ct.Block(t, "00000003aa",
			// Creates three balances `zswhq1`, `zswhq2`, `zswhq3` on `zswhq.token`
			ct.TrxTrace(t,
				ct.TableOp(t, "insert", "zswhq.token/accounts/zswhq1", "zswhq1"),
				ct.DBOp(t, "insert", "zswhq.token/accounts/zswhq1/eos", "/zswhq1", `/{"balance":"1.0000 EOS"}`, zswhqTokenABI1),

				ct.TableOp(t, "insert", "zswhq.token/accounts/zswhq2", "zswhq2"),
				ct.DBOp(t, "insert", "zswhq.token/accounts/zswhq2/eos", "/zswhq2", `/{"balance":"2.0000 EOS"}`, zswhqTokenABI1),

				ct.TableOp(t, "insert", "zswhq.token/accounts/zswhq3", "zswhq3"),
				ct.DBOp(t, "insert", "zswhq.token/accounts/zswhq3/eos", "/zswhq3", `/{"balance":"3.0000 EOS"}`, zswhqTokenABI1),
			),

			// Add three rows (keys `a`, `b` & `c`) to `zswhq.test` contract, on table `rows` under scope `s`, then update key `b` within same transaction
			ct.TrxTrace(t,
				ct.TableOp(t, "insert", "zswhq.test/rows/s", "s"),
				ct.DBOp(t, "insert", "zswhq.test/rows/s/a", "/s", `/{"from":"a"}`, zswhqTestABI1),
				ct.DBOp(t, "insert", "zswhq.test/rows/s/b", "/s", `/{"from":"b"}`, zswhqTestABI1),
				ct.DBOp(t, "insert", "zswhq.test/rows/s/c", "/s", `/{"from":"c"}`, zswhqTestABI1),
				ct.DBOp(t, "update", "zswhq.test/rows/s/b", "s/s", `{"from":"b"}/{"from":"b2"}`, zswhqTestABI1),
			),

			// Update balance of `zswhq2` on `zswhq.token` within same block, but in different transaction
			ct.TrxTrace(t,
				ct.DBOp(t, "update", "zswhq.token/accounts/zswhq2/eos", "zswhq2/zswhq2", `{"balance":"2.0000 EOS"}/{"balance":"20.0000 EOS"}`, zswhqTokenABI1),
			),
		),

		// Block #4
		ct.Block(t, "00000004aa",
			// Add a new token contract `zswhq.nekot` (to test `/tables/accounts` calls) and populate odd rows from `zswhq.token`
			ct.TrxTrace(t,
				ct.ActionTraceSetABI(t, "zswhq.nekot", zswhqNekotABI1),

				ct.TableOp(t, "insert", "zswhq.nekot/accounts/zswhq1", "zswhq1"),
				ct.DBOp(t, "insert", "zswhq.nekot/accounts/zswhq1/eos", "/zswhq1", `/{"balance":"1.0000 SOE"}`, zswhqNekotABI1),

				ct.TableOp(t, "insert", "zswhq.nekot/accounts/zswhq3", "zswhq3"),
				ct.DBOp(t, "insert", "zswhq.nekot/accounts/zswhq3/eos", "/zswhq3", `/{"balance":"3.0000 SOE"}`, zswhqNekotABI1),
			),

			// Modify `zswhq.token` `zswhq1` balance and delete `zswhq3`
			ct.TrxTrace(t,
				ct.DBOp(t, "update", "zswhq.token/accounts/zswhq1/eos", "zswhq1/zswhq1", `{"balance":"1.0000 EOS"}/{"balance":"10.0000 EOS"}`, zswhqTokenABI1),

				ct.DBOp(t, "remove", "zswhq.token/accounts/zswhq3/eos", "zswhq3/", `{"balance":"3.0000 EOS"}/`, zswhqTokenABI1),
				ct.TableOp(t, "remove", "zswhq.token/accounts/zswhq3", "zswhq3"),
			),
		),

		// Block #5
		ct.Block(t, "00000005aa",
			// Remove all rows (keys `a`, `b`) of `zswhq.test`
			ct.TrxTrace(t,
				ct.DBOp(t, "remove", "zswhq.test/rows/s/a", "s/", `{"from":"a"}/`, zswhqTestABI1),
				ct.DBOp(t, "remove", "zswhq.test/rows/s/b", "s/", `{"from":"b2"}/`, zswhqTestABI1),
				ct.DBOp(t, "remove", "zswhq.test/rows/s/b", "s/", `{"from":"c"}/`, zswhqTestABI1),
				ct.TableOp(t, "remove", "zswhq.test/rows/s", "s"),
			),

			// Set a new ABI on `zswhq.test`
			ct.TrxTrace(t, ct.ActionTraceSetABI(t, "zswhq.test", zswhqTestABI2)),

			// Re-add all rows on `zswhq.test` using new ABI
			ct.TrxTrace(t,
				ct.TableOp(t, "insert", "zswhq.test/rows2/s", "s"),
				ct.DBOp(t, "insert", "zswhq.test/rows2/s/a", "/s", `/{"to":1}`, zswhqTestABI2),
				ct.DBOp(t, "insert", "zswhq.test/rows2/s/b", "/s", `/{"to":2}`, zswhqTestABI2),
				ct.DBOp(t, "insert", "zswhq.test/rows2/s/c", "/s", `/{"to":3}`, zswhqTestABI2),
			),

			// Add a new token contract `zswhq.nekot` (to test `/tables/accounts` calls) and populate odd rows from `zswhq.token`
			ct.TrxTrace(t,
				ct.TableOp(t, "insert", "zswhq.nekot/accounts/zswhq5", "zswhq5"),
				ct.DBOp(t, "insert", "zswhq.nekot/accounts/zswhq5/........cpbp3", "/zswhq5", `/{"balance":"5.0000 SOE"}`, zswhqNekotABI1),
			),
		),

		// Block #6 | This block will be in the reversible segment, i.e. in the speculative writes
		ct.Block(t, "00000006aa",
			// Update balance of `zswhq2` on `zswhq.token`
			ct.TrxTrace(t,
				ct.DBOp(t, "update", "zswhq.token/accounts/zswhq2/eos", "zswhq2/zswhq2", `{"balance":"20.0000 EOS"}/{"balance":"22.0000 EOS"}`, zswhqTokenABI1),
			),

			// Delete rows `a` from `zswhq.test`, update `b` and add three new rows (`d`, `e` & `f`)
			ct.TrxTrace(t,
				ct.DBOp(t, "remove", "zswhq.test/rows2/s/a", "s/", `{"to":1}/`, zswhqTestABI2),

				ct.DBOp(t, "update", "zswhq.test/rows2/s/b", "s/s", `{"to":2}/{"to":20}`, zswhqTestABI2),

				ct.DBOp(t, "insert", "zswhq.test/rows2/s/d", "/s", `/{"to":4}`, zswhqTestABI2),
				ct.DBOp(t, "insert", "zswhq.test/rows2/s/e", "/s", `/{"to":5}`, zswhqTestABI2),
				ct.DBOp(t, "insert", "zswhq.test/rows2/s/f", "/s", `/{"to":6}`, zswhqTestABI2),
			),
		),
	}
}

func okQueryStateABI(e *httpexpect.Expect, account string, extraQuery string) (response *httpexpect.Object) {
	queryString := fmt.Sprintf("account=%s", account)
	if extraQuery != "" {
		queryString += "" + extraQuery
	}

	return okQuery(e, "/v0/state/abi", queryString)
}

func okQueryStateTable(e *httpexpect.Expect, table string, extraQuery string) (response *httpexpect.Object) {
	parts := strings.Split(table, "/")

	queryString := fmt.Sprintf("account=%s&table=%s&scope=%s", parts[0], parts[1], parts[2])
	if extraQuery != "" {
		queryString += "&" + extraQuery
	}

	return okQuery(e, "/v0/state/table", queryString)
}

func okQueryStateTableScopes(e *httpexpect.Expect, table string, extraQuery string) (response *httpexpect.Object) {
	parts := strings.Split(table, "/")

	queryString := fmt.Sprintf("account=%s&table=%s", parts[0], parts[1])
	if extraQuery != "" {
		queryString += "&" + extraQuery
	}

	return okQuery(e, "/v0/state/table_scopes", queryString)
}

func okQueryStateTablesForScopes(e *httpexpect.Expect, table string, extraQuery string) (response *httpexpect.Object) {
	parts := strings.Split(table, "/")

	queryString := fmt.Sprintf("account=%s&table=%s&scopes=%s", parts[0], parts[1], parts[2])
	if extraQuery != "" {
		queryString += "&" + extraQuery
	}

	return okQuery(e, "/v0/state/tables/scopes", queryString)
}

func okQueryStateTablesForAccounts(e *httpexpect.Expect, table string, extraQuery string) (response *httpexpect.Object) {
	parts := strings.Split(table, "/")

	queryString := fmt.Sprintf("accounts=%s&table=%s&scope=%s", parts[0], parts[1], parts[2])
	if extraQuery != "" {
		queryString += "&" + extraQuery
	}

	return okQuery(e, "/v0/state/tables/accounts", queryString)
}

func okQueryStateTableRow(e *httpexpect.Expect, table string, extraQuery string) (response *httpexpect.Object) {
	parts := strings.Split(table, "/")

	queryString := fmt.Sprintf("account=%s&table=%s&scope=%s&primary_key=%s", parts[0], parts[1], parts[2], parts[3])
	if extraQuery != "" {
		queryString += "&" + extraQuery
	}

	return okQuery(e, "/v0/state/table/row", queryString)
}

func okQuery(e *httpexpect.Expect, path string, queryString string) (response *httpexpect.Object) {
	return e.GET(path).
		WithQueryString(queryString).
		Expect().
		Status(http.StatusOK).JSON().Object()
}

func assertIrrBlockInfo(response *httpexpect.Object, libRef string) {
	lRef := bstream.NewBlockRefFromID(libRef)

	response.NotContainsKey("up_to_block_id")
	response.NotContainsKey("up_to_block_num")

	response.ValueEqual("last_irreversible_block_id", lRef.ID())
	response.ValueEqual("last_irreversible_block_num", lRef.Num())
}

func assertHeadBlockInfo(response *httpexpect.Object, blockRef string, libRef string) {
	bRef := bstream.NewBlockRefFromID(blockRef)
	lRef := bstream.NewBlockRefFromID(libRef)

	response.ValueEqual("up_to_block_id", bRef.ID())
	response.ValueEqual("up_to_block_num", bRef.Num())

	response.ValueEqual("last_irreversible_block_id", lRef.ID())
	response.ValueEqual("last_irreversible_block_num", lRef.Num())
}

// getFunctionName reads the program counter adddress and return the function
// name representing this address.
//
// The `FuncForPC` format is in the form of `github.com/.../.../package.func`.
// As such, we use `filepath.Base` to obtain the `package.func` part and then
// split it at the `.` to extract the function name.
func getFunctionName(i interface{}) string {
	pcIdentifier := runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
	baseName := filepath.Base(pcIdentifier)
	parts := strings.SplitN(baseName, ".", 2)
	if len(parts) <= 1 {
		return parts[0]
	}

	return parts[1]
}
