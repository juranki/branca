package branca

import "encoding/base64"

// StringEncoding can be used to specify custom string encoding for tokens.
type StringEncoding interface {
	Encode([]byte) string
	Decode(string) ([]byte, error)
}

// Base64URLEncoding can be used to stringify tokens with base64.URLEncoding.
// It's faster than base62, but the tokens don't comply with branca specification.
var Base64URLEncoding StringEncoding

func init() {
	Base64URLEncoding = base64URLEncoding{}
}

type base64URLEncoding struct{}

func (e base64URLEncoding) Encode(b []byte) string          { return base64.URLEncoding.EncodeToString(b) }
func (e base64URLEncoding) Decode(s string) ([]byte, error) { return base64.URLEncoding.DecodeString(s) }
