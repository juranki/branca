package base62

import (
	"encoding/binary"
	"fmt"
	"math"
	"math/big"
)

var (
	charSetBytes = []byte("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz")
	charValues   map[byte]uint64
	decodeBase   *big.Int
)

func init() {
	charValues = map[byte]uint64{}
	for i, v := range charSetBytes {
		charValues[v] = uint64(i)
	}
	decodeBase = big.NewInt(62)
	decodeBase.Exp(decodeBase, big.NewInt(10), nil)
}

// Encode bytes into base62 string
func Encode(b []byte) string {
	var slot *big.Int
	zero := big.NewInt(0)
	divident := big.NewInt(0)
	quotient := big.NewInt(0)
	remainder := big.NewInt(0)
	divisor := big.NewInt(62)

	divident.SetBytes(b)

	i := int(float64(divident.BitLen()) / math.Log2(float64(62)))
	bs := make([]byte, i+1)

	for divident.Cmp(zero) > 0 {
		quotient.DivMod(divident, divisor, remainder)
		slot = quotient
		quotient = divident
		divident = slot
		bs[i] = charSetBytes[remainder.Int64()]
		i--
	}
	return string(bs[i+1:])
}

// Decode base62 string
func Decode(s string) ([]byte, error) {
	l := len(s)
	bs := []byte(s)
	v := big.NewInt(0)
	acc := big.NewInt(0)
	mul := big.NewInt(1)
	uint64Bytes := make([]byte, 8)
	for l > 0 {
		low := l - 10
		if low < 0 {
			low = 0
		}
		err := decodeSmall(bs[low:l], uint64Bytes)
		if err != nil {
			return nil, err
		}
		v.SetBytes(uint64Bytes)
		v.Mul(v, mul)
		acc.Add(acc, v)
		mul.Mul(mul, decodeBase)
		l -= 10
	}
	return acc.Bytes(), nil
}

func decodeSmall(s []byte, t []byte) error {
	l := len(s)
	var acc, mul, base uint64
	base = 62
	mul = 1
	for i := l - 1; i >= 0; i-- {
		v, ok := charValues[s[i]]
		if !ok {
			return fmt.Errorf("invalid character in %s at index %d", s, i)
		}
		acc += mul * v
		if i > 0 {
			mul *= base
		}
	}
	binary.BigEndian.PutUint64(t, acc)
	return nil
}
