package client

import (
	"encoding/json"
	"fmt"
	"github.com/roson9527/vault_eth_wallet/modules"
	"github.com/roson9527/vault_eth_wallet/path/doc"
	"github.com/roson9527/vault_eth_wallet/path/storage"
)

type policy struct {
	*core
}

func newPolicy(core *core) *policy {
	return &policy{core: core}
}

func (c *policy) Update(project, chain string, policy *modules.Policy) error {
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

	_, err = c.Meta.Logical().Write(c.conf.SecretPath+fmt.Sprintf(storage.PatternPolicy, project, doc.CryptoSECP256K1, chain), payload)
	return err
}

func (c *policy) Read(project, chain string) (*modules.Policy, error) {
	sec, err := c.Meta.Logical().Read(c.conf.SecretPath + fmt.Sprintf(storage.PatternPolicy, project, doc.CryptoSECP256K1, chain))
	if err != nil {
		return nil, err
	}

	// 没有数据
	if sec == nil {
		return nil, nil
	}

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
