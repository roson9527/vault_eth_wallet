package addresses

import (
	"context"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

func pathList(pattern string) *framework.Path {
	return &framework.Path{
		Pattern: pattern,
		// 字段
		Fields: nil,
		// 执行的位置，有read，list，create，update
		Operations: map[logical.Operation]framework.OperationHandler{
			logical.ListOperation: &framework.PathOperation{
				Callback: listAll,
			},
		},
		HelpSynopsis:    pathListSyn,
		HelpDescription: pathListDesc,
	}
}

func listAll(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	val, err := req.Storage.List(ctx, patternStr)
	if err != nil {
		return nil, err
	}
	return logical.ListResponse(val), nil
}
