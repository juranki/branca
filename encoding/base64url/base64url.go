// Package base64url implements branca.StringEncoding using
// base64.URLEncoding from the standard library.
// It's faster than base62, but not compatible with the branca specification.
package base64url

import (
	"encoding/base64"

	"github.com/juranki/branca/encoding"
)

// New returns a StringEncoding based on base64.URLEncoding
func New() encoding.StringEncoding {
	return base64URLEncoding{}
}

type base64URLEncoding struct{}

func (e base64URLEncoding) Encode(b []byte) string          { return base64.URLEncoding.EncodeToString(b) }
func (e base64URLEncoding) Decode(s string) ([]byte, error) { return base64.URLEncoding.DecodeString(s) }
