package encoding

// StringEncoding can be used to specify custom string encoding for tokens.
type StringEncoding interface {
	Encode([]byte) string
	Decode(string) ([]byte, error)
}
