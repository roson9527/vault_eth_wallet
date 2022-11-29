package eth

import (
	"context"
	"errors"
	"github.com/fatih/structs"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/mitchellh/mapstructure"
	"github.com/roson9527/vault_eth_wallet/path/doc"
)

func (cb *callback) put(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	policy, err := cb.Storage.Policy.Write(ctx, req, data)
	if err != nil {
		return nil, err
	}

	return &logical.Response{Data: structs.Map(policy)}, nil
}

func (cb *callback) read(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	namespace := data.Get(doc.FieldNameSpace).(string)
	chain := data.Get(doc.FieldChain).(string)
	// TODO:做policy合并 global policy + namespace policy
	policy, err := cb.Storage.Policy.Read(ctx, req, namespace, doc.CryptoSECP256K1, chain)
	if err != nil {
		return nil, err
	}
	if policy == nil {
		return nil, errors.New("policy not found")
	}

	var payload map[string]any
	err = mapstructure.Decode(policy, &payload)
	if err != nil {
		return nil, err
	}

	return &logical.Response{
		Data: map[string]any{
			doc.FieldPolicy: payload,
		},
	}, nil
}
