txscript
========

[![Build Status](https://github.com/hdfchain/hdfd/workflows/Build%20and%20Test/badge.svg)](https://github.com/hdfchain/hdfd/actions)
[![ISC License](https://img.shields.io/badge/license-ISC-blue.svg)](http://copyfree.org)
[![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg)](https://godoc.org/github.com/hdfchain/hdfd/txscript)

Package txscript implements the Decred transaction script language.  There is
a comprehensive test suite.

This package has intentionally been designed so it can be used as a standalone
package for any projects needing to use or validate Decred transaction scripts.

## Decred Scripts

Decred provides a stack-based, FORTH-like language for the scripts in
the Decred transactions.  This language is not turing complete
although it is still fairly powerful.

## Installation and Updating

```bash
$ go get -u github.com/hdfchain/hdfd/txscript
```

## Examples

* [Standard Pay-to-pubkey-hash Script](https://godoc.org/github.com/hdfchain/hdfd/txscript#example-PayToAddrScript)
  Demonstrates creating a script which pays to a Decred address.  It also
  prints the created script hex and uses the DisasmString function to display
  the disassembled script.

* [Extracting Details from Standard Scripts](https://godoc.org/github.com/hdfchain/hdfd/txscript#example-ExtractPkScriptAddrs)
  Demonstrates extracting information from a standard public key script.

* [Manually Signing a Transaction Output](https://godoc.org/github.com/hdfchain/hdfd/txscript#example-SignTxOutput)
  Demonstrates manually creating and signing a redeem transaction.

* [Counting Opcodes in Scripts](https://godoc.org/github.com/hdfchain/hdfd/txscript#example-ScriptTokenizer)
  Demonstrates creating a script tokenizer instance and using it to count the
  number of opcodes a script contains.

## License

Package txscript is licensed under the [copyfree](http://copyfree.org) ISC
License.
