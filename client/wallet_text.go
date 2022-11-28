package client

import (
	"fmt"
	"github.com/mitchellh/mapstructure"
	"github.com/roson9527/vault_eth_wallet/modules"
	"github.com/roson9527/vault_eth_wallet/path/doc"
	"github.com/roson9527/vault_eth_wallet/path/storage"
)

type walletText struct {
	*core
}

func newWalletText(core *core) *walletText {
	return &walletText{core: core}
}

func (c *walletText) List(project, chain string) ([]string, error) {
	sec, err := c.Meta.Logical().List(c.conf.SecretPath + fmt.Sprintf(storage.PatternWalletByChain,
		project, doc.CryptoTEXT, chain, ""))
	if err != nil {
		return nil, err
	}

	// 没有数据
	if sec == nil {
		return []string{}, nil
	}

	out := make([]string, 0)
	err = mapstructure.Decode(sec.Data[doc.FieldKeys], &out)
	return out, err
}

func (c *walletText) Read(project, chain, address string) (*modules.WalletExtra, error) {
	sec, err := c.Meta.Logical().Read(c.conf.SecretPath +
		fmt.Sprintf(storage.PatternWalletByChain, project, doc.CryptoTEXT, chain, address))
	if err != nil {
		return nil, err
	}

	fmt.Println("sec", sec)

	// 没有数据
	if sec == nil {
		return nil, nil
	}

	var out modules.WalletExtra
	err = mapstructure.WeakDecode(sec.Data, &out)
	if err != nil {
		return nil, err
	}

	return &out, nil
}
