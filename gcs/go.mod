module github.com/hdfchain/hdfd/gcs/v2

go 1.11

require (
	github.com/dchest/siphash v1.2.1
	github.com/hdfchain/hdfd/blockchain/stake/v3 v3.0.0-20200215031403-6b2ce76f0986
	github.com/hdfchain/hdfd/chaincfg/chainhash v1.0.2
	github.com/hdfchain/hdfd/crypto/blake256 v1.0.0
	github.com/hdfchain/hdfd/txscript/v3 v3.0.0-20200215031403-6b2ce76f0986
	github.com/hdfchain/hdfd/wire v1.3.0
)

replace (
	github.com/hdfchain/hdfd/blockchain/stake/v3 => ../blockchain/stake
	github.com/hdfchain/hdfd/chaincfg/v3 => ../chaincfg
	github.com/hdfchain/hdfd/dcrec/secp256k1/v3 => ../dcrec/secp256k1
	github.com/hdfchain/hdfd/dcrutil/v3 => ../dcrutil
	github.com/hdfchain/hdfd/txscript/v3 => ../txscript
)
