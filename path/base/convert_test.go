package base

import (
	"encoding/base64"
	"fmt"
	"github.com/btcsuite/btcutil/base58"
	"testing"
)

func TestFormatData(t *testing.T) {
	data := "0xd46e8dd67c5d32be8d46e8dd67c5d32be8058bb8eb970870f072445675058bb8eb970870f072445675"
	ret, _ := FormatData(data, "hex")
	//
	//fmt.Println(hex.EncodeToString(ret))
	ret2 := base64.StdEncoding.EncodeToString(ret)
	fmt.Println(ret2)
	ret3 := base58.Encode(ret)
	fmt.Println(ret3)
}
