Certgen
======

[![Build Status](https://github.com/hdfchain/hdfd/workflows/Build%20and%20Test/badge.svg)](https://github.com/hdfchain/hdfd/actions)
[![ISC License](https://img.shields.io/badge/license-ISC-blue.svg)](http://copyfree.org)
[![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg)](https://godoc.org/github.com/hdfchain/hdfd/certgen)

## Overview

This package contains functions for creating self-signed TLS certificate from
random new key pairs, typically used for encrypting RPC and websocket
communications.

ECDSA certificates are supported on all Go versions.  Beginning with Go 1.13,
this package additionally includes support for Ed25519 certificates.

## Installation and Updating

```bash
$ go get -u github.com/hdfchain/hdfd/certgen
```

## License

Package certgen is licensed under the [copyfree](http://copyfree.org) ISC
License.
