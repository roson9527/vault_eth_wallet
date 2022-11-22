package client

import (
	"encoding/json"
	"fmt"
	"github.com/roson9527/vault_eth_wallet/modules"
)

func (c *Client) CreateWallet(wallet *Wallet) (*Wallet, error) {
	payload := make(map[string]any)
	if wallet != nil {
		payload[fieldNameSpaces] = wallet.NameSpaces
		payload[fieldPrivateKey] = wallet.PrivateKey
	}

	sec, err := c.Meta.Logical().Write(c.conf.SecretPath+fmt.Sprintf(PatternWallet, NameSpaceGlobal, "new"), payload)
	if err != nil {
		panic(err)
	}

	out := &Wallet{
		Address:   sec.Data[fieldAddress].(string),
		PublicKey: sec.Data[fieldPublicKey].(string),
		Network:   sec.Data[fieldNetwork].(string),
	}
	out.UpdateTime, _ = sec.Data[fieldUpdateTime].(json.Number).Int64()
	return out, nil
}

func (c *Client) DeleteWallet(address string) error {
	_, err := c.Meta.Logical().Delete(c.conf.SecretPath + fmt.Sprintf(PatternWallet, NameSpaceGlobal, address))
	return err
}

func (c *Client) UpdateWallet(wallet *Wallet) error {
	payload := make(map[string]any)
	payload[fieldNameSpaces] = wallet.NameSpaces
	payload[fieldNetwork] = wallet.Network

	_, err := c.Meta.Logical().Write(c.conf.SecretPath+fmt.Sprintf(PatternWallet, NameSpaceGlobal, wallet.Address), payload)
	return err
}

func (c *Client) WalletExport(project, address string) (*Wallet, error) {
	sec, err := c.Meta.Logical().Read(c.conf.SecretPath + fmt.Sprintf(PatternWallet, project, address) + "/export")
	if err != nil {
		return nil, err
	}

	// 没有数据
	if sec == nil {
		return nil, nil
	}

	namespaces := make([]string, 0)
	for _, v := range sec.Data[fieldNameSpaces].([]any) {
		namespaces = append(namespaces, fmt.Sprintf("%s", v))
	}

	out := &Wallet{
		Address:    sec.Data[fieldAddress].(string),
		PublicKey:  sec.Data[fieldPublicKey].(string),
		PrivateKey: sec.Data[fieldPrivateKey].(string),
		NameSpaces: namespaces,
		Network:    sec.Data[fieldNetwork].(string),
	}
	out.UpdateTime, _ = sec.Data[fieldUpdateTime].(json.Number).Int64()
	return out, nil
}

func (c *Client) UpdatePolicy(project string, policy *modules.Policy) error {
	data, err := json.Marshal(policy)
	if err != nil {
		return err
	}
	policyMap := make(map[string]any)
	err = json.Unmarshal(data, &policyMap)
	if err != nil {
		return err
	}

	payload := make(map[string]any)
	payload[fieldPolicy] = policyMap

	_, err = c.Meta.Logical().Write(c.conf.SecretPath+fmt.Sprintf(PatternPolicy, project), payload)
	return err
}

func (c *Client) ReadPolicy(project string) (*modules.Policy, error) {
	sec, err := c.Meta.Logical().Read(c.conf.SecretPath + fmt.Sprintf(PatternPolicy, project))
	if err != nil {
		return nil, err
	}

	// 没有数据
	if sec == nil {
		return nil, nil
	}

	//fmt.Println(sec.Data[fieldPolicy])
	data, err := json.Marshal(sec.Data[fieldPolicy])
	if err != nil {
		return nil, err
	}

	out := &modules.Policy{}
	err = json.Unmarshal(data, out)
	if err != nil {
		return nil, err
	}
	return out, nil
}
