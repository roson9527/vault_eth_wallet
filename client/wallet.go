package client

import (
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
)

func (c *Client) WalletList(project string) ([]string, error) {
	sec, err := c.Meta.Logical().List(c.conf.SecretPath + fmt.Sprintf(PatternWallet, project, ""))
	if err != nil {
		return nil, err
	}

	// 没有数据
	if sec == nil {
		return []string{}, nil
	}

	out := make([]string, 0)
	for _, v := range sec.Data["keys"].([]any) {
		out = append(out, fmt.Sprintf("%s", v))
	}

	return out, nil
}

func (c *Client) ReadWallet(project, address string) (*Wallet, error) {
	sec, err := c.Meta.Logical().Read(c.conf.SecretPath + fmt.Sprintf(PatternWallet, project, address))
	if err != nil {
		return nil, err
	}

	// 没有数据
	if sec == nil {
		return nil, nil
	}

	out := &Wallet{
		Address:   sec.Data["address"].(string),
		PublicKey: sec.Data["public_key"].(string),
		Network:   sec.Data["network"].(string),
	}
	out.UpdateTime, _ = sec.Data["update_time"].(json.Number).Int64()
	return out, nil
}

func (c *Client) SignTx(project, address string, unsignTx *types.Transaction) (*types.Transaction, error) {
	payload := make(map[string]any)
	data, err := unsignTx.MarshalBinary()
	if err != nil {
		return nil, err
	}
	payload[fieldTxBinary] = hexutil.Encode(data)

	sec, err := c.Meta.Logical().Write(c.conf.SecretPath+fmt.Sprintf(PatternWallet, project, address)+"/sign_tx", payload)
	if err != nil {
		return nil, err
	}

	var signedTx types.Transaction
	_ = signedTx.UnmarshalBinary(hexutil.MustDecode(sec.Data["tx_binary"].(string)))

	return &signedTx, nil
}
