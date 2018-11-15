// Package branca encodes and decodes branca tokens.
//
// https://github.com/tuupola/branca-spec
package branca

import (
	"crypto/cipher"
	"crypto/rand"
	"encoding/binary"
	"errors"
	"time"

	"github.com/eknkc/basex"
	"golang.org/x/crypto/chacha20poly1305"
)

const (
	version byte = 0xBA
)

// Codec encodes/decodes branca tokens
type Codec struct {
	aead           cipher.AEAD
	stringEncoding StringEncoding
}

// New creates a codec. The key must be exactly 32 bytes long.
// Tokens are stringified with base62.
func New(key string) (*Codec, error) {
	enc, err := basex.NewEncoding("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz")
	if err != nil {
		return nil, err
	}
	return NewWithEncoding(key, enc)
}

// NewWithEncoding creates a codec. The key must be exactly 32 bytes long.
// Tokens are stringified with provided encoding.
//
// WARNING!! I'll probably remove support for alternative string encoders once
// I find a fast base62 encoder that supports the branca spec.
func NewWithEncoding(key string, stringEncoding StringEncoding) (*Codec, error) {
	aead, err := chacha20poly1305.NewX([]byte(key))
	if err != nil {
		return nil, err
	}
	return &Codec{
		aead:           aead,
		stringEncoding: stringEncoding,
	}, nil
}

// Encode message
func (c *Codec) Encode(message []byte) (string, error) {
	nonce := make([]byte, 24)
	if _, err := rand.Read(nonce); err != nil {
		return "", err
	}
	return c.stringEncoding.Encode(
		encode(c.aead, nonce, message, time.Now()),
	), nil
}

// Decode message. Returns payload and creation time of token.
func (c *Codec) Decode(token string) ([]byte, time.Time, error) {
	tokenBytes, err := c.stringEncoding.Decode(token)
	if err != nil {
		return nil, time.Time{}, err
	}
	if tokenBytes[0] != version {
		return nil, time.Time{}, errors.New("invalid version")
	}
	message, err := c.aead.Open(nil, tokenBytes[5:29], tokenBytes[29:], tokenBytes[0:29])
	if err != nil {
		return nil, time.Time{}, err
	}
	ts := binary.BigEndian.Uint32(tokenBytes[1:5])
	return message, time.Unix(int64(ts), 0), nil
}

// encode assumes that slices are exactly the right length
func encode(aead cipher.AEAD, nonce, message []byte, ts time.Time) []byte {
	header := make([]byte, 29, 29+len(message)+aead.Overhead())
	header[0] = version
	binary.BigEndian.PutUint32(header[1:], uint32(ts.Unix()))
	copy(header[5:], nonce)
	return aead.Seal(header, nonce, message, header)
}
