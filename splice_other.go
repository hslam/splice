// Copyright (c) 2020 Meng Huang (mhboy@outlook.com)
// This package is licensed under a MIT license that can be found in the LICENSE file.

// +build !darwin,!linux,!dragonfly,!freebsd,!netbsd,!openbsd

package splice

import (
	"net"
	"syscall"
	"time"
)

// Splice wraps the splice system call.
//
// splice() moves data between two file descriptors without copying between
// kernel address space and user address space. It transfers up to len bytes
// of data from the file descriptor rfd to the file descriptor wfd,
// where one of the descriptors must refer to a pipe.
func Splice(dst, src net.Conn, len int64) (n int64, err error) {
	return spliceBuffer(dst, src, len)
}
