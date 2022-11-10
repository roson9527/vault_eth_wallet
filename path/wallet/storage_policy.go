package wallet

import (
	"context"
	"errors"
	"fmt"
	"github.com/hashicorp/hcl"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/mitchellh/mapstructure"
	"github.com/roson9527/vault_eth_wallet/modules"
	"sync"
)

const (
	patternPolicy = "%s/policy"
)

type policyStorage struct {
	lock  sync.RWMutex
	cache *modules.Policy
}

func newPolicyStorage() *policyStorage {
	return &policyStorage{
		cache: &modules.Policy{},
	}
}

func (as *policyStorage) readPolicy(ctx context.Context, req *logical.Request, namespace string) (*modules.Policy, error) {
	as.lock.RLock()
	defer as.lock.RUnlock()

	path := fmt.Sprintf(patternPolicy, namespace)
	entry, err := req.Storage.Get(ctx, path)

	if err != nil {
		return nil, err
	}

	if entry == nil {
		return nil, errors.New("policy not found")
	}

	var policy modules.Policy
	err = entry.DecodeJSON(&policy)
	if err != nil {
		return nil, err
	}

	return &policy, nil
}

func (as *policyStorage) writePolicy(ctx context.Context, req *logical.Request, data *framework.FieldData) (*modules.Policy, error) {
	as.lock.Lock()
	defer as.lock.Unlock()

	namespace := data.Get(fieldNameSpace).(string)
	path := fmt.Sprintf(patternPolicy, namespace)
	policyData := data.Get(fieldPolicy).(map[string]any)
	policyHCL := data.Get(fieldPolicyHCL).(string)

	var p modules.Policy
	switch policyHCL {
	case "":
		// 从policyData中解析
		err := mapstructure.Decode(policyData, &p)
		if err != nil {
			return nil, err
		}
	default:
		// 从policyHCL中解析
		err := hcl.Decode(&p, policyHCL)
		if err != nil {
			return nil, err
		}
	}

	//hclog.Default().Info("data", p)
	entry, err := logical.StorageEntryJSON(path, p)

	err = req.Storage.Put(ctx, entry)
	if err != nil {
		return nil, err
	}

	return &p, nil
}
