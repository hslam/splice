// Copyright (c) 2020 Meng Huang (mhboy@outlook.com)
// This package is licensed under a MIT license that can be found in the LICENSE file.

// +build darwin dragonfly freebsd netbsd openbsd

package splice

import (
	"github.com/hslam/shm"
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
	bufferSize := maxSpliceSize
	if bufferSize < int(len) {
		bufferSize = int(len)
	}
	shmid, buf, err := shm.GetAttach(shm.IPC_PRIVATE, bufferSize, 0)
	if err != nil {
		return spliceBuffer(dst, src, len)
	}
	defer shm.Remove(shmid)
	defer shm.Detach(buf)
	var retain int
	retain, err = src.Read(buf)
	if err != nil {
		return 0, err
	}
	var out int
	var pos int
	for retain > 0 {
		out, err = dst.Write(buf[pos : pos+retain])
		if out > 0 {
			retain -= out
			n += int64(out)
			pos += out
			continue
		}
		if err != syscall.EAGAIN {
			return n, err
		}
		time.Sleep(time.Microsecond * 10)
	}
	return n, nil
}
