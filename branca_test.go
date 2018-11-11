package branca

import (
	"testing"
	"time"
)

func Test_encode(t *testing.T) {
	codec, err := NewCodec("supersecretkeyyoushouldnotcommit")
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
	}
	for _, tt := range tests {
		t.Run(tt.message, func(t *testing.T) {
			tokenBytes, err := encode(codec.aead, nonce, []byte(tt.message), ts)
			if (err != nil) != tt.wantErr {
				t.Errorf("encode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			token := codec.base62.Encode(tokenBytes)
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
