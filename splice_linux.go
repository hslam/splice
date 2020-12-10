// Copyright (c) 2020 Meng Huang (mhboy@outlook.com)
// This package is licensed under a MIT license that can be found in the LICENSE file.

// +build linux

package splice

import (
	"errors"
	"net"
	"syscall"
	"time"
)

// ErrSyscallConn will be returned when the net.Conn do not implements the syscall.Conn interface.
var ErrSyscallConn = errors.New("The net.Conn do not implements the syscall.Conn interface")

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
		return 0, err
	}
	srcFd, err = netFd(src)
	if err != nil {
		return 0, err
	}
	var p [2]int
	syscall.ForkLock.RLock()
	err = syscall.Pipe(p[0:])
	if err != nil {
		syscall.ForkLock.RUnlock()
		return 0, err
	}
	syscall.CloseOnExec(p[0])
	syscall.CloseOnExec(p[1])
	syscall.ForkLock.RUnlock()
	rFd := int(p[0])
	wFd := int(p[1])
	if len > maxSpliceSize {
		len = maxSpliceSize
	}
	var retain int64
	retain, err = splice(srcFd, nil, wFd, nil, int(len), spliceNonblock)
	if err != nil {
		return 0, err
	}
	var out int64
	for retain > 0 {
		out, err = splice(rFd, nil, dstFd, nil, int(retain), spliceNonblock)
		if out > 0 {
			retain -= out
			n += out
			continue
		}
		if err != syscall.EAGAIN {
			return n, err
		}
		time.Sleep(time.Microsecond * 10)
	}
	return n, nil
}

func netFd(conn net.Conn) (int, error) {
	syscallConn, dstOk := conn.(syscall.Conn)
	if !dstOk {
		return 0, ErrSyscallConn
	}
	return fd(syscallConn)
}

func fd(c syscall.Conn) (int, error) {
	var nfd int
	dstRaw, err := c.SyscallConn()
	if err != nil {
		return 0, err
	}
	dstRaw.Control(func(fd uintptr) {
		nfd = int(fd)
	})
	return nfd, nil
}

func splice(rfd int, roff *int64, wfd int, woff *int64, len int, flags int) (n int64, err error) {
	return syscall.Splice(rfd, roff, wfd, woff, len, flags)
}

func tee(rfd int, wfd int, len int, flags int) (n int64, err error) {
	return syscall.Tee(rfd, wfd, len, flags)
}
