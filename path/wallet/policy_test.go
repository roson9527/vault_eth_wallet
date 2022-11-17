package wallet

import (
	"fmt"
	"testing"
)

func TestPolicy_Read(t *testing.T) {
	cli := createLocalClient()
	sec, err := cli.Logical().Read("mock_web3/" + fmt.Sprintf(patternPolicy, "test_b"))
	if err != nil {
		panic(err)
	}
	if sec == nil {
		// 没有读取到任何数据
	}
	fmt.Println(sec, err)
}

type M map[string]any

func TestPolicy_Write(t *testing.T) {
	cli := createLocalClient()
	payload := map[string]any{
		"policy": M{
			"chain_ids": []int{1},
			"contract": M{
				"0xA69babEF1cA67A37Ffaf7a485DfFF3382056e78C": M{
					"func_sign": M{
						"balance_of": M{
							"sign":      "0x1cff79cd",
							"max_value": "0x5c0200",
						},
					},
				},
			},
		},
	}
	sec, err := cli.Logical().Write("mock_web3/"+fmt.Sprintf(patternPolicy, "test_b"), payload)
	if err != nil {
		panic(err)
	}
	if sec == nil {
		// 没有读取到任何数据
	}
	fmt.Println(sec, err)
}
