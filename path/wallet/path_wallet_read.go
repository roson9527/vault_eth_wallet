package wallet

import (
	"context"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/roson9527/vault_eth_wallet/path/base"
)

func (pmgr *pathWallet) readWalletPath(pattern string) *framework.Path {
	return &framework.Path{
		Pattern: pattern,
		// 字段
		Fields: map[string]*framework.FieldSchema{
			fieldNameSpace: {Type: framework.TypeString},
			fieldAddress:   {Type: framework.TypeString},
		},
		ExistenceCheck: base.PathExistenceCheck,
		// 执行的位置，有read，listWallet，createWallet，update
		Operations: map[logical.Operation]framework.OperationHandler{
			logical.ReadOperation: &framework.PathOperation{
				Callback: pmgr.readCallBack,
			},
		},
		HelpSynopsis:    pathReadSyn,
		HelpDescription: pathReadDesc,
	}
}

func (pmgr *pathWallet) readCallBack(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	namespace := data.Get(fieldNameSpace).(string)
	address := data.Get(fieldAddress).(string)
	// 获取目标钱包
	wallet, err := pmgr.walletStorage.readWallet(ctx, req, namespace, address)
	if err != nil {
		return nil, err
	}
	return &logical.Response{
		Data: walletResponseData(wallet, false),
	}, nil
}
