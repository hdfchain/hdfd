module github.com/hdfchain/hdfd/hdkeychain/v3

go 1.13

require (
	github.com/hdfchain/base58 v1.0.5
	github.com/hdfchain/hdfd/chaincfg/v3 v3.0.0-20200215031403-6b2ce76f0986
	github.com/hdfchain/hdfd/crypto/blake256 v1.0.0
	github.com/hdfchain/hdfd/crypto/ripemd160 v1.0.0
	github.com/hdfchain/hdfd/dcrec v1.0.0
	github.com/hdfchain/hdfd/dcrec/secp256k1/v3 v3.0.0-20200215031403-6b2ce76f0986
	github.com/hdfchain/hdfd/dcrutil/v3 v3.0.0-20200215031403-6b2ce76f0986
)

replace (
	github.com/hdfchain/hdfd/chaincfg/v3 => ../chaincfg
	github.com/hdfchain/hdfd/dcrec/secp256k1/v3 => ../dcrec/secp256k1
	github.com/hdfchain/hdfd/dcrutil/v3 => ../dcrutil
)
