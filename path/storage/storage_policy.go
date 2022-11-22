package storage

import (
	"context"
	"errors"
	"fmt"
	"github.com/hashicorp/hcl"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/mitchellh/mapstructure"
	"github.com/roson9527/vault_eth_wallet/modules"
	"github.com/roson9527/vault_eth_wallet/path/doc"
)

const (
	PatternPolicy = "%s/policy"
)

type policyStorage struct {
}

func newPolicyStorage() *policyStorage {
	return &policyStorage{}
}

func (as *policyStorage) Read(ctx context.Context, req *logical.Request, namespace string) (*modules.Policy, error) {
	path := fmt.Sprintf(PatternPolicy, namespace)
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

func (as *policyStorage) Write(ctx context.Context, req *logical.Request, data *framework.FieldData) (*modules.Policy, error) {
	namespace := data.Get(doc.FieldNameSpace).(string)
	path := fmt.Sprintf(PatternPolicy, namespace)
	policyData := data.Get(doc.FieldPolicy).(map[string]any)
	policyHCL := data.Get(doc.FieldPolicyHCL).(string)

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

	entry, err := logical.StorageEntryJSON(path, p)

	err = req.Storage.Put(ctx, entry)
	if err != nil {
		return nil, err
	}

	return &p, nil
}
