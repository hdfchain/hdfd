connmgr
=======

[![Build Status](https://github.com/hdfchain/hdfd/workflows/Build%20and%20Test/badge.svg)](https://github.com/hdfchain/hdfd/actions)
[![ISC License](https://img.shields.io/badge/license-ISC-blue.svg)](http://copyfree.org)
[![Doc](https://img.shields.io/badge/doc-reference-blue.svg)](https://pkg.go.dev/github.com/hdfchain/hdfd/connmgr/v3)

Package connmgr implements a generic Hdfchain network connection manager.

## Overview

This package handles all the general connection concerns such as maintaining a
set number of outbound connections, sourcing peers, banning, limiting max
connections, tor lookup, etc.

The package provides a generic connection manager which is able to accept
connection requests from a source or a set of given addresses, dial them and
notify the caller on connections.  The main intended use is to initialize a pool
of active connections and maintain them to remain connected to the P2P network.

In addition the connection manager provides the following utilities:

- Notifications on connections or disconnections
- Handle failures and retry new addresses from the source
- Connect only to specified addresses
- Permanent connections with increasing backoff retry timers
- Disconnect or Remove an established connection

## Installation and Updating

```bash
$ go get -u github.com/hdfchain/hdfd/connmgr
```

## License

Package connmgr is licensed under the [copyfree](http://copyfree.org) ISC License.
