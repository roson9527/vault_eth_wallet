package base

import (
	"encoding/hex"
	"fmt"
)

func FormatData(data, encoding string) ([]byte, error) {
	var txDataToSign []byte
	var err error

	if encoding == "hex" {
		if has0xPrefix(data) {
			data = data[2:]
		}
		if len(data)%2 == 1 {
			data = "0" + data
		}
		txDataToSign, err = hex.DecodeString(data)

	} else if encoding == "utf8" {
		txDataToSign = []byte(data)

	} else {
		err = fmt.Errorf("invalid encoding encountered - %s", encoding)
	}

	return txDataToSign, err
}

func has0xPrefix(str string) bool {
	return len(str) >= 2 && str[0] == '0' && (str[1] == 'x' || str[1] == 'X')
}
