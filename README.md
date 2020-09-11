hdfd
====

[![Build Status](https://github.com/hdfchain/hdfd/workflows/Build%20and%20Test/badge.svg)](https://github.com/hdfchain/hdfd/actions)
[![ISC License](https://img.shields.io/badge/license-ISC-blue.svg)](http://copyfree.org)
[![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg)](https://godoc.org/github.com/hdfchain/hdfd)
[![Go Report Card](https://goreportcard.com/badge/github.com/hdfchain/hdfd)](https://goreportcard.com/report/github.com/hdfchain/hdfd)

## Hdfchain Overview

Hdfchain is a blockchain-based cryptocurrency with a strong focus on community
input, open governance, and sustainable funding for development. It utilizes a
hybrid proof-of-work and proof-of-stake mining system to ensure that a small
group cannot dominate the flow of transactions or make changes to Hdfchain without
the input of the community.  A unit of the currency is called a `hdfchain` (DCR).

https://clkj.ltd

## Latest Downloads

https://clkj.ltd/downloads

## What is hdfd?

hdfd is a full node implementation of Hdfchain written in Go (golang).

It acts as a fully-validating chain daemon for the Hdfchain cryptocurrency.  hdfd
maintains the entire past transactional ledger of Hdfchain and allows relaying of
transactions to other Hdfchain nodes around the world.

This software is currently under active development.  It is extremely stable and
has been in production use since February 2016.

The software was originally forked from [btcd](https://github.com/btcsuite/btcd),
which is a bitcoin full node implementation that is still under active
development.  To gain the benefit of btcd's ongoing upgrades, including improved
peer and connection handling, database optimization, and other blockchain
related technology improvements, hdfd is continuously synced with the btcd
codebase.

## What is a full node?

The term 'full node' is short for 'fully-validating node' and refers to software
that fully validates all transactions and blocks, as opposed to trusting a 3rd
party.  In addition to validating transactions and blocks, nearly all full nodes
also participate in relaying transactions and blocks to other full nodes around
the world, thus forming the peer-to-peer network that is the backbone of the
Hdfchain cryptocurrency.

The full node distinction is important, since full nodes are not the only type
of software participating in the Hdfchain peer network. For instance, there are
'lightweight nodes' which rely on full nodes to serve the transactions, blocks,
and cryptographic proofs they require to function, as well as relay their
transactions to the rest of the global network.

## Why run hdfd?

As described in the previous section, the Hdfchain cryptocurrency relies on having
a peer-to-peer network of nodes that fully validate all transactions and blocks
and then relay them to other full nodes.

Running a full node with hdfd contributes to the overall security of the
network, increases the available paths for transactions and blocks to relay,
and helps ensure there are an adequate number of nodes available to serve
lightweight clients, such as Simplified Payment Verification (SPV) wallets.

Without enough full nodes, the network could be unable to expediently serve
users of lightweight clients which could force them to have to rely on
centralized services that significantly reduce privacy and are vulnerable to
censorship.

In terms of individual benefits, since hdfd fully validates every block and
transaction, it provides the highest security and privacy possible when used in
conjunction with a wallet that also supports directly connecting to it in full
validation mode, such as [hdfwallet (CLI)](https://github.com/hdfchain/hdfwallet)
and [Hdfchainiton (GUI)](https://github.com/hdfchain/hdfchainiton).

## Minimum Recommended Specifications (hdfd only)

* 10 GB disk space (as of September 2018, increases over time)
* 1GB memory (RAM)
* ~150MB/day download, ~1.5GB/day upload
  * Plus one-time initial download of the entire block chain
* Windows 7/8.x/10 (server preferred), macOS, Linux
* High uptime

## Getting Started

So, you've decided to help the network by running a full node.  Great!  Running
hdfd is simple.  All you need to do is install hdfd on a machine that is
connected to the internet and meets the minimum recommended specifications, and
launch it.

Also, make sure your firewall is configured to allow inbound connections to port
9108.

<a name="Installation" />

## Installing and updating

### Binaries (Windows/Linux/macOS)

Binary releases are provided for common operating systems and architectures:

https://clkj.ltd/downloads

### Build from source (all platforms)

Building or updating from source requires the following build dependencies:

- **Go 1.11 or 1.12**

  Installation instructions can be found here: https://golang.org/doc/install.
  It is recommended to add `$GOPATH/bin` to your `PATH` at this point.

- **Git**

  Installation instructions can be found at https://git-scm.com or
  https://gitforwindows.org.

To build and install from a checked-out repo, run `go install . ./cmd/...` in
the repo's root directory.  Some notes:

* Set the `GO111MODULE=on` environment variable if building from within
  `GOPATH`.

* The `hdfd` executable will be installed to `$GOPATH/bin`.  `GOPATH`
  defaults to `$HOME/go` (or `%USERPROFILE%\go` on Windows) if unset.


### Example of obtaining and building from source on Windows 10:

```PowerShell
PS> git clone https://github.com/hdfchain/hdfd $env:USERPROFILE\src\hdfd
PS> cd $env:USERPROFILE\src\hdfd
PS> go install . .\cmd\...
PS> & "$(go env GOPATH)\bin\hdfd" -V

```

## Docker

### Running hdfd

You can run a hdfchain node from inside a docker container.  To build the image
yourself, use the following command:

```
docker build -t hdfchain/hdfd .
```

Or you can create an alpine based image (requires Docker 17.05 or higher):

```
docker build -t hdfchain/hdfd:alpine -f Dockerfile.alpine .
```

You can then run the image using:

```
docker run hdfchain/hdfd
```

You may wish to use an external volume to customise your config and persist the
data in an external volume:

```
docker run --rm -v /home/user/dcrdata:/root/.hdfd/data hdfchain/hdfd
```

For a minimal image, you can use the hdfchain/hdfd:alpine tag.  This is typically
a more secure option while also being a much smaller image.

You can run hdfctl from inside the image.  For example, run an image (mounting
your data from externally) with:

```
docker run --rm -ti --name=hdfd-1 -v /home/user/.hdfd:/root/.hdfd \
  hdfchain/hdfd:alpine
```

And then run hdfctl commands against it.  For example:

```
docker exec -ti hdfd-1 hdfctl getbestblock
```

### Running Tests

All tests and linters may be run in a docker (or podman) container using the
script `run_tests.sh` by specifying either `docker` or `podman` as the first
parameter.  This script defaults to using the current latest supported version
of Go, but it also respects the `GOVERSION` environment variable set to the
major version of Go to allow testing on a previous version of Go.  Generally,
Hdfchain only supports the current and previous major versions of Go.

```
./run_tests.sh docker
```

To run the tests locally without docker on the latest supported version of Go:

```
./run_tests.sh
```

To run the tests locally without docker on Go 1.11:

```
GOVERSION=1.11 ./run_tests.sh
```

## Contact

If you have any further questions you can find us at:

https://clkj.ltd/community

## Issue Tracker

The [integrated github issue tracker](https://github.com/hdfchain/hdfd/issues)
is used for this project.

## Documentation

The documentation for hdfd is a work-in-progress.  It is located in the
[docs](https://github.com/hdfchain/hdfd/tree/master/docs) folder.

## License

hdfd is licensed under the [copyfree](http://copyfree.org) ISC License.
