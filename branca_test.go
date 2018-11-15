package branca

import (
	"encoding/base64"
	"fmt"
	"runtime"
	"sync"
	"testing"
	"time"

	"github.com/eknkc/basex"
	"github.com/gbrlsnchs/jwt/v2"
	"golang.org/x/crypto/sha3"
)

var (
	k     = "supersecretkeyyoushouldnotcommit"
	nonce = []byte{
		0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c,
		0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c,
	}
	ts = time.Unix(123206400, 0)
)

func Test_encode(t *testing.T) {
	codec, err := New(k)
	if err != nil {
		t.Error(err)
	}
	tests := []struct {
		message string
		want    string
	}{
		{
			"Hello world!",
			"875GH233T7IYrxtgXxlQBYiFobZMQdHAT51vChKsAIYCFxZtL1evV54vYqLyZtQ0ekPHt8kJHQp0a",
		},
		{
			"1234567890123456789012345678901234567890",
			"1h4IciYOEawvyw9yCwKTDUnuQ6BTck6xQxYecjVIOdGbRhZfQvuqDCcDywvrDEEXFY7vwKuwfYL8aQSmg0LKH6PuqAryBB0iqPgzTtrxp8ZIu6kGhJv",
		},
		{
			"                    ",
			"2sLAhjtzkx9Wt8rZmN5KHw3HlK45sGsl1etFJ7a7wXqJcNWsMhwCFU5GH01zFy23LYwU3VBX5dEOkNcNXZ9GcQ0h",
		},
	}
	for _, tt := range tests {
		t.Run(tt.message, func(t *testing.T) {
			tokenBytes := encode(codec.aead, nonce, []byte(tt.message), ts)
			token := codec.stringEncoding.Encode(tokenBytes)
			if token != tt.want {
				t.Errorf("encode() = %v, want %v", token, tt.want)
			}
			origMsg, origTs, err := codec.Decode(token)
			if err != nil {
				t.Errorf("error decoding %s, original message was %s", token, tt.message)
				return
			}
			if string(origMsg) != tt.message {
				t.Errorf("decode() = %s, want %s", origMsg, tt.message)
				return
			}
			if origTs != ts {
				t.Errorf("decode(time) = %v, want %v", origTs, ts)
				return
			}
		})
	}
}

func Test_encodehashicorp(t *testing.T) {
	codec, err := NewWithEncoding(k, HashicorpBase62Encoding)
	if err != nil {
		t.Error(err)
	}
	tests := []struct {
		message string
		want    string
	}{
		{
			"Hello world!",
			"875gh233t7iyRXTGxXLqbyIfOBzmqDhat51VcHkSaiycfXzTl1EVv54VyQlYzTq0EKphT8KjhqP0A",
		},
		{
			"1234567890123456789012345678901234567890",
			"1H4iCIyoeAWVYW9YcWktduNUq6btCK6XqXyECJvioDgBrHzFqVUQdcCdYWVRdeexfy7VWkUWFyl8AqsMG0lkh6pUQaRYbb0IQpGZtTRXP8ziU6KgHjV",
		},
		{
			"                    ",
			"2SlaHJTZKX9wT8RzMn5khW3hLk45SgSL1ETfj7A7WxQjCnwSmHWcfu5gh01ZfY23lyWu3vbx5DeoKnCnxz9gCq0H",
		},
	}
	for _, tt := range tests {
		t.Run(tt.message, func(t *testing.T) {
			tokenBytes := encode(codec.aead, nonce, []byte(tt.message), ts)
			token := codec.stringEncoding.Encode(tokenBytes)
			if token != tt.want {
				t.Errorf("encode() = %v, want %v", token, tt.want)
			}
			origMsg, origTs, err := codec.Decode(token)
			if err != nil {
				t.Errorf("error decoding %s, original message was %s", token, tt.message)
				return
			}
			if string(origMsg) != tt.message {
				t.Errorf("decode() = %s, want %s", origMsg, tt.message)
				return
			}
			if origTs != ts {
				t.Errorf("decode(time) = %v, want %v", origTs, ts)
				return
			}
		})
	}
}

func TestBase64(t *testing.T) {
	codec, err := NewWithEncoding(k, Base64URLEncoding)
	if err != nil {
		t.Error(err)
	}
	tests := []struct {
		message string
		want    string
	}{
		{
			"Hello world!",
			"ugdX-wABAgMEBQYHCAkKCwwBAgMEBQYHCAkKCwyx3HE9f90h8oQ4xlWFeXJsXPeBrLpX5HsC4a5Q",
		},
		{
			"1234567890123456789012345678901234567890",
			"ugdX-wABAgMEBQYHCAkKCwwBAgMEBQYHCAkKCwzIiy5lJcthpc9kk0YibE-ioO5sCX8LEa-Tqh4g2PsgmlyaLHfcOb3To1iRofr2ToWNdOe84u1gHw==",
		},
		{
			"                    ",
			"ugdX-wABAgMEBQYHCAkKCwwBAgMEBQYHCAkKCwzZmT1xMN12vdZ0glQxeFq0t_Z1GbcxGzDYVd7MzElP0cWT-oM=",
		},
	}
	for _, tt := range tests {
		t.Run(tt.message, func(t *testing.T) {
			tokenBytes := encode(codec.aead, nonce, []byte(tt.message), ts)
			token := codec.stringEncoding.Encode(tokenBytes)
			if token != tt.want {
				t.Errorf("encode() = %v, want %v", token, tt.want)
			}
			origMsg, origTs, err := codec.Decode(token)
			if err != nil {
				t.Errorf("error decoding %s, original message was %s", token, tt.message)
				return
			}
			if string(origMsg) != tt.message {
				t.Errorf("decode() = %s, want %s", origMsg, tt.message)
				return
			}
			if origTs != ts {
				t.Errorf("decode(time) = %v, want %v", origTs, ts)
				return
			}
		})
	}
}

func Test_parallel(t *testing.T) {
	var wg sync.WaitGroup
	codec, err := New(k)
	if err != nil {
		t.Error(err)
	}
	for i := 0; i < 5000; i++ {
		wg.Add(1)
		go func(msg_nro int) {
			cText, err := codec.Encode([]byte(fmt.Sprintf("message_%d", msg_nro)))
			if err != nil {
				t.Errorf("encode failed: %d\n%#v", msg_nro, err)
				return
			}
			// force yield
			runtime.Gosched()
			clearText, _, err := codec.Decode(cText)
			if err != nil {
				t.Errorf("decode failed: %d\n%#v", msg_nro, err)
				return
			}
			if fmt.Sprintf("message_%d", msg_nro) != string(clearText) {
				t.Errorf("decode expected: %s got:%s\n", fmt.Sprintf("message_%d", msg_nro), clearText)
				return
			}
		}(i)
	}
}

func benchmarkEncode(encoding StringEncoding, payload string, b *testing.B) {
	var codec *Codec
	var err error
	if encoding == nil {
		codec, err = New("supersecretkeyyoushouldnotcommit")
		if err != nil {
			b.Error(err)
		}
	} else {
		codec, err = NewWithEncoding("supersecretkeyyoushouldnotcommit", encoding)
		if err != nil {
			b.Error(err)
		}
	}
	message := []byte(payload)
	for n := 0; n < b.N; n++ {
		codec.Encode(message)
	}
}

func benchmarkDecode(encoding StringEncoding, token string, b *testing.B) {
	var codec *Codec
	var err error
	if encoding == nil {
		codec, err = New("supersecretkeyyoushouldnotcommit")
		if err != nil {
			b.Error(err)
		}
	} else {
		codec, err = NewWithEncoding("supersecretkeyyoushouldnotcommit", encoding)
		if err != nil {
			b.Error(err)
		}
	}
	for n := 0; n < b.N; n++ {
		codec.Decode(token)
	}
}

var (
	s20       = "012345678901234567890"
	s50       = "01234567890123456789012345678901234567890123456789"
	s100      = "0123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789"
	t20       = "2seo65HD4ERrLj2eDcDYNZKy0YVuSfRRTkF9Nv9sOtcdZw2m5aUWDDKJhWPyH3QbTfvAzPlqZwZ79Vh4PNJabEqF"
	t50       = "AG8pFcEs2TYYzplnIGSlMv0DafABhxbhrDmkAesCD928BrgesqIffH7eg89a0zUtu6gQnKBakI20dfMwCQmySXZutJkteHwiUVzIFj0dLvDkaat3Mt3hDUMDVZPgsQWe"
	t100      = "LWTFSDXQP8ybC1i5JNqOBg2qYs6Ae3Z4qAbxBaS499FIPTViAWy56Ev98c4gLxdRKKVCcADWW60ziHhcXISDDy1q18eXu5L3ruyAF7NBLNnKSPNSZQYqTonOwmkYPRqlbj5lx3dg2h2Ju28wNYcMmdP519VlndsBQT0X6ZEc3iUoaXjBwgQMLBCzsI6Q9tkP7Yy"
	t20_64    = "ulvtriZmUqpErrb_OzCVgIfmvrEUi8mR666J6XjVsSYrCauAupdbuLOrRevXe1gL0MAhGi_rrpmE5EQinEEe_Nio"
	t50_64    = "ulvtribN4ztMFeA2RExcGFIq89hr3b-M5iH0K33Dw84OuQN2-BINbKoXgttBMqB3rBp_Y_j732xEzN1V0dBoYDG-F69ZA3lk9A4e973b0bRmf1svatclrvNXgX4pbhU="
	t100_64   = "ulvtrib-tetyqBBT20FOw_kcKPwSAPI2A09dmbLombvOCGmnuguvsmyo0FWqYt3EwA-qC9T0Z8eRQteSs2BvPvRLJaSklmXQ626cmgJWVjXWJ0dasPftd7Dj64h4UrwCUBvbk4u-8dZiVZZqNjl6odyjRAehTuKLx8Kx6-5CUHANFkP5cAEx0gSKAX9JWCIe-Q=="
	t20Hashi  = "bTGvKHANQupg1eei7YL0bmge2analXpC9Xw9oGPDhC15Za2CXyvbf6k5rf3bO7O72Iyoh3eDgs79fodJ4Y5eAnmIe"
	t50Hashi  = "ag8Q1wk07kLj3T7QJsXt5xgjsIwcTGMSxkoc0k1CFY2WFYZnbXnrclp1w6E3kZxMZmJ2EIwztaXQP2ncmzmMHbLwAYgbV2MhM7P2cCDfUXP6dOOhHleL2e3iSFt9SnUT"
	t100Hashi = "lwth4G0ROVfABvbrOtNRLZJOC0fS7NIJ6wAs62GJepMNPfsRO788dUwZUkzumvPGQuVepcEyEqNV2s68KPETCSBM1ltHZqkeoeKsstEin4lzkKXCp6aNr4ko2eB3yTDBriV2liZuNNM854t1F2bw3eHxTbPg5QETPoirQ20294MuVzTw9yBQQgsH02oB2wH1LCA"
)

func BenchmarkEncode20Bytes(b *testing.B)        { benchmarkEncode(nil, s20, b) }
func BenchmarkEncode50Bytes(b *testing.B)        { benchmarkEncode(nil, s50, b) }
func BenchmarkEncode100Bytes(b *testing.B)       { benchmarkEncode(nil, s100, b) }
func BenchmarkDecode20Bytes(b *testing.B)        { benchmarkDecode(nil, t20, b) }
func BenchmarkDecode50Bytes(b *testing.B)        { benchmarkDecode(nil, t50, b) }
func BenchmarkDecode100Bytes(b *testing.B)       { benchmarkDecode(nil, t100, b) }
func BenchmarkEncode20BytesBase64(b *testing.B)  { benchmarkEncode(Base64URLEncoding, s20, b) }
func BenchmarkEncode50BytesBase64(b *testing.B)  { benchmarkEncode(Base64URLEncoding, s50, b) }
func BenchmarkEncode100BytesBase64(b *testing.B) { benchmarkEncode(Base64URLEncoding, s100, b) }
func BenchmarkDecode20BytesBase64(b *testing.B)  { benchmarkDecode(Base64URLEncoding, t20_64, b) }
func BenchmarkDecode50BytesBase64(b *testing.B)  { benchmarkDecode(Base64URLEncoding, t50_64, b) }
func BenchmarkDecode100BytesBase64(b *testing.B) { benchmarkDecode(Base64URLEncoding, t100_64, b) }
func BenchmarkEncode20BytesHashi(b *testing.B)   { benchmarkEncode(HashicorpBase62Encoding, s20, b) }
func BenchmarkEncode50BytesHashi(b *testing.B)   { benchmarkEncode(HashicorpBase62Encoding, s50, b) }
func BenchmarkEncode100BytesHashi(b *testing.B)  { benchmarkEncode(HashicorpBase62Encoding, s100, b) }
func BenchmarkDecode20BytesHashi(b *testing.B)   { benchmarkDecode(HashicorpBase62Encoding, t20Hashi, b) }
func BenchmarkDecode50BytesHashi(b *testing.B)   { benchmarkDecode(HashicorpBase62Encoding, t50Hashi, b) }
func BenchmarkDecode100BytesHashi(b *testing.B) {
	benchmarkDecode(HashicorpBase62Encoding, t100Hashi, b)
}

func BenchmarkSHA3Sign100BytesBase64URL(b *testing.B) {
	// enc, _ := basex.NewEncoding(base62)
	k := []byte("supersecretkeyyoushouldnotcommitsupersecretkeyyoushouldnotcommit")
	buf := []byte("and this is some data to authenticate")
	h := make([]byte, 64)
	for n := 0; n < b.N; n++ {
		shake := sha3.NewShake256()
		shake.Write(k)
		shake.Write(buf)
		shake.Read(h)
		base64.URLEncoding.EncodeToString(h)
		// enc.Encode(h)
	}
}

func BenchmarkSHA3Sign100BytesBase62(b *testing.B) {
	enc, _ := basex.NewEncoding("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz")
	k := []byte("supersecretkeyyoushouldnotcommitsupersecretkeyyoushouldnotcommit")
	buf := []byte("and this is some data to authenticate")
	h := make([]byte, 64)
	for n := 0; n < b.N; n++ {
		shake := sha3.NewShake256()
		shake.Write(k)
		shake.Write(buf)
		shake.Read(h)
		// base64.URLEncoding.EncodeToString(h)
		enc.Encode(h)
	}
}

func BenchmarkJWTVerify56Bytes(b *testing.B) {
	type Token struct {
		*jwt.JWT
		IsLoggedIn  bool   `json:"isLoggedIn"`
		CustomField string `json:"customField,omitempty"`
	}
	hs256 := jwt.NewHS256("secret")
	// paylaod is 56 bytes {"sub":"1234567890","name":"John Doe","iat":1516239022}
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9." +
		"eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ." +
		"lZ1zDoGNAv3u-OclJtnoQKejE8_viHlMtGlAxE8AE0Q"

	for n := 0; n < b.N; n++ {
		payload, sig, _ := jwt.Parse(token)
		hs256.Verify(payload, sig)
		var jot Token
		jwt.Unmarshal(payload, &jot)
	}
}
