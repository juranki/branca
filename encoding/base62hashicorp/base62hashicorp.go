package base62hashicorp

import (
	hashiBase62 "github.com/hashicorp/vault/helper/base62"
	"github.com/juranki/branca/encoding"
)

// New base62 encoding based on github.com/hashicorp/vault/helper/base62.
// It's the fastest of the base62 encodings, but not compatible with branca
// spec because it uses different character set.
func New() encoding.StringEncoding {
	return hashicorpBase62{}
}

type hashicorpBase62 struct{}

func (e hashicorpBase62) Encode(b []byte) string          { return hashiBase62.Encode(b) }
func (e hashicorpBase62) Decode(s string) ([]byte, error) { return hashiBase62.Decode(s), nil }
