package client

import (
	"errors"
	"fmt"
	"github.com/mitchellh/mapstructure"
	"github.com/roson9527/vault_eth_wallet/modules"
	"github.com/roson9527/vault_eth_wallet/path/doc"
	"github.com/roson9527/vault_eth_wallet/path/storage"
	"gopkg.in/ffmt.v1"
)

func (c *walletText) Update(payload *modules.WalletExtra) error {
	var data map[string]any
	if payload == nil {
		return errors.New("payload is nil")
	}

	err := mapstructure.Decode(payload, &data)
	if err != nil {
		return err
	}

	ffmt.Print(data)

	_, err = c.Meta.Logical().Write(c.conf.SecretPath+fmt.Sprintf(storage.PatternWallet, doc.NameSpaceGlobal,
		doc.CryptoTEXT, payload.Address), data)
	return err
}

func (c *walletText) Delete(address string) error {
	_, err := c.Meta.Logical().Delete(c.conf.SecretPath + fmt.Sprintf(storage.PatternWallet,
		doc.NameSpaceGlobal, doc.CryptoTEXT, address))
	return err
}

func (c *walletText) Export(project, chain, address string) (*modules.WalletExtra, error) {
	sec, err := c.Meta.Logical().Read(c.conf.SecretPath +
		fmt.Sprintf(storage.PatternWalletByChain, project, doc.CryptoTEXT, chain, address) + doc.PathSubExport)
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
