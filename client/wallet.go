package client

import (
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/roson9527/vault_eth_wallet/modules"
	"github.com/roson9527/vault_eth_wallet/path/doc"
	"github.com/roson9527/vault_eth_wallet/path/storage"
)

func (c *Client) WalletList(project string) ([]string, error) {
	sec, err := c.Meta.Logical().List(c.conf.SecretPath + fmt.Sprintf(storage.PatternWallet, project, ""))
	if err != nil {
		return nil, err
	}

	// 没有数据
	if sec == nil {
		return []string{}, nil
	}

	out := make([]string, 0)
	for _, v := range sec.Data[doc.FieldKeys].([]any) {
		out = append(out, fmt.Sprintf("%s", v))
	}

	return out, nil
}

func (c *Client) ReadWallet(project, address string) (*modules.Wallet, error) {
	sec, err := c.Meta.Logical().Read(c.conf.SecretPath + fmt.Sprintf(storage.PatternWallet, project, address))
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
		Network:   sec.Data[doc.FieldNetwork].(string),
	}
	out.UpdateTime, _ = sec.Data[doc.FieldUpdateTime].(json.Number).Int64()
	return out, nil
}

func (c *Client) SignTx(project, address string, unsignTx *types.Transaction) (*types.Transaction, error) {
	payload := make(map[string]any)
	data, err := unsignTx.MarshalBinary()
	if err != nil {
		return nil, err
	}
	payload[doc.FieldTxBinary] = hexutil.Encode(data)

	sec, err := c.Meta.Logical().Write(c.conf.SecretPath+fmt.Sprintf(storage.PatternWallet, project, address)+doc.PathSubSignTx, payload)
	if err != nil {
		return nil, err
	}

	var signedTx types.Transaction
	_ = signedTx.UnmarshalBinary(hexutil.MustDecode(sec.Data[doc.FieldTxBinary].(string)))

	return &signedTx, nil
}
