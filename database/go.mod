module github.com/hdfchain/hdfd/database/v2

go 1.11

require (
	github.com/hdfchain/hdfd/chaincfg/chainhash v1.0.2
	github.com/hdfchain/hdfd/chaincfg/v3 v3.0.0-20200215031403-6b2ce76f0986
	github.com/hdfchain/hdfd/dcrutil/v3 v3.0.0-20200215031403-6b2ce76f0986
	github.com/hdfchain/hdfd/wire v1.3.0
	github.com/hdfchain/slog v1.0.0
	github.com/golang/protobuf v1.3.2 // indirect
	github.com/jessevdk/go-flags v1.4.0
	github.com/onsi/ginkgo v1.11.0 // indirect
	github.com/onsi/gomega v1.8.1 // indirect
	github.com/syndtr/goleveldb v1.0.1-0.20190923125748-758128399b1d
	golang.org/x/sys v0.0.0-20191010194322-b09406accb47 // indirect
)

replace (
	github.com/hdfchain/hdfd/chaincfg/v3 => ../chaincfg
	github.com/hdfchain/hdfd/dcrec/secp256k1/v3 => ../dcrec/secp256k1
	github.com/hdfchain/hdfd/dcrutil/v3 => ../dcrutil
)
