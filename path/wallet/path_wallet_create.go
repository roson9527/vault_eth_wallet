package wallet

import (
	"context"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/roson9527/vault_eth_wallet/modules"
	"github.com/roson9527/vault_eth_wallet/path/base"
)

func (pmgr *pathWallet) createWalletPath(pattern string) *framework.Path {
	return &framework.Path{
		Pattern: pattern,
		// 字段
		Fields: map[string]*framework.FieldSchema{
			fieldNameSpaces: {Type: framework.TypeCommaStringSlice, Default: []string{}},
			fieldPrivateKey: {Type: framework.TypeString},
			fieldAddress:    {Type: framework.TypeString},
			fieldNetwork:    {Type: framework.TypeString, Default: networkETH},
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
	overwrite := &modules.Wallet{}
	var err error

	privateKey := data.Get(fieldPrivateKey).(string)
	if privateKey != "" {
		overwrite, err = base.PrivateToWallet(privateKey)
		if err != nil {
			return nil, err
		}
	}
	overwrite.NameSpaces = data.Get(fieldNameSpaces).([]string)
	overwrite.Network = data.Get(fieldNetwork).(string)

	// 获取所有的钱包
	wallet, err := pmgr.walletStorage.createWallet(ctx, req, overwrite)
	if err != nil {
		return nil, err
	}

	err = pmgr.walletStorage.updateAlias(ctx, req, wallet.Address, []string{}, wallet.NameSpaces) // 更新别名
	if err != nil {
		return nil, err
	}

	return &logical.Response{
		Data: walletResponseData(wallet, false),
	}, nil
}
