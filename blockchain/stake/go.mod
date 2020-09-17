module github.com/hdfchain/hdfd/blockchain/stake/v3

go 1.13

require (
	github.com/hdfchain/hdfd/chaincfg/chainhash v1.0.2
	github.com/hdfchain/hdfd/database/v2 v2.0.1
	github.com/hdfchain/hdfd/dcrec v1.0.0
	github.com/hdfchain/hdfd/wire v1.3.0
	github.com/hdfchain/slog v1.0.0
)

replace (
	github.com/hdfchain/hdfd/chaincfg/v3 => ../../chaincfg
	github.com/hdfchain/hdfd/dcrec/secp256k1/v3 => ../../dcrec/secp256k1
	github.com/hdfchain/hdfd/dcrutil/v3 => ../../dcrutil
	github.com/hdfchain/hdfd/txscript/v3 => ../../txscript
)
