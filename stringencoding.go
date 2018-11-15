package branca

import (
	"encoding/base64"

	"github.com/hashicorp/vault/helper/base62"
)

// StringEncoding can be used to specify custom string encoding for tokens.
type StringEncoding interface {
	Encode([]byte) string
	Decode(string) ([]byte, error)
}

// Base64URLEncoding can be used to stringify tokens with base64.URLEncoding from go library.
// It's faster than base62, but the tokens don't comply with branca specification.
var Base64URLEncoding StringEncoding

// HashicorpBase62Encoding can be used to stringify tokens with github.com/hashicorp/vault/helper/base62.
var HashicorpBase62Encoding StringEncoding

func init() {
	Base64URLEncoding = base64URLEncoding{}
	HashicorpBase62Encoding = hashicorpBase62{}
}

type base64URLEncoding struct{}

func (e base64URLEncoding) Encode(b []byte) string          { return base64.URLEncoding.EncodeToString(b) }
func (e base64URLEncoding) Decode(s string) ([]byte, error) { return base64.URLEncoding.DecodeString(s) }

type hashicorpBase62 struct{}

func (e hashicorpBase62) Encode(b []byte) string          { return base62.Encode(b) }
func (e hashicorpBase62) Decode(s string) ([]byte, error) { return base62.Decode(s), nil }
