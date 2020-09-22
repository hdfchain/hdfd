blockchain
==========

[![Build Status](https://github.com/hdfchain/hdfd/workflows/Build%20and%20Test/badge.svg)](https://github.com/hdfchain/hdfd/actions)
[![ISC License](https://img.shields.io/badge/license-ISC-blue.svg)](http://copyfree.org)
[![Doc](https://img.shields.io/badge/doc-reference-blue.svg)](https://pkg.go.dev/github.com/hdfchain/hdfd/blockchain/v3)

Package blockchain implements Hdfchain block handling and chain selection rules.
The test coverage is currently only around 60%, but will be increasing over
time. See `test_coverage.txt` for the gocov coverage report.  Alternatively, if
you are running a POSIX OS, you can run the `cov_report.sh` script for a
real-time report.  Package blockchain is licensed under the liberal ISC license.

There is an associated blog post about the release of this package
[here](https://blog.conformal.com/btcchain-the-bitcoin-chain-package-from-bctd/).

This package has intentionally been designed so it can be used as a standalone
package for any projects needing to handle processing of blocks into the hdfchain
block chain.

## Installation and Updating

```bash
$ go get -u github.com/hdfchain/hdfd/blockchain
```

## Hdfchain Chain Processing Overview

Before a block is allowed into the block chain, it must go through an intensive
series of validation rules.  The following list serves as a general outline of
those rules to provide some intuition into what is going on under the hood, but
is by no means exhaustive:

 - Reject duplicate blocks
 - Perform a series of sanity checks on the block and its transactions such as
   verifying proof of work, timestamps, number and character of transactions,
   transaction amounts, script complexity, and merkle root calculations
 - Compare the block against predetermined checkpoints for expected timestamps
   and difficulty based on elapsed time since the checkpoint
 - Perform a series of more thorough checks that depend on the block's position
   within the block chain such as verifying block difficulties adhere to
   difficulty retarget rules, timestamps are after the median of the last
   several blocks, all transactions are finalized, checkpoint blocks match, and
   block versions are in line with the previous blocks
 - Determine how the block fits into the chain and perform different actions
   accordingly in order to ensure any side chains which have higher difficulty
   than the main chain become the new main chain
 - When a block is being connected to the main chain (either through
   reorganization of a side chain to the main chain or just extending the
   main chain), perform further checks on the block's transactions such as
   verifying transaction duplicates, script complexity for the combination of
   connected scripts, coinbase maturity, double spends, and connected
   transaction values
 - Run the transaction scripts to verify the spender is allowed to spend the
   coins
 - Insert the block into the block database

## Examples

* [ProcessBlock Example](https://pkg.go.dev/github.com/hdfchain/hdfd/blockchain/v3#example-BlockChain.ProcessBlock)
  Demonstrates how to create a new chain instance and use ProcessBlock to
  attempt to add a block to the chain.  This example intentionally
  attempts to insert a duplicate genesis block to illustrate how an invalid
  block is handled.

## License

Package blockchain is licensed under the [copyfree](http://copyfree.org) ISC
License.
