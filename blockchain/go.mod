module github.com/hdfchain/hdfd/blockchain

require (
	github.com/hdfchain/hdfd/blockchain/stake v1.0.0
	github.com/hdfchain/hdfd/chaincfg v1.0.1
	github.com/hdfchain/hdfd/chaincfg/chainhash v1.0.1
	github.com/hdfchain/hdfd/database v1.0.0
	github.com/hdfchain/hdfd/dcrec v0.0.0-20180801202239-0761de129164
	github.com/hdfchain/hdfd/dcrec/edwards v0.0.0-20180721031028-5369a485acf6
	github.com/hdfchain/hdfd/dcrec/secp256k1 v1.0.0
	github.com/hdfchain/hdfd/dcrjson v1.0.0
	github.com/hdfchain/hdfd/dcrutil v1.0.0
	github.com/hdfchain/hdfd/gcs v1.0.0
	github.com/hdfchain/hdfd/txscript v1.0.0
	github.com/hdfchain/hdfd/wire v1.0.1
	github.com/hdfchain/slog v1.0.0
)

replace (
	github.com/hdfchain/hdfd/blockchain/stake => ./stake
	github.com/hdfchain/hdfd/chaincfg => ../chaincfg
	github.com/hdfchain/hdfd/chaincfg/chainhash => ../chaincfg/chainhash
	github.com/hdfchain/hdfd/database => ../database
	github.com/hdfchain/hdfd/dcrec => ../dcrec
	github.com/hdfchain/hdfd/dcrec/edwards => ../dcrec/edwards
	github.com/hdfchain/hdfd/dcrec/secp256k1 => ../dcrec/secp256k1
	github.com/hdfchain/hdfd/dcrjson => ../dcrjson
	github.com/hdfchain/hdfd/dcrutil => ../dcrutil
	github.com/hdfchain/hdfd/gcs => ../gcs
	github.com/hdfchain/hdfd/txscript => ../txscript
	github.com/hdfchain/hdfd/wire => ../wire
)
