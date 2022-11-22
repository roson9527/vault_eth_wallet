package policy

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/fatih/structs"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/roson9527/vault_eth_wallet/path/doc"
	"github.com/roson9527/vault_eth_wallet/path/storage"
)

type pathPolicy struct {
	Storage *storage.Core
}

func (pmgr *pathPolicy) policyPath(pattern string) *framework.Path {
	return &framework.Path{
		Pattern: pattern,
		Fields: map[string]*framework.FieldSchema{
			doc.FieldNameSpace: {
				Type:        framework.TypeString,
				Description: "Namespace",
				Required:    true,
			},
			doc.FieldPolicyHCL: {Type: framework.TypeString, Default: ""},
			doc.FieldPolicy:    {Type: framework.TypeMap, Default: map[string]any{}},
		},

		Operations: map[logical.Operation]framework.OperationHandler{
			logical.UpdateOperation: &framework.PathOperation{
				Callback: pmgr.updateCallBack,
			},
			logical.CreateOperation: &framework.PathOperation{
				Callback: pmgr.updateCallBack,
			},
			logical.ReadOperation: &framework.PathOperation{
				Callback: pmgr.readCallBack,
			},
		},
		HelpSynopsis:    "",
		HelpDescription: "",
	}
}

func (pmgr *pathPolicy) updateCallBack(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	policy, err := pmgr.Storage.Policy.Write(ctx, req, data)
	if err != nil {
		return nil, err
	}

	return &logical.Response{Data: structs.Map(policy)}, nil
}

func (pmgr *pathPolicy) readCallBack(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	namespace := data.Get(doc.FieldNameSpace).(string)
	// TODO:做policy合并 global policy + namespace policy
	policy, err := pmgr.Storage.Policy.Read(ctx, req, namespace)
	if err != nil {
		return nil, err
	}
	if policy == nil {
		return nil, errors.New("policy not found")
	}

	var payload map[string]any
	jsonData, err := json.Marshal(policy)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(jsonData, &payload)
	if err != nil {
		return nil, err
	}

	return &logical.Response{
		Data: map[string]any{
			doc.FieldPolicy: payload,
		},
	}, nil
}
