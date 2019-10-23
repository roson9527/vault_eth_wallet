package base

import (
	"github.com/roson9527/vault_eth_wallet/modules"
	"math/big"
	"testing"
)

func TestSignatureTx(t *testing.T) {
	testAccount, _ := GenerateKey()
	params := modules.SignParams{
		Nonce:     4,
		ToAddress: nil,
		Amount:    big.NewInt(1),
		GasLimit:  1,
		GasPrice:  big.NewInt(1),
		Data:      nil,
		ChainId:   big.NewInt(1),
	}
	ret, err := SignatureTx(testAccount, &params)
	if err != nil {
		t.Error(err)
	}
	t.Log(ret.TransactionHash, ret.Signed)
}
