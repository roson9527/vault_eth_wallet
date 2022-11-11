package modules

import (
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/roson9527/vault_eth_wallet/utils"
	"math/big"
)

// Wallet is an Ethereum Wallet
type Wallet struct {
	PrivateKey string   `json:"private_key"`          // PrivateKey is the private key of the wallet
	PublicKey  string   `json:"public_key,omitempty"` // PublicKey is the public key of the wallet
	Address    string   `json:"address"`
	UpdateTime int64    `json:"update_time"` // key pair update time
	NameSpaces []string `json:"namespaces,omitempty"`
}

func (w *Wallet) SignEthTx(unsignTx *types.Transaction, chainId int64) (*types.Transaction, error) {
	privateKey, err := crypto.HexToECDSA(w.PrivateKey)
	if err != nil {
		return nil, err
	}
	// 防止没有进入回收被检索到
	defer utils.ZeroKey(privateKey)

	// 做chainId约束
	cId := big.NewInt(chainId)
	if unsignTx.Type() != types.LegacyTxType {
		if unsignTx.ChainId() == nil || unsignTx.ChainId().Cmp(cId) != 0 {
			return nil, types.ErrInvalidChainId
		}
	}

	// 签名工具
	transactor, err := bind.NewKeyedTransactorWithChainID(privateKey, cId)
	if err != nil {
		return nil, err
	}

	return transactor.Signer(transactor.From, unsignTx)
}
