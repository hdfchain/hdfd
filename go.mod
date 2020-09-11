module github.com/hdfchain/hdfd

go 1.14

replace (
	github.com/hdfchain/hdfd/addrmgr => ./addrmgr
	github.com/hdfchain/hdfd/blockchain => ./blockchain
	github.com/hdfchain/hdfd/blockchain/stake => ./blockchain/stake
	github.com/hdfchain/hdfd/certgen => ./certgen
	github.com/hdfchain/hdfd/chaincfg/chainhash => ./chaincfg/chainhash
	github.com/hdfchain/hdfd/chaincfg/v2 => ./chaincfg
	github.com/hdfchain/hdfd/connmgr => ./connmgr
	github.com/hdfchain/hdfd/database => ./database
	github.com/hdfchain/hdfd/dcrec => ./dcrec
	github.com/hdfchain/hdfd/dcrjson/v3 => ./dcrjson
	github.com/hdfchain/hdfd/dcrutil/v2 => ./dcrutil
	github.com/hdfchain/hdfd/fees => ./fees
	github.com/hdfchain/hdfd/gcs => ./gcs
	github.com/hdfchain/hdfd/hdkeychain/v2 => ./hdkeychain
	github.com/hdfchain/hdfd/internal/limits => ./limits
	github.com/hdfchain/hdfd/lru => ./lru
	github.com/hdfchain/hdfd/mempool/v2 => ./mempool
	github.com/hdfchain/hdfd/mining => ./mining
	github.com/hdfchain/hdfd/peer => ./peer
	github.com/hdfchain/hdfd/rpc/jsonrpc/types => ./rpc/jsonrpc/types
	github.com/hdfchain/hdfd/rpcclient/v5 => ./rpcclient
	github.com/hdfchain/hdfd/txscript/v2 => ./txscript
	github.com/hdfchain/hdfd/wire => ./wire
)

replace github.com/hdfchain/hdfd/dcrjson/v3 v3.0.1 => github.com/jrick/btcd/dcrjson/v2 v2.0.0-20190715200557-9fffa6c80ab0

replace github.com/hdfchain/hdfwallet/rpc/jsonrpc/types v1.1.0 => github.com/jrick/btcwallet/rpc/jsonrpc/types v0.0.0-20190715193601-785bca9161e7
