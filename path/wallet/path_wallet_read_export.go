package wallet

import (
	"context"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

func (pmgr *pathWallet) walletExportPath(pattern string) *framework.Path {
	return &framework.Path{
		Pattern: pattern,
		// 字段
		Fields: map[string]*framework.FieldSchema{
			fieldNameSpace: {Type: framework.TypeString, Required: true},
			fieldAddress:   {Type: framework.TypeString, Required: true},
		},
		// 执行的位置，有read，listWallet，createWallet，update
		Operations: map[logical.Operation]framework.OperationHandler{
			logical.ReadOperation: &framework.PathOperation{
				Callback: pmgr.exportCallBack,
			},
		},
		HelpSynopsis:    pathReadSyn,
		HelpDescription: pathReadSyn,
	}
}

func (pmgr *pathWallet) exportCallBack(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	namespace := data.Get(fieldNameSpace).(string)
	address := data.Get(fieldAddress).(string)
	// 获取目标钱包
	wallet, err := pmgr.walletStorage.readWallet(ctx, req, namespace, address)
	if err != nil {
		return nil, err
	}
	return &logical.Response{
		Data: walletResponseData(wallet, true),
	}, nil
}
