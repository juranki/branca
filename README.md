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

Output:
```
Encrypted token: Q2BbddzrYDzFTQnBsdILsj5CDU9JAp3VDKfe3CjnLciSIIbjNM7dCealKpfmuXzYqM6wsh
Payload: payload
Creation time: 2018-11-15 00:47:03 +0200 EET
```

## Using base64URL encoding instead of base62

In the previous example, change
```go
codec, err := branca.New("01234567890123456789012345678901")
```
to
```go
codec, err := branca.NewWithEncoding("01234567890123456789012345678901", branca.Base64URLEncoding)
```

Output:
```
Encrypted token: ulvspfVaU_BlhqCk-G-XDwG6CLfZPoKGbcJpRh_B97qawKqYZYWtNsFmkL7wttVxZDWV-A==
Payload: payload
Creation time: 2018-11-15 00:47:17 +0200 EET
```

Base64 is significantly faster than base62, but it doesn't comply with branca spec. Here 
are the benchmark results from my laptop:

```
BenchmarkEncode20Bytes-8               	   30000	     59175 ns/op	    2496 B/op	      12 allocs/op
BenchmarkEncode50Bytes-8               	   10000	    121749 ns/op	    2544 B/op	      12 allocs/op
BenchmarkEncode100Bytes-8              	    5000	    280559 ns/op	    5024 B/op	      14 allocs/op
BenchmarkDecode20Bytes-8               	  300000	      5066 ns/op	     640 B/op	       8 allocs/op
BenchmarkDecode50Bytes-8               	  200000	      9001 ns/op	     832 B/op	       8 allocs/op
BenchmarkDecode100Bytes-8              	  100000	     18078 ns/op	    1520 B/op	       9 allocs/op
BenchmarkEncode20BytesBase64-8         	 1000000	      1319 ns/op	     304 B/op	       4 allocs/op
BenchmarkEncode50BytesBase64-8         	 1000000	      1349 ns/op	     384 B/op	       4 allocs/op
BenchmarkEncode100BytesBase64-8        	 1000000	      1614 ns/op	     608 B/op	       4 allocs/op
BenchmarkDecode20BytesBase64-8         	 5000000	       353 ns/op	     192 B/op	       3 allocs/op
BenchmarkDecode50BytesBase64-8         	 3000000	       435 ns/op	     240 B/op	       3 allocs/op
BenchmarkDecode100BytesBase64-8        	 3000000	       631 ns/op	     360 B/op	       3 allocs/op
```