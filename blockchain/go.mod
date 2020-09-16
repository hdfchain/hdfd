module github.com/hdfchain/hdfd/blockchain/v2

go 1.11

require (
	github.com/dchest/blake256 v1.1.0 // indirect
	github.com/hdfchain/hdfd/blockchain/stake/v2 v2.0.1
	github.com/hdfchain/hdfd/blockchain/standalone v1.0.0
	github.com/hdfchain/hdfd/chaincfg/chainhash v1.0.2
	github.com/hdfchain/hdfd/chaincfg/v2 v2.2.0
	github.com/hdfchain/hdfd/database/v2 v2.0.0
	github.com/hdfchain/hdfd/dcrec v1.0.0
	github.com/hdfchain/hdfd/dcrec/secp256k1/v2 v2.0.0
	github.com/hdfchain/hdfd/dcrutil/v2 v2.0.0
	github.com/hdfchain/hdfd/gcs/v2 v2.0.0-00010101000000-000000000000
	github.com/hdfchain/hdfd/txscript/v2 v2.0.0
	github.com/hdfchain/hdfd/wire v1.2.0
	github.com/hdfchain/slog v1.0.0
)

replace (
	github.com/hdfchain/hdfd/blockchain/standalone => ./standalone
	github.com/hdfchain/hdfd/chaincfg/v2 => ../chaincfg
	github.com/hdfchain/hdfd/gcs/v2 => ../gcs
)
