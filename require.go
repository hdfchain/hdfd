// Copyright (c) 2019 The Decred developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

// +build ignore

// This file exists to prevent go mod tidy from removing requires for newer
// module versions that are not yet fully integrated and to allow them to be
// automatically discovered by the testing infrastructure.
//
// It is excluded from the build to avoid including unused modules in the final
// binary.

package main

import (
	_ "github.com/hdfchain/hdfd/chaincfg/v2"
	_ "github.com/hdfchain/hdfd/dcrutil/v2"
	_ "github.com/hdfchain/hdfd/txscript/v2"
)
