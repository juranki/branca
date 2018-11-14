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

func Test_encode(t *testing.T) {
	codec, err := New("supersecretkeyyoushouldnotcommit")
	if err != nil {
		t.Error(err)
	}
	nonce := []byte{
		0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c,
		0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c,
	}
	ts := time.Unix(123206400, 0)
	tests := []struct {
		message string
		want    string
		wantErr bool
	}{
		{"Hello world!", "875GH233T7IYrxtgXxlQBYiFobZMQdHAT51vChKsAIYCFxZtL1evV54vYqLyZtQ0ekPHt8kJHQp0a", false},
		{
			"1234567890123456789012345678901234567890",
			"1h4IciYOEawvyw9yCwKTDUnuQ6BTck6xQxYecjVIOdGbRhZfQvuqDCcDywvrDEEXFY7vwKuwfYL8aQSmg0LKH6PuqAryBB0iqPgzTtrxp8ZIu6kGhJv",
			false,
		},
		{"                    ", "2sLAhjtzkx9Wt8rZmN5KHw3HlK45sGsl1etFJ7a7wXqJcNWsMhwCFU5GH01zFy23LYwU3VBX5dEOkNcNXZ9GcQ0h", false},
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

func Test_withstringencoding(t *testing.T) {
	codec, err := NewWithEncoding("supersecretkeyyoushouldnotcommit", Base64URLEncoding)
	if err != nil {
		t.Error(err)
	}
	nonce := []byte{
		0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c,
		0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c,
	}
	ts := time.Unix(123206400, 0)
	tests := []struct {
		message string
		want    string
		wantErr bool
	}{
		{"Hello world!", "ugdX-wABAgMEBQYHCAkKCwwBAgMEBQYHCAkKCwyx3HE9f90h8oQ4xlWFeXJsXPeBrLpX5HsC4a5Q", false},
		{
			"1234567890123456789012345678901234567890",
			"ugdX-wABAgMEBQYHCAkKCwwBAgMEBQYHCAkKCwzIiy5lJcthpc9kk0YibE-ioO5sCX8LEa-Tqh4g2PsgmlyaLHfcOb3To1iRofr2ToWNdOe84u1gHw==",
			false,
		},
		{"                    ", "ugdX-wABAgMEBQYHCAkKCwwBAgMEBQYHCAkKCwzZmT1xMN12vdZ0glQxeFq0t_Z1GbcxGzDYVd7MzElP0cWT-oM=", false},
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
	codec, err := New("supersecretkeyyoushouldnotcommit")
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

func BenchmarkEncode20Bytes(b *testing.B) {
	codec, err := New("supersecretkeyyoushouldnotcommit")
	if err != nil {
		b.Error(err)
	}
	message := []byte("012345678901234567890")
	for n := 0; n < b.N; n++ {
		codec.Encode(message)
	}
}

func BenchmarkEncode50Bytes(b *testing.B) {
	codec, err := New("supersecretkeyyoushouldnotcommit")
	if err != nil {
		b.Error(err)
	}
	message := []byte("01234567890123456789012345678901234567890123456789")
	for n := 0; n < b.N; n++ {
		codec.Encode(message)
	}
}

func BenchmarkEncode100Bytes(b *testing.B) {
	codec, err := New("supersecretkeyyoushouldnotcommit")
	if err != nil {
		b.Error(err)
	}
	message := []byte("0123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789")
	for n := 0; n < b.N; n++ {
		codec.Encode(message)
	}
}

func BenchmarkDecode20Bytes(b *testing.B) {
	codec, err := New("supersecretkeyyoushouldnotcommit")
	if err != nil {
		b.Error(err)
	}

	token := "2seo65HD4ERrLj2eDcDYNZKy0YVuSfRRTkF9Nv9sOtcdZw2m5aUWDDKJhWPyH3QbTfvAzPlqZwZ79Vh4PNJabEqF"
	for n := 0; n < b.N; n++ {
		codec.Decode(token)
	}
}

func BenchmarkDecode50Bytes(b *testing.B) {
	codec, err := New("supersecretkeyyoushouldnotcommit")
	if err != nil {
		b.Error(err)
	}

	token := "AG8pFcEs2TYYzplnIGSlMv0DafABhxbhrDmkAesCD928BrgesqIffH7eg89a0zUtu6gQnKBakI20dfMwCQmySXZutJkteHwiUVzIFj0dLvDkaat3Mt3hDUMDVZPgsQWe"
	for n := 0; n < b.N; n++ {
		codec.Decode(token)
	}
}

func BenchmarkDecode100Bytes(b *testing.B) {
	codec, err := New("supersecretkeyyoushouldnotcommit")
	if err != nil {
		b.Error(err)
	}

	token := "LWTFSDXQP8ybC1i5JNqOBg2qYs6Ae3Z4qAbxBaS499FIPTViAWy56Ev98c4gLxdRKKVCcADWW60ziHhcXISDDy1q18eXu5L3ruyAF7NBLNnKSPNSZQYqTonOwmkYPRqlbj5lx3dg2h2Ju28wNYcMmdP519VlndsBQT0X6ZEc3iUoaXjBwgQMLBCzsI6Q9tkP7Yy"
	for n := 0; n < b.N; n++ {
		codec.Decode(token)
	}
}

func BenchmarkEncode20BytesBase64(b *testing.B) {
	codec, err := NewWithEncoding("supersecretkeyyoushouldnotcommit", Base64URLEncoding)
	if err != nil {
		b.Error(err)
	}
	message := []byte("012345678901234567890")
	for n := 0; n < b.N; n++ {
		codec.Encode(message)
	}
}

func BenchmarkEncode50BytesBase64(b *testing.B) {
	codec, err := NewWithEncoding("supersecretkeyyoushouldnotcommit", Base64URLEncoding)
	if err != nil {
		b.Error(err)
	}
	message := []byte("01234567890123456789012345678901234567890123456789")
	for n := 0; n < b.N; n++ {
		codec.Encode(message)
	}
}

func BenchmarkEncode100BytesBase64(b *testing.B) {
	codec, err := NewWithEncoding("supersecretkeyyoushouldnotcommit", Base64URLEncoding)
	if err != nil {
		b.Error(err)
	}
	message := []byte("0123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789")
	for n := 0; n < b.N; n++ {
		codec.Encode(message)
	}
}

func BenchmarkDecode20BytesBase64(b *testing.B) {
	codec, err := NewWithEncoding("supersecretkeyyoushouldnotcommit", Base64URLEncoding)
	if err != nil {
		b.Error(err)
	}

	token := "2seo65HD4ERrLj2eDcDYNZKy0YVuSfRRTkF9Nv9sOtcdZw2m5aUWDDKJhWPyH3QbTfvAzPlqZwZ79Vh4PNJabEqF"
	for n := 0; n < b.N; n++ {
		codec.Decode(token)
	}
}

func BenchmarkDecode50BytesBase64(b *testing.B) {
	codec, err := NewWithEncoding("supersecretkeyyoushouldnotcommit", Base64URLEncoding)
	if err != nil {
		b.Error(err)
	}

	token := "AG8pFcEs2TYYzplnIGSlMv0DafABhxbhrDmkAesCD928BrgesqIffH7eg89a0zUtu6gQnKBakI20dfMwCQmySXZutJkteHwiUVzIFj0dLvDkaat3Mt3hDUMDVZPgsQWe"
	for n := 0; n < b.N; n++ {
		codec.Decode(token)
	}
}

func BenchmarkDecode100BytesBase64(b *testing.B) {
	codec, err := NewWithEncoding("supersecretkeyyoushouldnotcommit", Base64URLEncoding)
	if err != nil {
		b.Error(err)
	}

	token := "LWTFSDXQP8ybC1i5JNqOBg2qYs6Ae3Z4qAbxBaS499FIPTViAWy56Ev98c4gLxdRKKVCcADWW60ziHhcXISDDy1q18eXu5L3ruyAF7NBLNnKSPNSZQYqTonOwmkYPRqlbj5lx3dg2h2Ju28wNYcMmdP519VlndsBQT0X6ZEc3iUoaXjBwgQMLBCzsI6Q9tkP7Yy"
	for n := 0; n < b.N; n++ {
		codec.Decode(token)
	}
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
	// fmt.Printf("%x\n", h)
	// fmt.Printf("%s\n", base64.StdEncoding.EncodeToString(h))
	// fmt.Printf("%s\n", base64.URLEncoding.EncodeToString(h))
	// fmt.Printf("%s\n", enc.Encode(h))
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
	// fmt.Printf("%x\n", h)
	// fmt.Printf("%s\n", base64.StdEncoding.EncodeToString(h))
	// fmt.Printf("%s\n", base64.URLEncoding.EncodeToString(h))
	// fmt.Printf("%s\n", enc.Encode(h))
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
