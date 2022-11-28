package client

import (
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/mitchellh/mapstructure"
	"github.com/roson9527/vault_eth_wallet/modules"
	"github.com/roson9527/vault_eth_wallet/path/doc"
	"github.com/roson9527/vault_eth_wallet/path/storage"
)

type walletETH struct {
	*walletAdmin
}

func newWalletETH(core *core) *walletETH {
	return &walletETH{&walletAdmin{
		core,
	}}
}

func (c *walletETH) List(project string) ([]string, error) {
	sec, err := c.Meta.Logical().List(c.conf.SecretPath + fmt.Sprintf(storage.PatternWalletByChain, project, doc.CryptoSECP256K1, doc.ChainETH, ""))
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

func (c *walletETH) Read(project, address string) (*modules.Wallet, error) {
	sec, err := c.Meta.Logical().Read(c.conf.SecretPath + fmt.Sprintf(storage.PatternWalletByChain, project, doc.CryptoSECP256K1, doc.ChainETH, address))
	if err != nil {
		return nil, err
	}

	// 没有数据
	if sec == nil {
		return nil, nil
	}

	out := &modules.Wallet{
		Address:   sec.Data[doc.FieldAddress].(string),
		PublicKey: sec.Data[doc.FieldPublicKey].(string),
	}
	out.UpdateTime, _ = sec.Data[doc.FieldUpdateTime].(json.Number).Int64()
	return out, nil
}

func (c *walletETH) SignTx(project, address string, unsignTx *types.Transaction) (*types.Transaction, error) {
	payload := make(map[string]any)
	data, err := unsignTx.MarshalBinary()
	if err != nil {
		return nil, err
	}
	payload[doc.FieldTxBinary] = hexutil.Encode(data)

	sec, err := c.Meta.Logical().Write(c.conf.SecretPath+
		fmt.Sprintf(storage.PatternWalletByChain, project, doc.CryptoSECP256K1, doc.ChainETH, address)+doc.PathSubSignTx,
		payload)
	if err != nil {
		return nil, err
	}

	var signedTx types.Transaction
	_ = signedTx.UnmarshalBinary(hexutil.MustDecode(sec.Data[doc.FieldTxBinary].(string)))

	return &signedTx, nil
}
