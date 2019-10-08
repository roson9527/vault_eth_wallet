package accounts

import (
	"context"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/roson9527/vault_eth_wallet/path/base"
)

func pathAccount(pattern string) *framework.Path {
	return &framework.Path{
		Pattern: pattern,
		// 字段
		Fields: map[string]*framework.FieldSchema{
			fieldName: {Type: framework.TypeString},
		},
		//Fields:         nil,
		ExistenceCheck: base.PathExistenceCheck,
		// 执行的位置，有read，list，create，update
		Operations: map[logical.Operation]framework.OperationHandler{
			logical.ReadOperation: &framework.PathOperation{
				Callback: read,
			},
			logical.CreateOperation: &framework.PathOperation{
				Callback: createCrossReference,
			},
			logical.DeleteOperation: &framework.PathOperation{
				Callback: deleteCrossReference,
			},
		},
		HelpSynopsis:    pathListSyn,
		HelpDescription: pathListDesc,
	}
}

func createCrossReference(
	ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {

	ret, err := create(ctx, req, data)
	if err != nil {
		return nil, err
	}

	addr := ret.Data[fieldAddress].(string)
	name := data.Get(fieldName).(string)
	err = addressWrite(ctx, req, addr, name)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

func deleteCrossReference(
	ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {

	name := data.Get(fieldName).(string)
	err := addressDelete(ctx, req, name)
	if err != nil {
		return nil, err
	}

	return delete(ctx, req, data)
}
