# branca

[![Build Status](https://travis-ci.org/juranki/branca.svg?branch=master)](https://travis-ci.org/juranki/branca)
[![GoDoc](https://godoc.org/github.com/juranki/branca?status.svg)](https://godoc.org/github.com/juranki/branca)

Implements encoding and decoding for [branca tokens](https://github.com/tuupola/branca-spec).

## Install

```
go get -u github.com/juranki/branca
```

## Usage

```go
package main

import (
	"fmt"
	"log"

	"github.com/juranki/branca"
)

func main() {
	// The encryption key must be exactly 32 bytes.
	codec, err := branca.New("01234567890123456789012345678901")
	if err != nil {
		log.Fatalln(err)
	}

	token, err := codec.Encode([]byte("payload"))
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("Encrypted token: %s\n", token)

	payload, createTime, err := codec.Decode(token)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Payload: %s\n", payload)
	fmt.Printf("Creation time: %s\n", createTime)
}

```