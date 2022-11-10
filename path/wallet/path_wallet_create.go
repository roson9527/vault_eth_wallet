package wallet

import (
	"context"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

func (pmgr *pathWallet) createWalletPath(pattern string) *framework.Path {
	return &framework.Path{
		Pattern: pattern,
		// 字段
		Fields: map[string]*framework.FieldSchema{
			fieldNameSpaces: {Type: framework.TypeCommaStringSlice, Default: []string{}},
			fieldPrivateKey: {Type: framework.TypeString},
			fieldPublicKey:  {Type: framework.TypeString},
			fieldAddress:    {Type: framework.TypeString},
		},
		// 执行的位置，有read，listWallet，createWallet，update
		Operations: map[logical.Operation]framework.OperationHandler{
			logical.UpdateOperation: &framework.PathOperation{
				Callback: pmgr.createCallBack,
			},
			logical.CreateOperation: &framework.PathOperation{
				Callback: pmgr.createCallBack,
			},
		},
		HelpSynopsis:    pathCreateSyn,
		HelpDescription: pathCreateDesc,
	}
}

func (pmgr *pathWallet) createCallBack(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	// 获取所有的钱包
	wallet, err := pmgr.walletStorage.createWallet(ctx, req, data)
	if err != nil {
		return nil, err
	}
	return &logical.Response{
		Data: walletResponseData(wallet, false),
	}, nil
}
