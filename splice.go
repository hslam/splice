// Copyright (c) 2020 Meng Huang (mhboy@outlook.com)
// This package is licensed under a MIT license that can be found in the LICENSE file.

// Package splice wraps the splice system call.
package splice

import (
	"io"
	"net"
	"sync"
	"sync/atomic"
	"syscall"
	"time"
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

var (
	// EOF is the error returned by Read when no more input is available.
	// Functions should return EOF only to signal a graceful end of input.
	// If the EOF occurs unexpectedly in a structured data stream,
	// the appropriate error is either ErrUnexpectedEOF or some other error
	// giving more detail.
	EOF = io.EOF

	buffers = sync.Map{}
	assign  int32
)

func assignPool(size int) *sync.Pool {
	for {
		if p, ok := buffers.Load(size); ok {
			return p.(*sync.Pool)
		}
		if atomic.CompareAndSwapInt32(&assign, 0, 1) {
			var pool = &sync.Pool{New: func() interface{} {
				return make([]byte, size)
			}}
			buffers.Store(size, pool)
			atomic.StoreInt32(&assign, 0)
			return pool
		}
	}
}

// Context represents a splice context.
type Context struct {
	buffer []byte
	writer int
	reader int
	shmid  int
	pool   *sync.Pool
}

func spliceBuffer(dst, src net.Conn, ctx *Context, len int64) (n int64, err error) {
	bufferSize := maxSpliceSize
	if bufferSize < int(len) {
		bufferSize = int(len)
	}
	var buf []byte
	if ctx != nil {
		buf = ctx.buffer[:bufferSize]
	} else {
		pool := assignPool(bufferSize)
		buf = pool.Get().([]byte)
		defer pool.Put(buf)
	}
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
