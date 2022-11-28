package client

import (
	"errors"
	"fmt"
	"github.com/mitchellh/mapstructure"
	"github.com/roson9527/vault_eth_wallet/modules"
	"github.com/roson9527/vault_eth_wallet/path/doc"
	"github.com/roson9527/vault_eth_wallet/path/storage"
)

type walletAdmin struct {
	*core
}

func newWalletAdmin(c *core) *walletAdmin {
	return &walletAdmin{c}
}

func (c *walletAdmin) Create(wallet *modules.WalletExtra) (*modules.WalletExtra, error) {
	payload := make(map[string]any)
	if wallet != nil {
		payload[doc.FieldNameSpaces] = wallet.NameSpaces
		payload[doc.FieldPrivateKey] = wallet.PrivateKey
	}

	sec, err := c.Meta.Logical().Write(c.conf.SecretPath+fmt.Sprintf(storage.PatternWalletAdmin, doc.NameSpaceGlobal, doc.CryptoSECP256K1), payload)
	if err != nil {
		return nil, err
	}

	var out modules.WalletExtra
	err = mapstructure.Decode(sec.Data, &out)
	//fmt.Println("create:", sec.Data, out)

	//out.UpdateTime, _ = sec.Data[doc.FieldUpdateTime].(json.Number).Int64()
	return &out, err
}

func (c *walletAdmin) Delete(address string) error {
	_, err := c.Meta.Logical().Delete(c.conf.SecretPath + fmt.Sprintf(storage.PatternWallet, doc.NameSpaceGlobal, doc.CryptoSECP256K1, address))
	return err
}

func (c *walletAdmin) Update(wallet *modules.WalletExtra) error {
	payload := make(map[string]any)
	payload[doc.FieldAddress] = wallet.Address
	payload[doc.FieldNameSpaces] = wallet.NameSpaces
	payload[doc.FieldAddressAlias] = wallet.AddressAlias

	//fmt.Println("update:", payload)

	_, err := c.Meta.Logical().Write(c.conf.SecretPath+fmt.Sprintf(storage.PatternWallet, doc.NameSpaceGlobal,
		doc.CryptoSECP256K1, wallet.Address), payload)
	return err
}

func (c *walletAdmin) Export(project, address string) (*modules.WalletExtra, error) {
	sec, err := c.Meta.Logical().Read(c.conf.SecretPath +
		fmt.Sprintf(storage.PatternWalletByChain, project, doc.CryptoSECP256K1, doc.ChainETH, address) + doc.PathSubExport)
	if err != nil {
		return nil, err
	}

	// 没有数据
	if sec == nil {
		return nil, errors.New("not found")
	}

	var out modules.WalletExtra
	err = mapstructure.WeakDecode(sec.Data, &out)
	if err != nil {
		return nil, err
	}

	return &out, nil
}
