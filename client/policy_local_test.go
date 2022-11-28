package client

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/roson9527/vault_eth_wallet/modules"
	"gopkg.in/ffmt.v1"
	"math/big"
	"testing"
)

func TestPolicy_Update(t *testing.T) {
	cli := createLocalClient()
	_ = cli.Policy.Update("test_c", "eth", &modules.Policy{
		ChainIds:             []uint64{1},
		EnableCreateContract: false,
		Contract: map[string]modules.ContractConfig{
			"VotingEscrow": {
				Address: "0x0e42acBD23FAee03249DAFF896b78d7e79fBD58E",
				FuncSigns: map[string]modules.FuncSign{
					"create_lock": {
						Sign:     "0x65fc3873",
						MaxValue: "50",
					},
				},
			},
		},
	})
	res, _ := cli.Policy.Read("test_c", "eth")
	ffmt.Print(res)
}

// 看看是否能拦截到禁止的合约
func TestPolicy_Verify(t *testing.T) {
	cli := createLocalClient()
	wallets, _ := cli.Wallet.ETH.List("test_c")
	fmt.Println(wallets)

	to := common.HexToAddress("0x0e42acBD23FAee03249DAFF896b78d7e79fBD58E")
	// 构造一个交易
	tx := &types.DynamicFeeTx{
		ChainID:   big.NewInt(2),
		GasTipCap: big.NewInt(1e9),
		GasFeeCap: big.NewInt(1e9),
		Gas:       100000,
		To:        &to,
		Value:     big.NewInt(50),
		Nonce:     11,
		Data:      hexutil.MustDecode("0x65fc3873"),
	}
	unsignTx := types.NewTx(tx)

	signTx, err := cli.Wallet.ETH.SignTx("test_c", wallets[0], unsignTx)
	fmt.Println(signTx, err)
}
