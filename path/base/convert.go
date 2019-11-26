package base

import (
	"fmt"
	"github.com/roson9527/vault_eth_wallet/utils"
)

func FormatData(data, encoding string) ([]byte, error) {
	var txDataToSign []byte
	var err error

	if encoding == "hex" {
		if len(data) >= 2 && data[:2] == "0x" {
			data = data[2:]
		}
		txDataToSign, err = utils.Decode([]byte(data))
	} else if encoding == "utf8" {
		txDataToSign = []byte(data)
	} else {
		err = fmt.Errorf("invalid encoding encountered - %s", encoding)
	}

	return txDataToSign, err
}
