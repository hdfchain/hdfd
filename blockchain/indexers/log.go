// Copyright (c) 2016 The btcsuite developers
// Copyright (c) 2016-2019 The Hdfchain developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package indexers

import "github.com/hdfchain/slog"

// log is a logger that is initialized with no output filters.  This
// means the package will not perform any logging by default until the caller
// requests it.
// The default amount of logging is none.
var log = slog.Disabled

// UseLogger uses a specified Logger to output package logging info.
// This should be used in preference to SetLogWriter if the caller is also
// using slog.
func UseLogger(logger slog.Logger) {
	log = logger
}
