module github.com/hdfchain/hdfd/rpcclient/v5

go 1.11

require (
	github.com/davecgh/go-spew v1.1.1
	github.com/hdfchain/hdfd/chaincfg/chainhash v1.0.1
	github.com/hdfchain/hdfd/dcrjson/v3 v3.0.0
	github.com/hdfchain/hdfd/dcrutil/v2 v2.0.0
	github.com/hdfchain/hdfd/gcs/v2 v2.0.0-00010101000000-000000000000
	github.com/hdfchain/hdfd/hdkeychain/v2 v2.0.1
	github.com/hdfchain/hdfd/rpc/jsonrpc/types v1.0.0
	github.com/hdfchain/hdfd/rpc/jsonrpc/types/v2 v2.0.0-00010101000000-000000000000
	github.com/hdfchain/hdfd/wire v1.2.0
	github.com/hdfchain/hdfwallet/rpc/jsonrpc/types v1.1.0
	github.com/hdfchain/go-socks v1.0.0
	github.com/hdfchain/slog v1.0.0
	github.com/gorilla/websocket v1.4.0
)

replace (
	github.com/hdfchain/hdfd/gcs/v2 => ../gcs
	github.com/hdfchain/hdfd/rpc/jsonrpc/types/v2 => ../rpc/jsonrpc/types
)