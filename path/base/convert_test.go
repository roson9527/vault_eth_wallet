package base

import (
	"fmt"
	"testing"
)

func TestFormatData(t *testing.T) {
	data:="0xd46e8dd67c5d32be8d46e8dd67c5d32be8058bb8eb970870f072445675058bb8eb970870f072445675"
	ret, _ :=FormatData(data, "hex")
	fmt.Println(ret)
}
