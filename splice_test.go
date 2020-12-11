// Copyright (c) 2020 Meng Huang (mhboy@outlook.com)
// This package is licensed under a MIT license that can be found in the LICENSE file.

package splice

import (
	"io/ioutil"
	"net"
	"sync"
	"syscall"
	"testing"
	"time"
)

func TestSplice(t *testing.T) {
	addr := "127.0.0.1:8888"
	proxyAddr := "127.0.0.1:9999"
	contents := "Hello world"
	wg := sync.WaitGroup{}
	// Start server listening on a socket.
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		t.Error(err)
	}
	defer lis.Close()
	wg.Add(1)
	go func() {
		defer wg.Done()
		conn, err := lis.Accept()
		if err != nil {
			t.Error(err)
		}
		defer conn.Close()
		b, _ := ioutil.ReadAll(conn)
		if string(b) != contents {
			t.Errorf("contents not transmitted: got %s (len=%d), want %s\n", string(b), len(b), contents)
		}
	}()

	// Start server listening on a socket.
	plis, err := net.Listen("tcp", proxyAddr)
	if err != nil {
		t.Error(err)
	}
	defer plis.Close()
	wg.Add(1)
	go func() {
		defer wg.Wait()
		defer wg.Done()
		conn, err := plis.Accept()
		if err != nil {
			t.Error(err)
		}
		defer conn.Close()
		proxy, err := net.Dial("tcp", addr)
		if err != nil {
			t.Error(err)
		}
		defer proxy.Close()
		time.Sleep(time.Millisecond * 100)
		written, err := Splice(proxy, conn, nil, 1024)
		if err != nil && err != syscall.EAGAIN && err != EOF {
			t.Error(err)
		}
		if int(written) != len(contents) {
			t.Error()
		}
	}()
	// Send source file to server.
	conn, err := net.Dial("tcp", proxyAddr)
	if err != nil {
		t.Error(err)
	}
	conn.Write([]byte(contents))
	conn.Close()
	wg.Wait()
}

func TestSpliceContext(t *testing.T) {
	addr := "127.0.0.1:8888"
	proxyAddr := "127.0.0.1:9999"
	contents := "Hello world"
	wg := sync.WaitGroup{}
	// Start server listening on a socket.
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		t.Error(err)
	}
	defer lis.Close()
	wg.Add(1)
	go func() {
		defer wg.Done()
		conn, err := lis.Accept()
		if err != nil {
			t.Error(err)
		}
		defer conn.Close()
		b, _ := ioutil.ReadAll(conn)
		if string(b) != contents {
			t.Errorf("contents not transmitted: got %s (len=%d), want %s\n", string(b), len(b), contents)
		}
	}()

	// Start server listening on a socket.
	plis, err := net.Listen("tcp", proxyAddr)
	if err != nil {
		t.Error(err)
	}
	defer plis.Close()
	wg.Add(1)
	go func() {
		defer wg.Wait()
		defer wg.Done()
		conn, err := plis.Accept()
		if err != nil {
			t.Error(err)
		}
		defer conn.Close()
		proxy, err := net.Dial("tcp", addr)
		if err != nil {
			t.Error(err)
		}
		defer proxy.Close()
		time.Sleep(time.Millisecond * 100)
		ctx, err := NewContext()
		if err != nil {
			t.Error(err)
		}
		defer ctx.Close()
		written, err := Splice(proxy, conn, ctx, 1024)
		if err != nil && err != syscall.EAGAIN && err != EOF {
			t.Error(err)
		}
		if int(written) != len(contents) {
			t.Error()
		}
	}()
	// Send source file to server.
	conn, err := net.Dial("tcp", proxyAddr)
	if err != nil {
		t.Error(err)
	}
	conn.Write([]byte(contents))
	conn.Close()
	wg.Wait()
}

func TestSpliceBuffer(t *testing.T) {
	addr := "127.0.0.1:8888"
	proxyAddr := "127.0.0.1:9999"
	contents := "Hello world"
	wg := sync.WaitGroup{}
	// Start server listening on a socket.
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		t.Error(err)
	}
	defer lis.Close()
	wg.Add(1)
	go func() {
		defer wg.Done()
		conn, err := lis.Accept()
		if err != nil {
			t.Error(err)
		}
		defer conn.Close()
		b, _ := ioutil.ReadAll(conn)
		if string(b) != contents {
			t.Errorf("contents not transmitted: got %s (len=%d), want %s\n", string(b), len(b), contents)
		}
	}()

	// Start server listening on a socket.
	plis, err := net.Listen("tcp", proxyAddr)
	if err != nil {
		t.Error(err)
	}
	defer plis.Close()
	wg.Add(1)
	go func() {
		defer wg.Wait()
		defer wg.Done()
		conn, err := plis.Accept()
		if err != nil {
			t.Error(err)
		}
		defer conn.Close()
		proxy, err := net.Dial("tcp", addr)
		if err != nil {
			t.Error(err)
		}
		defer proxy.Close()
		time.Sleep(time.Millisecond * 100)
		written, err := spliceBuffer(proxy, conn, nil, 1024)
		if err != nil && err != syscall.EAGAIN && err != EOF {
			t.Error(err)
		}
		if int(written) != len(contents) {
			t.Error()
		}
	}()
	// Send source file to server.
	conn, err := net.Dial("tcp", proxyAddr)
	if err != nil {
		t.Error(err)
	}
	conn.Write([]byte(contents))
	conn.Close()
	wg.Wait()
}
