module github.com/hdfchain/hdfd/chaincfg

require (
	github.com/davecgh/go-spew v1.1.0
	github.com/hdfchain/hdfd/chaincfg/chainhash v1.0.1
	github.com/hdfchain/hdfd/dcrec/edwards v0.0.0-20181208004914-a0816cf4301f
	github.com/hdfchain/hdfd/dcrec/secp256k1 v1.0.1
	github.com/hdfchain/hdfd/wire v1.2.0
)

replace (
	github.com/hdfchain/hdfd/chaincfg/chainhash => ./chainhash
	github.com/hdfchain/hdfd/dcrec/edwards => ../dcrec/edwards
	github.com/hdfchain/hdfd/dcrec/secp256k1 => ../dcrec/secp256k1
	github.com/hdfchain/hdfd/wire => ../wire
)
