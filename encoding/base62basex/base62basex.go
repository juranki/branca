package base62basex

import (
	"log"

	"github.com/juranki/branca/encoding"

	"github.com/eknkc/basex"
)

// New StringEncoding, based on github.com/eknkc/basex.
// It's compatible with the branca spec, but a little slow.
func New() encoding.StringEncoding {
	basexBase62Encoding, err := basex.NewEncoding("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz")
	if err != nil {
		log.Fatal(err)
	}
	return basexBase62Encoding
}
