module github.com/hdfchain/hdfd/peer/v2

go 1.11

require (
	github.com/davecgh/go-spew v1.1.1
	github.com/dchest/blake256 v1.1.0 // indirect
	github.com/hdfchain/hdfd/chaincfg/chainhash v1.0.2
	github.com/hdfchain/hdfd/lru v1.0.0
	github.com/hdfchain/hdfd/txscript/v2 v2.0.0
	github.com/hdfchain/hdfd/wire v1.2.0
	github.com/hdfchain/go-socks v1.0.0
	github.com/hdfchain/slog v1.0.0
)

replace github.com/hdfchain/hdfd/wire => ../wire
