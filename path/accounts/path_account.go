package accounts

import (
	"context"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

func pathAccount(pattern string) *framework.Path {
	return &framework.Path{
		Pattern: pattern,
		// 字段
		Fields: map[string]*framework.FieldSchema{
			fieldName: {Type: framework.TypeString},
		},
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

		addr:=data.Get(fieldAddress).(string)
		name:=data.Get(fieldName).(string)
	err:= addressWrite(ctx, req, addr, name)
	if err!=nil {
		return nil, err
	}

	return create(ctx, req, data)
}

func deleteCrossReference(
	ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	name:=data.Get(fieldName).(string)
	err:= addressDelete(ctx, req, name)
	if err!=nil {
		return nil, err
	}
		return delete(ctx, req, data)

}
