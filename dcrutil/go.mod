module github.com/hdfchain/hdfd/dcrutil

require (
	github.com/davecgh/go-spew v1.1.0
	github.com/hdfchain/base58 v1.0.0
	github.com/hdfchain/hdfd/chaincfg v1.1.1
	github.com/hdfchain/hdfd/chaincfg/chainhash v1.0.1
	github.com/hdfchain/hdfd/dcrec v0.0.0-20180721005212-59fe2b293f69
	github.com/hdfchain/hdfd/dcrec/edwards v0.0.0-20180721005212-59fe2b293f69
	github.com/hdfchain/hdfd/dcrec/secp256k1 v1.0.0
	github.com/hdfchain/hdfd/wire v1.1.0
	golang.org/x/crypto v0.0.0-20180718160520-a2144134853f
)

replace (
	github.com/hdfchain/hdfd/chaincfg => ../chaincfg
	github.com/hdfchain/hdfd/chaincfg/chainhash => ../chaincfg/chainhash
	github.com/hdfchain/hdfd/dcrec => ../dcrec
	github.com/hdfchain/hdfd/dcrec/edwards => ../dcrec/edwards
	github.com/hdfchain/hdfd/dcrec/secp256k1 => ../dcrec/secp256k1
	github.com/hdfchain/hdfd/wire => ../wire
)
