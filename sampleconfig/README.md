sampleconfig
============

[![Build Status](https://img.shields.io/travis/decred/hdfd.svg)](https://travis-ci.org/decred/hdfd)
[![ISC License](https://img.shields.io/badge/license-ISC-blue.svg)](http://copyfree.org)
[![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg)](https://godoc.org/github.com/hdfchain/hdfd/sampleconfig)

Package sampleconfig provides a single constant that contains the contents of
the sample configuration file for hdfd.  This is provided for tools that perform
automatic configuration and would like to ensure the generated configuration
file not only includes the specifically configured values, but also provides
samples of other configuration options.

## Installation and Updating

```bash
$ go get -u github.com/hdfchain/hdfd/sampleconfig
```

## License

Package sampleconfig is licensed under the [copyfree](http://copyfree.org) ISC
License.
