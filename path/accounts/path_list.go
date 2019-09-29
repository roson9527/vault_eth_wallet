package accounts

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

//Pattern:         pattern,
//// 字段
//Fields:          nil,
//// 之行的位置，有read，list，create，update
//Operations:  map[logical.Operation]framework.OperationHandler{
//logical.ListOperation: &framework.PathOperation{
//Callback:    listAll,
//},
//},
//// 存在性检测，没有则会在update时先处理create
//ExistenceCheck:  nil,
//// 如果已经实现，则验证是否为路径提供了功能
//FeatureRequired: 0,
//// 是否弃用，一般不用管
////Deprecated:      false,
//HelpSynopsis: "List all the Ethereum accounts at a path",
//HelpDescription: `
//			All the Ethereum accounts will be listed.
//			`,
//// 用于生成 OPEN API 使用
//DisplayAttrs:    nil,
