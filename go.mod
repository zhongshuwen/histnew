module github.com/zhongshuwen/histnew

go 1.16

require (
	contrib.go.opencensus.io/exporter/stackdriver v0.13.10
	github.com/GeertJohan/go.rice v1.0.2
	github.com/ShinyTrinkets/overseer v0.3.0
	github.com/acarl005/stripansi v0.0.0-20180116102854-5a71ef0e047d
	github.com/andreyvit/diff v0.0.0-20170406064948-c7f18ee00883
	github.com/araddon/dateparse v0.0.0-20190622164848-0fb0a474d195
	github.com/arpitbbhayani/tripod v0.0.0-20170425181942-66807adce3a5
	github.com/auth0/go-jwt-middleware v0.0.0-20190805220309-36081240882b
	github.com/blevesearch/bleve v1.0.14
	github.com/bmizerany/assert v0.0.0-20160611221934-b7ed37b82869 // indirect
	github.com/coreos/etcd v3.3.25+incompatible // indirect
	github.com/coreos/go-systemd v0.0.0-20191104093116-d3cd4ed1dbcf // indirect
	github.com/daaku/go.zipexe v1.0.1 // indirect
	github.com/davecgh/go-spew v1.1.1
	github.com/dustin/go-humanize v1.0.0
	github.com/eoscanada/eos-go v0.9.1-0.20210812015252-984fc96878b6 // indirect
	github.com/eoscanada/eosc v1.4.0
	github.com/eoscanada/pitreos v1.1.1-0.20210811185752-fa06394508d0
	github.com/francoispqt/gojay v1.2.13
	github.com/gavv/httpexpect/v2 v2.0.3
	github.com/golang-collections/collections v0.0.0-20130729185459-604e922904d3
	github.com/golang/protobuf v1.5.2
	github.com/google/cel-go v0.4.1
	github.com/gorilla/handlers v1.4.2
	github.com/gorilla/mux v1.8.0
	github.com/gorilla/websocket v1.4.2
	github.com/graph-gophers/graphql-go v0.0.0-20191115155744-f33e81362277
	github.com/invisible-train-40/client-go v0.1.2
	github.com/invisible-train-40/eosio-boot v0.1.0
	github.com/invisible-train-40/eosws-go v0.1.0
	github.com/klauspost/compress v1.10.2
	github.com/lithammer/dedent v1.1.0
	github.com/logrusorgru/aurora v2.0.3+incompatible
	github.com/lytics/lifecycle v0.0.0-20130117214539-7b4c4028d422 // indirect
	github.com/lytics/ordpool v0.0.0-20130426221837-8d833f097fe7
	github.com/manifoldco/promptui v0.8.0
	github.com/mitchellh/go-testing-interface v1.14.1
	github.com/paulbellamy/ratecounter v0.2.0
	github.com/pkg/errors v0.9.1
	github.com/sergi/go-diff v1.0.1-0.20180205163309-da645544ed44 // indirect
	github.com/sirupsen/logrus v1.8.1 // indirect
	github.com/spf13/cobra v1.1.3
	github.com/spf13/viper v1.7.0
	github.com/streamingfast/blockmeta v0.0.2-0.20210811194956-90dc4202afda
	github.com/streamingfast/bstream v0.0.2-0.20211006133726-b4a4315e6708
	github.com/streamingfast/cli v0.0.3-0.20210811201236-5c00ec55462d
	github.com/streamingfast/dauth v0.0.0-20210812020920-1c83ba29add1
	github.com/streamingfast/dbin v0.0.0-20210809205249-73d5eca35dc5
	github.com/streamingfast/derr v0.0.0-20220301163149-de09cb18fc70
	github.com/streamingfast/dgraphql v0.0.2-0.20220307143518-466192441cfe
	github.com/streamingfast/dgrpc v0.0.0-20220301153539-536adf71b594
	github.com/streamingfast/dhammer v0.0.0-20210811180702-456c4cf0a840
	github.com/streamingfast/dipp v1.0.1-0.20210811200841-d2cca4e058e6
	github.com/streamingfast/dlauncher v0.0.0-20210811194929-f06e488e63da
	github.com/streamingfast/dmesh v0.0.0-20210811181323-5a37ad73216b
	github.com/streamingfast/dmetering v0.0.0-20220301165106-a642bb6a21bd
	github.com/streamingfast/dmetrics v0.0.0-20210811180524-8494aeb34447
	github.com/streamingfast/dstore v0.1.1-0.20210811180812-4db13e99cc22
	github.com/streamingfast/dtracing v0.0.0-20220305214756-b5c0e8699839
	github.com/streamingfast/firehose v0.1.1-0.20211202153816-44577bee52dd
	github.com/streamingfast/fluxdb v0.0.0-20210811195408-0515ef659298
	github.com/streamingfast/jsonpb v0.0.0-20210811021341-3670f0aa02d0
	github.com/streamingfast/kvdb v0.0.2-0.20210811194032-09bf862bd2e3
	github.com/streamingfast/logging v0.0.0-20220304214715-bc750a74b424
	github.com/streamingfast/merger v0.0.3-0.20210820210545-ca8b1a40ae2a
	github.com/streamingfast/node-manager v0.0.2-0.20210830135731-4b00105a1479
	github.com/streamingfast/opaque v0.0.0-20210811180740-0c01d37ea308
	github.com/streamingfast/pbgo v0.0.6-0.20220228185940-1bbaafec7d8a
	github.com/streamingfast/relayer v0.0.2-0.20210812020310-adcf15941b23
	github.com/streamingfast/search v0.0.2-0.20220307144412-f4c2c6fabd9b
	github.com/streamingfast/search-client v0.0.0-20210811200417-677bdb765983
	github.com/streamingfast/shutter v1.5.0
	github.com/streamingfast/validator v0.0.0-20210812013448-b9da5752ce14
	github.com/stretchr/testify v1.7.0
	github.com/teris-io/shortid v0.0.0-20201117134242-e59966efd125
	github.com/thedevsaddam/govalidator v1.9.9
	github.com/tidwall/gjson v1.14.0
	github.com/tidwall/sjson v1.0.4
	github.com/urfave/negroni v1.0.0 // indirect
	github.com/zhongshuwen/zswchain-go v1.12.11
	go.opencensus.io v0.23.0
	go.uber.org/atomic v1.9.0
	go.uber.org/zap v1.21.0
	golang.org/x/crypto v0.0.0-20220315160706-3147a52a75dd
	golang.org/x/oauth2 v0.0.0-20211005180243-6b3c2da341f1
	golang.org/x/time v0.0.0-20191024005414-555d28b269f0
	google.golang.org/grpc v1.44.0
	google.golang.org/protobuf v1.27.1
	gopkg.in/olivere/elastic.v3 v3.0.75
	gopkg.in/yaml.v2 v2.4.0
	gotest.tools v2.2.0+incompatible
)

// to solve "github.com/ugorji/go/codec: ambiguous import: found package github.com/ugorji/go/codec in multiple modules:"
replace github.com/ugorji/go/codec => github.com/ugorji/go v1.1.2

replace github.com/graph-gophers/graphql-go => github.com/streamingfast/graphql-go v0.0.0-20210204202750-0e485a040a3c

replace github.com/census-instrumentation/opencensus-proto v0.1.0-0.20181214143942-ba49f56771b8 => github.com/census-instrumentation/opencensus-proto v0.0.3-0.20181214143942-ba49f56771b8

replace github.com/ShinyTrinkets/overseer => github.com/dfuse-io/overseer v0.2.1-0.20210326144022-ee491780e3ef

// The go-testing-interface version matches the Golang version to compile against, in this case, we want
// compatibility with 1.14 which is our minimum version. So we enforce a strict version to v1.14.1 now.

// replace github.com/streamingfast/dauth => github.com/eosnationftw/dauth v0.0.0-20210818092801-fb989ed88d76

// replace github.com/streamingfast/firehose => github.com/EOS-Nation/firehose v0.1.1-0.20211125122622-4d3db7b50f2c
replace github.com/dfuse-io/dfuse-eosio => github.com/zhongshuwen/histnew v0.2.12

replace github.com/streamingfast/pbgo => github.com/historyz/pbgo v0.1.0
replace github.com/streamingfast/bstream => github.com/streamingfast/bstream v0.0.2-0.20211006133726-b4a4315e6708