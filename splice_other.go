// Copyright (c) 2020 Meng Huang (mhboy@outlook.com)
// This package is licensed under a MIT license that can be found in the LICENSE file.

// +build !darwin,!linux,!dragonfly,!freebsd,!netbsd,!openbsd

package splice

import (
	"net"
)

// NewContext returns a new context.
func NewContext() (*Context, error) {
	pool := assignPool(maxSpliceSize)
	buf := pool.Get().([]byte)
	return &Context{buffer: buf, pool: pool}, nil
}

// Close closes the context.
func (ctx *Context) Close() {
	ctx.pool.Put(ctx.buffer[:cap(ctx.buffer)])
}

// Splice wraps the splice system call.
//
// splice() moves data between two file descriptors without copying between
// kernel address space and user address space. It transfers up to len bytes
// of data from the file descriptor rfd to the file descriptor wfd,
// where one of the descriptors must refer to a pipe.
func Splice(dst, src net.Conn, ctx *Context, len int64) (n int64, err error) {
	return spliceBuffer(dst, src, ctx, len)
}
