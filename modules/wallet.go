package modules

import (
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/mitchellh/mapstructure"
	"github.com/roson9527/vault_eth_wallet/path/doc"
	"github.com/roson9527/vault_eth_wallet/utils"
	"math/big"
)

type Extra struct {
	IsLock   bool            `json:"is_lock" mapstructure:"is_lock"`
	Tags     []string        `json:"tags,omitempty" mapstructure:"tags,omitempty"`
	Metadata AddressAliasMap `json:"metadata,omitempty" mapstructure:"metadata,omitempty"`
}

func (w *Extra) Decode(m map[string]any) error {
	err := mapstructure.Decode(m, w)
	return err
}

var zeroInt = big.NewInt(0)

type AddressAliasMap = map[string]string

// Wallet is an Ethereum Wallet
type Wallet struct {
	PublicKey  string `json:"public_key,omitempty" mapstructure:"public_key,omitempty"` // PublicKey is the public key of the wallet
	Address    string `json:"address" mapstructure:"address"`
	UpdateTime int64  `json:"update_time" mapstructure:"update_time,omitempty"` // key pair update time
}

// WalletExtra is an Ethereum Wallet
type WalletExtra struct {
	Wallet     `mapstructure:",squash"`
	Mnemonic   string `json:"mnemonic" mapstructure:"mnemonic,omitempty"`       // Mnemonic is the mnemonic of the wallet
	PrivateKey string `json:"private_key" mapstructure:"private_key,omitempty"` // PrivateKey is the private key of the wallet
	//PublicKey    string          `json:"public_key,omitempty"`    // PublicKey is the public key of the wallet
	AddressAlias AddressAliasMap `json:"address_alias,omitempty" mapstructure:"address_alias,omitempty"` // Address is the address of the wallet
	NameSpaces   []string        `json:"namespaces,omitempty" mapstructure:"namespaces,omitempty"`       // 用于项目区分
	CryptoType   string          `json:"crypto_type,omitempty" mapstructure:"crypto_type,omitempty"`     // 加密曲线区分
	Extra        Extra           `json:"extra" mapstructure:"extra,omitempty"`                           // 用于标签
}

func (w *WalletExtra) GetAddress(chain string) string {
	addr := w.AddressAlias[chain]
	if addr != "" {
		return addr
	}
	if w.Address != "" {
		return w.Address
	}
	return w.AddressAlias[doc.ChainDefault]
}

func (w *WalletExtra) SignEthTx(unsignTx *types.Transaction) (*types.Transaction, error) {
	privateKey, err := crypto.HexToECDSA(w.PrivateKey)
	if err != nil {
		return nil, err
	}
	// 防止没有进入回收被检索到
	defer utils.ZeroKey(privateKey)

	// 不支持旧交易
	if unsignTx.ChainId() == nil || unsignTx.ChainId().Cmp(zeroInt) == 0 {
		return nil, fmt.Errorf("not support old transaction")
	}

	// 签名工具
	transactor, err := bind.NewKeyedTransactorWithChainID(privateKey, unsignTx.ChainId())
	if err != nil {
		return nil, err
	}

	return transactor.Signer(transactor.From, unsignTx)
}
