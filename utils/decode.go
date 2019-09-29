package utils

import "encoding/hex"

// Decode will decode the hex
func Decode(src []byte) ([]byte, error) {
	raw := make([]byte, hex.EncodedLen(len(src)))
	n, err := hex.Decode(raw, src)
	if err != nil {
		return nil, err
	}
	raw = raw[:n]
	return raw[:], nil
}
