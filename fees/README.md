fees
=======


[![Build Status](https://github.com/hdfchain/hdfd/workflows/Build%20and%20Test/badge.svg)](https://github.com/hdfchain/hdfd/actions)
[![ISC License](https://img.shields.io/badge/license-ISC-blue.svg)](http://copyfree.org)
[![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg)](https://godoc.org/github.com/hdfchain/hdfd/fees)

Package fees provides hdfchain-specific methods for tracking and estimating fee
rates for new transactions to be mined into the network. Fee rate estimation has
two main goals:

- Ensuring transactions are mined within a target _confirmation range_
  (expressed in blocks);
- Attempting to minimize fees while maintaining be above restriction.

This package was started in order to resolve issue hdfchain/hdfd#1412 and related.
See that issue for discussion of the selected approach.

This package was developed for hdfd, a full-node implementation of Decred which
is under active development.  Although it was primarily written for
hdfd, this package has intentionally been designed so it can be used as a
standalone package for any projects needing the functionality provided.

## Installation and Updating

```bash
$ go get -u github.com/hdfchain/hdfd/fees
```

## License

Package dcrutil is licensed under the [copyfree](http://copyfree.org) ISC
License.
