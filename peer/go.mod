module github.com/hdfchain/hdfd/peer/v2

go 1.11

require (
	github.com/davecgh/go-spew v1.1.1
	github.com/hdfchain/hdfd/chaincfg/chainhash v1.0.2
	github.com/hdfchain/hdfd/lru v1.0.0
	github.com/hdfchain/hdfd/txscript/v3 v3.0.0-20200215031403-6b2ce76f0986
	github.com/hdfchain/hdfd/wire v1.3.0
	github.com/hdfchain/go-socks v1.1.1
	github.com/hdfchain/slog v1.0.0
)

replace (
	github.com/hdfchain/hdfd/chaincfg/v3 => ../chaincfg
	github.com/hdfchain/hdfd/dcrec/secp256k1/v3 => ../dcrec/secp256k1
	github.com/hdfchain/hdfd/dcrutil/v3 => ../dcrutil
	github.com/hdfchain/hdfd/txscript/v3 => ../txscript
)
