// Copyright (c) 2020 Meng Huang (mhboy@outlook.com)
// This package is licensed under a MIT license that can be found in the LICENSE file.

// +build linux

package splice

import (
	"errors"
	"net"
	"syscall"
)

const (
	// spliceNonblock makes calls to splice(2) non-blocking.
	spliceNonblock = 0x2

	// maxSpliceSize is the maximum amount of data Splice asks
	// the kernel to move in a single call to splice(2).
	maxSpliceSize = 4 << 20
)

// ErrSyscallConn will be returned when the net.Conn do not implements the syscall.Conn interface.
var ErrSyscallConn = errors.New("The net.Conn do not implements the syscall.Conn interface")

// newContext returns a new context.
func newContext(b *bucket) (ctx *context, err error) {
	var p [2]int
	syscall.ForkLock.RLock()
	err = syscall.Pipe(p[0:])
	if err == nil {
		syscall.CloseOnExec(p[0])
		syscall.CloseOnExec(p[1])
		ctx = &context{reader: int(p[0]), writer: int(p[1]), bucket: b}
	}
	syscall.ForkLock.RUnlock()
	return ctx, err
}

// Close closes the context.
func (ctx *context) Close() {
	syscall.Close(ctx.reader)
	syscall.Close(ctx.writer)
}

// Splice wraps the splice system call.
//
// splice() moves data between two file descriptors without copying between
// kernel address space and user address space. It transfers up to len bytes
// of data from the file descriptor rfd to the file descriptor wfd,
// where one of the descriptors must refer to a pipe.
func Splice(dst, src net.Conn, len int64) (n int64, err error) {
	var srcFd, dstFd int
	dstFd, err = netFd(dst)
	if err != nil {
		return spliceBuffer(dst, src, len)
	}
	srcFd, err = netFd(src)
	if err != nil {
		return spliceBuffer(dst, src, len)
	}
	var rFd, wFd int
	b := assignBucket(dstFd).GetInstance()
	var ctx *context
	ctx, err = b.Get()
	if err != nil {
		return 0, ErrNotHandled
	}
	defer b.Free(ctx)
	ctx.alive = false
	rFd = ctx.reader
	wFd = ctx.writer
	if len > maxSpliceSize {
		len = maxSpliceSize
	}
	var remain int64
	// If remain == 0 && err == nil, src is at EOF, and the
	// transfer is complete.
	remain, err = splice(srcFd, nil, wFd, nil, int(len), spliceNonblock)
	if err != nil {
		return 0, err
	}
	if remain == 0 {
		return 0, EOF
	}
	var out int64
	for remain > 0 {
		out, err = splice(rFd, nil, dstFd, nil, int(remain), spliceNonblock)
		if out > 0 {
			remain -= out
			n += out
			continue
		}
		if err != syscall.EAGAIN {
			return n, EOF
		}
	}
	ctx.alive = true
	return n, nil
}

func netFd(conn net.Conn) (int, error) {
	syscallConn, ok := conn.(syscall.Conn)
	if !ok {
		return 0, ErrSyscallConn
	}
	return fd(syscallConn)
}

func fd(c syscall.Conn) (int, error) {
	var nfd int
	raw, err := c.SyscallConn()
	if err != nil {
		return 0, err
	}
	raw.Control(func(fd uintptr) {
		nfd = int(fd)
	})
	return nfd, nil
}

func splice(rfd int, roff *int64, wfd int, woff *int64, len int, flags int) (n int64, err error) {
	return syscall.Splice(rfd, roff, wfd, woff, len, flags)
}
