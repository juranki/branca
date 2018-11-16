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
BenchmarkEncode20Bytes-8               	  100000	     15206 ns/op	     784 B/op	      13 allocs/op
BenchmarkEncode50Bytes-8               	   50000	     27561 ns/op	     896 B/op	      13 allocs/op
BenchmarkEncode100Bytes-8              	   30000	     57926 ns/op	    1248 B/op	      13 allocs/op
BenchmarkDecode20Bytes-8               	  300000	      4070 ns/op	     688 B/op	      12 allocs/op
BenchmarkDecode50Bytes-8               	  200000	      5953 ns/op	    1152 B/op	      15 allocs/op
BenchmarkDecode100Bytes-8              	  200000	      9393 ns/op	    1824 B/op	      18 allocs/op

BenchmarkEncode20BytesBasex-8          	   30000	     58547 ns/op	    2496 B/op	      12 allocs/op
BenchmarkEncode50BytesBasex-8          	   10000	    117016 ns/op	    2544 B/op	      12 allocs/op
BenchmarkEncode100BytesBasex-8         	    5000	    267366 ns/op	    5024 B/op	      14 allocs/op
BenchmarkDecode20BytesBasex-8          	  300000	      4814 ns/op	     640 B/op	       8 allocs/op
BenchmarkDecode50BytesBasex-8          	  200000	      8787 ns/op	     832 B/op	       8 allocs/op
BenchmarkDecode100BytesBasex-8         	  100000	     17546 ns/op	    1520 B/op	       9 allocs/op

BenchmarkEncode20BytesBase64-8         	 1000000	      1329 ns/op	     304 B/op	       4 allocs/op
BenchmarkEncode50BytesBase64-8         	 1000000	      1441 ns/op	     384 B/op	       4 allocs/op
BenchmarkEncode100BytesBase64-8        	 1000000	      1707 ns/op	     608 B/op	       4 allocs/op
BenchmarkDecode20BytesBase64-8         	 2000000	       691 ns/op	     208 B/op	       3 allocs/op
BenchmarkDecode50BytesBase64-8         	 2000000	       813 ns/op	     288 B/op	       3 allocs/op
BenchmarkDecode100BytesBase64-8        	 1000000	      1042 ns/op	     480 B/op	       3 allocs/op

BenchmarkEncode20BytesHashi-8          	  300000	      4559 ns/op	    1120 B/op	      15 allocs/op
BenchmarkEncode50BytesHashi-8          	  300000	      5583 ns/op	    1296 B/op	      15 allocs/op
BenchmarkEncode100BytesHashi-8         	  200000	      8559 ns/op	    2322 B/op	      19 allocs/op
BenchmarkDecode20BytesHashi-8          	 1000000	      1553 ns/op	     272 B/op	       5 allocs/op
BenchmarkDecode50BytesHashi-8          	 1000000	      2046 ns/op	     448 B/op	       6 allocs/op
BenchmarkDecode100BytesHashi-8         	  500000	      2938 ns/op	     720 B/op	       7 allocs/op
```