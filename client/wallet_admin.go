package client

import (
	"encoding/json"
	"fmt"
	"github.com/roson9527/vault_eth_wallet/modules"
	"github.com/roson9527/vault_eth_wallet/path/doc"
	"github.com/roson9527/vault_eth_wallet/path/storage"
)

func (c *Client) CreateWallet(wallet *modules.Wallet) (*modules.Wallet, error) {
	payload := make(map[string]any)
	if wallet != nil {
		payload[doc.FieldNameSpaces] = wallet.NameSpaces
		payload[doc.FieldPrivateKey] = wallet.PrivateKey
	}

	sec, err := c.Meta.Logical().Write(c.conf.SecretPath+fmt.Sprintf(storage.PatternWallet, doc.NameSpaceGlobal, doc.PathSubNew), payload)
	if err != nil {
		panic(err)
	}

	out := &modules.Wallet{
		Address:   sec.Data[doc.FieldAddress].(string),
		PublicKey: sec.Data[doc.FieldPublicKey].(string),
		Network:   sec.Data[doc.FieldNetwork].(string),
	}
	out.UpdateTime, _ = sec.Data[doc.FieldUpdateTime].(json.Number).Int64()
	return out, nil
}

func (c *Client) DeleteWallet(address string) error {
	_, err := c.Meta.Logical().Delete(c.conf.SecretPath + fmt.Sprintf(storage.PatternWallet, doc.NameSpaceGlobal, address))
	return err
}

func (c *Client) UpdateWallet(wallet *modules.Wallet) error {
	payload := make(map[string]any)
	payload[doc.FieldNameSpaces] = wallet.NameSpaces
	payload[doc.FieldNetwork] = wallet.Network

	_, err := c.Meta.Logical().Write(c.conf.SecretPath+fmt.Sprintf(storage.PatternWallet, doc.NameSpaceGlobal, wallet.Address), payload)
	return err
}

func (c *Client) WalletExport(project, address string) (*modules.Wallet, error) {
	sec, err := c.Meta.Logical().Read(c.conf.SecretPath + fmt.Sprintf(storage.PatternWallet, project, address) + "/export")
	if err != nil {
		return nil, err
	}

	// 没有数据
	if sec == nil {
		return nil, nil
	}

	namespaces := make([]string, 0)
	for _, v := range sec.Data[doc.FieldNameSpaces].([]any) {
		namespaces = append(namespaces, fmt.Sprintf("%s", v))
	}

	out := &modules.Wallet{
		Address:    sec.Data[doc.FieldAddress].(string),
		PublicKey:  sec.Data[doc.FieldPublicKey].(string),
		PrivateKey: sec.Data[doc.FieldPrivateKey].(string),
		NameSpaces: namespaces,
		Network:    sec.Data[doc.FieldNetwork].(string),
	}
	out.UpdateTime, _ = sec.Data[doc.FieldUpdateTime].(json.Number).Int64()
	if err = out.Extra.Decode(sec.Data[doc.FieldExtra].(map[string]any)); err != nil {
		return nil, err
	}

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
	payload[doc.FieldPolicy] = policyMap

	_, err = c.Meta.Logical().Write(c.conf.SecretPath+fmt.Sprintf(storage.PatternPolicy, project), payload)
	return err
}

func (c *Client) ReadPolicy(project string) (*modules.Policy, error) {
	sec, err := c.Meta.Logical().Read(c.conf.SecretPath + fmt.Sprintf(storage.PatternPolicy, project))
	if err != nil {
		return nil, err
	}

	// 没有数据
	if sec == nil {
		return nil, nil
	}

	//fmt.Println(sec.Data[fieldPolicy])
	data, err := json.Marshal(sec.Data[doc.FieldPolicy])
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
