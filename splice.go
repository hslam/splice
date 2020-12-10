// Copyright (c) 2020 Meng Huang (mhboy@outlook.com)
// This package is licensed under a MIT license that can be found in the LICENSE file.

// Package splice wraps the splice system call.
package splice

import (
	"syscall"
)

const (
	// spliceNonblock makes calls to splice(2) non-blocking.
	spliceNonblock = 0x2

	// maxSpliceSize is the maximum amount of data Splice asks
	// the kernel to move in a single call to splice(2).
	maxSpliceSize = 64 << 10

	// EAGAIN will be returned when resource temporarily unavailable.
	EAGAIN = syscall.EAGAIN
)
