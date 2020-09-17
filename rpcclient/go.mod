module github.com/hdfchain/hdfd/rpcclient/v6

go 1.13

require (
	github.com/hdfchain/hdfd/chaincfg/chainhash v1.0.2
	github.com/hdfchain/hdfd/dcrjson/v3 v3.0.1
	github.com/hdfchain/hdfd/gcs/v2 v2.0.1
	github.com/hdfchain/hdfd/wire v1.3.0
	github.com/hdfchain/go-socks v1.1.0
	github.com/hdfchain/slog v1.0.0
	github.com/gorilla/websocket v1.4.2
)

replace (
	github.com/hdfchain/hdfd/chaincfg/v3 => ../chaincfg
	github.com/hdfchain/hdfd/dcrec/secp256k1/v3 => ../dcrec/secp256k1
	github.com/hdfchain/hdfd/dcrutil/v3 => ../dcrutil
	github.com/hdfchain/hdfd/hdkeychain/v3 => ../hdkeychain
	github.com/hdfchain/hdfd/rpc/jsonrpc/types/v2 => ../rpc/jsonrpc/types
)
