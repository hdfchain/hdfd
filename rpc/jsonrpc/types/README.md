jsonrpc/types
=============

[![Build Status](https://github.com/hdfchain/hdfd/workflows/Build%20and%20Test/badge.svg)](https://github.com/hdfchain/hdfd/actions)
[![ISC License](https://img.shields.io/badge/license-ISC-blue.svg)](http://copyfree.org)
[![Doc](https://img.shields.io/badge/doc-reference-blue.svg)](https://pkg.go.dev/github.com/hdfchain/hdfd/rpc/jsonrpc/types/v2)

Package types implements concrete types for marshalling to and from the hdfd
JSON-RPC commands, return values, and notifications.  A comprehensive suite of
tests is provided to ensure proper functionality.

The provided types are automatically registered with
[dcrjson](https://github.com/hdfchain/hdfd/tree/master/dcrjson) when the package
is imported.  Although this package was primarily written for hdfd, it has
intentionally been designed so it can be used as a standalone package for any
projects needing to marshal to and from hdfd JSON-RPC requests and responses.

## Installation and Updating

```bash
$ go get -u github.com/hdfchain/hdfd/rpc/jsonrpc/types
```

## License

Package types is licensed under the [copyfree](http://copyfree.org) ISC License.
