module github.com/hdfchain/hdfd

go 1.14

replace (
	github.com/hdfchain/hdfd/addrmgr => ./addrmgr
	github.com/hdfchain/hdfd/bech32 => ./bech32
	github.com/hdfchain/hdfd/blockchain/stake/v2 => ./blockchain/stake
	github.com/hdfchain/hdfd/blockchain/standalone => ./blockchain/standalone
	github.com/hdfchain/hdfd/blockchain/v2 => ./blockchain
	github.com/hdfchain/hdfd/certgen => ./certgen
	github.com/hdfchain/hdfd/chaincfg/chainhash => ./chaincfg/chainhash
	github.com/hdfchain/hdfd/chaincfg/v2 => ./chaincfg
	github.com/hdfchain/hdfd/connmgr/v2 => ./connmgr
	github.com/hdfchain/hdfd/crypto/blake256 => ./crypto/blake256
	github.com/hdfchain/hdfd/crypto/ripemd160 => ./crypto/ripemd160
	github.com/hdfchain/hdfd/database/v2 => ./database
	github.com/hdfchain/hdfd/dcrec => ./dcrec
	github.com/hdfchain/hdfd/dcrec/edwards/v2 => ./dcrec/edwards
	github.com/hdfchain/hdfd/dcrec/secp256k1/v2 => ./dcrec/secp256k1
	github.com/hdfchain/hdfd/dcrjson/v3 => ./dcrjson
	github.com/hdfchain/hdfd/dcrutil/v2 => ./dcrutil
	github.com/hdfchain/hdfd/fees/v2 => ./fees
	github.com/hdfchain/hdfd/gcs/v2 => ./gcs
	github.com/hdfchain/hdfd/hdkeychain/v2 => ./hdkeychain
	github.com/hdfchain/hdfd/internal/limits => ./limits
	github.com/hdfchain/hdfd/lru => ./lru
	github.com/hdfchain/hdfd/mempool/v3 => ./mempool
	github.com/hdfchain/hdfd/mining/v2 => ./mining
	github.com/hdfchain/hdfd/peer/v2 => ./peer
	github.com/hdfchain/hdfd/rpc/jsonrpc/types/v2 => ./rpc/jsonrpc/types
	github.com/hdfchain/hdfd/rpcclient/v5 => ./rpcclient
	github.com/hdfchain/hdfd/txscript/v2 => ./txscript
	github.com/hdfchain/hdfd/wire => ./wire
)

require (
	github.com/btcsuite/winsvc v1.0.0
	github.com/gorilla/websocket v1.4.2
	github.com/hdfchain/base58 v1.0.5
	github.com/hdfchain/go-socks v1.1.1
	github.com/hdfchain/hdfd/addrmgr v1.1.0
	github.com/hdfchain/hdfd/blockchain/stake/v2 v2.0.0-00010101000000-000000000000
	github.com/hdfchain/hdfd/blockchain/standalone v0.0.0-00010101000000-000000000000
	github.com/hdfchain/hdfd/blockchain/v2 v2.0.0-00010101000000-000000000000
	github.com/hdfchain/hdfd/certgen v0.0.0-20200910063743-aa6bbf11f012
	github.com/hdfchain/hdfd/chaincfg/chainhash v1.0.2
	github.com/hdfchain/hdfd/chaincfg/v2 v2.3.0
	github.com/hdfchain/hdfd/connmgr/v2 v2.0.0-00010101000000-000000000000
	github.com/hdfchain/hdfd/crypto/ripemd160 v1.0.0
	github.com/hdfchain/hdfd/database/v2 v2.0.0-20200910063743-aa6bbf11f012
	github.com/hdfchain/hdfd/dcrec v1.0.0
	github.com/hdfchain/hdfd/dcrec/secp256k1/v2 v2.0.0
	github.com/hdfchain/hdfd/dcrjson/v3 v3.0.1
	github.com/hdfchain/hdfd/dcrutil/v2 v2.0.1
	github.com/hdfchain/hdfd/fees/v2 v2.0.0-00010101000000-000000000000
	github.com/hdfchain/hdfd/gcs/v2 v2.0.0-20200910063743-aa6bbf11f012
	github.com/hdfchain/hdfd/hdkeychain/v2 v2.0.0-00010101000000-000000000000
	github.com/hdfchain/hdfd/lru v0.0.0-00010101000000-000000000000
	github.com/hdfchain/hdfd/mempool/v3 v3.0.0-00010101000000-000000000000
	github.com/hdfchain/hdfd/mining/v2 v2.0.0-00010101000000-000000000000
	github.com/hdfchain/hdfd/peer/v2 v2.0.0-00010101000000-000000000000
	github.com/hdfchain/hdfd/rpc/jsonrpc/types/v2 v2.0.0-20200910063743-aa6bbf11f012
	github.com/hdfchain/hdfd/rpc/jsonrpc/types v1.0.3
	github.com/hdfchain/hdfd/rpcclient/v5 v5.0.0-00010101000000-000000000000
	github.com/hdfchain/hdfd/txscript/v2 v2.0.0-00010101000000-000000000000
	github.com/hdfchain/hdfd/wire v1.3.0
	github.com/hdfchain/hdfwallet v0.0.0-20200910081459-e1799de0430f
	github.com/hdfchain/slog v1.1.1
	github.com/jessevdk/go-flags v1.4.0
	github.com/jrick/bitset v1.0.0
	github.com/jrick/logrotate v1.0.0
	golang.org/x/crypto v0.0.0-20200820211705-5c72a883971a
	golang.org/x/sync v0.0.0-20200625203802-6e8e738ad208
)
