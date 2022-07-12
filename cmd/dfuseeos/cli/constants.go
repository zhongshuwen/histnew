package cli

import (
	pbbstream "github.com/streamingfast/pbgo/dfuse/bstream/v1"
)

const (
	Protocol               pbbstream.Protocol = pbbstream.Protocol_EOS
	TrxDBDSN               string             = "badger://{dfuse-data-dir}/storage/trxdb"
	StateDBDSN             string             = "badger://{dfuse-data-dir}/storage/statedb-v1"
	AccountHistDSN         string             = "badger://{dfuse-data-dir}/storage/accounthist"
	MergedBlocksStoreURL   string             = "file://{dfuse-data-dir}/storage/merged-blocks"
	FilteredBlocksStoreURL string             = "file://{dfuse-data-dir}/storage/filtered-merged-blocks"
	IndicesStoreURL        string             = "file://{dfuse-data-dir}/storage/indexes"
	OneBlockStoreURL       string             = "file://{dfuse-data-dir}/storage/one-blocks"
	PitreosURL             string             = "file://{dfuse-data-dir}/storage/pitreos"
	SnapshotsURL           string             = "file://{dfuse-data-dir}/storage/snapshots"
	DmeshDSN               string             = "local://"
	DmeshServiceVersion    string             = "v1"
	NetworkID              string             = "eos-local"
	NodeosBinPath          string             = "nodeos"
	// Ports
	NodeManagerHTTPServingAddr  string = ":13008"
	MindreaderHTTPServingAddr   string = ":13009"
	MindreaderGRPCAddr          string = ":13010"
	RelayerServingAddr          string = ":13011"
	MergerServingAddr           string = ":13012"
	ABICodecServingAddr         string = ":13013"
	BlockmetaServingAddr        string = ":13014"
	ArchiveServingAddr          string = ":13015"
	ArchiveHTTPServingAddr      string = ":13016"
	LiveServingAddr             string = ":13017"
	RouterServingAddr           string = ":13018"
	RouterHTTPServingAddr       string = ":13019"
	KvdbHTTPServingAddr         string = ":13020"
	IndexerServingAddr          string = ":13021"
	IndexerHTTPServingAddr      string = ":13022"
	DgraphqlHTTPServingAddr     string = ":13023"
	DgraphqlGRPCServingAddr     string = ":13024"
	EoswsHTTPServingAddr        string = ":13026"
	ForkResolverServingAddr     string = ":13027"
	ForkResolverHTTPServingAddr string = ":13028"
	StateDBHTTPServingAddr      string = ":13029"
	StateDBGRPCServingAddr      string = ":13032"
	EosqHTTPServingAddr         string = ":13030"
	DashboardGRPCServingAddr    string = ":13031"
	FilteringRelayerServingAddr string = ":13032"
	AccountHistGRPCServingAddr  string = ":13034"
	FirehoseGRPCServingAddr     string = ":13035"
	TokenmetaGRPCServingAddr    string = ":14001"
	DashboardHTTPListenAddr     string = ":8081"
	APIProxyHTTPListenAddr      string = ":8080"
	MindreaderNodeosAPIAddr     string = ":9888"
	NodeosAPIAddr               string = ":8888"
	MetricsListenAddr           string = ":9102"

	DgraphqlAPIKey string = "web_0000"
	JWTIssuerURL   string = "null://dfuse"
	EosqAPIKey     string = "web_0000"
)
