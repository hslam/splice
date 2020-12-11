# splice
[![PkgGoDev](https://pkg.go.dev/badge/github.com/hslam/splice)](https://pkg.go.dev/github.com/hslam/splice)
[![Build Status](https://travis-ci.org/hslam/splice.svg?branch=master)](https://travis-ci.org/hslam/splice)
[![Go Report Card](https://goreportcard.com/badge/github.com/hslam/splice)](https://goreportcard.com/report/github.com/hslam/splice)
[![LICENSE](https://img.shields.io/github/license/hslam/splice.svg?style=flat-square)](https://github.com/hslam/splice/blob/master/LICENSE)

Package splice wraps the splice system call.

## Get started

### Install
```
go get github.com/hslam/splice
```
### Import
```
import "github.com/hslam/splice"
```
### Usage
#### Example
```go
package main

import (
	"fmt"
	"github.com/hslam/splice"
	"io"
	"net"
	"time"
)

func main() {
	contents := "Hello world"
	lis, err := net.Listen("tcp", ":9999")
	if err != nil {
		panic(err)
	}
	defer lis.Close()
	done := make(chan bool)
	go func() {
		conn, _ := lis.Accept()
		defer conn.Close()
		time.Sleep(time.Millisecond * 100)
		if _, err := splice.Splice(conn, conn, nil, 1024); err != nil && err != io.EOF {
			panic(err)
		}
		close(done)
	}()
	conn, _ := net.Dial("tcp", "127.0.0.1:9999")
	conn.Write([]byte(contents))
	buf := make([]byte, 64)
	n, _ := conn.Read(buf)
	fmt.Println(string(buf[:n]))
	conn.Close()
	<-done
}
```

### Output
```
Hello world
```

### License
This package is licensed under a MIT license (Copyright (c) 2020 Meng Huang)


### Author
splice was written by Meng Huang.


