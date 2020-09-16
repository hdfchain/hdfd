module github.com/hdfchain/hdfd

go 1.13

require (
	github.com/btcsuite/winsvc v1.0.0
	github.com/davecgh/go-spew v1.1.1
	github.com/gorilla/websocket v1.4.2
	github.com/hdfchain/base58 v1.0.5
	github.com/hdfchain/go-socks v1.1.1
	github.com/hdfchain/hdfd/addrmgr v1.1.0
	github.com/hdfchain/hdfd/bech32 v1.0.0
	github.com/hdfchain/hdfd/blockchain/stake/v3 v3.0.0-20200215031403-6b2ce76f0986
	github.com/hdfchain/hdfd/blockchain/standalone/v2 v2.0.0
	github.com/hdfchain/hdfd/blockchain/v3 v3.0.0-20200215031403-6b2ce76f0986
	github.com/hdfchain/hdfd/certgen v1.1.0
	github.com/hdfchain/hdfd/chaincfg/chainhash v1.0.2
	github.com/hdfchain/hdfd/chaincfg/v3 v3.0.0-20200215031403-6b2ce76f0986
	github.com/hdfchain/hdfd/connmgr/v3 v3.0.0-20200215031403-6b2ce76f0986
	github.com/hdfchain/hdfd/crypto/ripemd160 v1.0.0
	github.com/hdfchain/hdfd/database/v2 v2.0.1
	github.com/hdfchain/hdfd/dcrec v1.0.0
	github.com/hdfchain/hdfd/dcrec/secp256k1/v3 v3.0.0-20200215031403-6b2ce76f0986
	github.com/hdfchain/hdfd/dcrjson/v3 v3.0.1
	github.com/hdfchain/hdfd/dcrutil/v3 v3.0.0-20200503044000-76f6906e50e5
	github.com/hdfchain/hdfd/gcs/v2 v2.0.1
	github.com/hdfchain/hdfd/hdkeychain/v3 v3.0.0
	github.com/hdfchain/hdfd/lru v1.0.0
	github.com/hdfchain/hdfd/peer/v2 v2.1.0
	github.com/hdfchain/hdfd/rpc/jsonrpc/types/v2 v2.0.1-0.20200503044000-76f6906e50e5
	github.com/hdfchain/hdfd/rpcclient/v6 v6.0.0-20200215031403-6b2ce76f0986
	github.com/hdfchain/hdfd/txscript/v3 v3.0.0-20200215031403-6b2ce76f0986
	github.com/hdfchain/hdfd/wire v1.3.0
	github.com/hdfchain/slog v1.0.0
	github.com/jessevdk/go-flags v1.4.0
	github.com/jrick/bitset v1.0.0
	github.com/jrick/logrotate v1.0.0
	github.com/syndtr/goleveldb v1.0.1-0.20190923125748-758128399b1d
	golang.org/x/crypto v0.0.0-20190611184440-5c40567a22f8
	golang.org/x/sync v0.0.0-20181221193216-37e7f081c4d4
)

replace (
	github.com/hdfchain/hdfd/addrmgr => ./addrmgr
	github.com/hdfchain/hdfd/bech32 => ./bech32
	github.com/hdfchain/hdfd/blockchain/stake/v3 => ./blockchain/stake
	github.com/hdfchain/hdfd/blockchain/standalone/v2 => ./blockchain/standalone
	github.com/hdfchain/hdfd/blockchain/v3 => ./blockchain
	github.com/hdfchain/hdfd/certgen => ./certgen
	github.com/hdfchain/hdfd/chaincfg/chainhash => ./chaincfg/chainhash
	github.com/hdfchain/hdfd/chaincfg/v3 => ./chaincfg
	github.com/hdfchain/hdfd/connmgr/v3 => ./connmgr
	github.com/hdfchain/hdfd/crypto/blake256 => ./crypto/blake256
	github.com/hdfchain/hdfd/crypto/ripemd160 => ./crypto/ripemd160
	github.com/hdfchain/hdfd/database/v2 => ./database
	github.com/hdfchain/hdfd/dcrec => ./dcrec
	github.com/hdfchain/hdfd/dcrec/secp256k1/v3 => ./dcrec/secp256k1
	github.com/hdfchain/hdfd/dcrjson/v3 => ./dcrjson
	github.com/hdfchain/hdfd/dcrutil/v3 => ./dcrutil
	github.com/hdfchain/hdfd/gcs/v2 => ./gcs
	github.com/hdfchain/hdfd/hdkeychain/v3 => ./hdkeychain
	github.com/hdfchain/hdfd/limits => ./limits
	github.com/hdfchain/hdfd/lru => ./lru
	github.com/hdfchain/hdfd/peer/v2 => ./peer
	github.com/hdfchain/hdfd/rpc/jsonrpc/types/v2 => ./rpc/jsonrpc/types
	github.com/hdfchain/hdfd/rpcclient/v6 => ./rpcclient
	github.com/hdfchain/hdfd/txscript/v3 => ./txscript
	github.com/hdfchain/hdfd/wire => ./wire
)

replace github.com/hdfchain/hdfd/dcrec/edwards/v2 v2.0.0 => ./dcrec/edwards
