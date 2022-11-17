package wallet

import (
	"context"
	"errors"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/roson9527/vault_eth_wallet/modules"
	"github.com/roson9527/vault_eth_wallet/path/base"
	"time"
)

func (pmgr *pathWallet) readWalletPath(pattern string) *framework.Path {
	return &framework.Path{
		Pattern: pattern,
		// 字段
		Fields: map[string]*framework.FieldSchema{
			fieldNetwork:    {Type: framework.TypeString, Default: ""},
			fieldNameSpace:  {Type: framework.TypeString},
			fieldAddress:    {Type: framework.TypeString},
			fieldNameSpaces: {Type: framework.TypeCommaStringSlice, Default: []string{}},
		},
		ExistenceCheck: base.PathExistenceCheck,
		// 执行的位置，有read，listWallet，createWallet，update
		Operations: map[logical.Operation]framework.OperationHandler{
			logical.ReadOperation: &framework.PathOperation{
				Callback: pmgr.readCallBack,
			},
			logical.UpdateOperation: &framework.PathOperation{
				Callback: pmgr.updateCallBack,
			},
			logical.DeleteOperation: &framework.PathOperation{
				Callback: pmgr.deleteCallBack,
			},
		},
		HelpSynopsis:    pathReadSyn,
		HelpDescription: pathReadDesc,
	}
}

func (pmgr *pathWallet) deleteCallBack(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	namespace := data.Get(fieldNameSpace).(string)
	if namespace != NameSpaceGlobal {
		return nil, errors.New("only global namespace can be deleted")
	}

	address := data.Get(fieldAddress).(string)
	// 获取目标钱包
	oldWallet, err := pmgr.walletStorage.readWallet(ctx, req, NameSpaceGlobal, address)
	if err != nil {
		return nil, err
	}

	err = pmgr.walletStorage.deleteWallet(ctx, req, address)
	if err != nil {
		return nil, err
	}

	err = pmgr.walletStorage.updateAlias(ctx, req, address, oldWallet.NameSpaces, []string{})
	if err != nil {
		return nil, err
	}

	return &logical.Response{
		Data: map[string]any{
			fieldAddress: address,
		},
	}, nil
}

func (pmgr *pathWallet) updateCallBack(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	namespace := data.Get(fieldNameSpace).(string)
	if namespace != NameSpaceGlobal {
		return nil, errors.New("only global namespace can be updated")
	}

	overwrite := modules.Wallet{
		Address:    data.Get(fieldAddress).(string),
		NameSpaces: data.Get(fieldNameSpaces).([]string),
		Network:    data.Get(fieldNetwork).(string),
		UpdateTime: time.Now().Unix(),
	}
	// 获取目标钱包
	oldWallet, err := pmgr.walletStorage.readWallet(ctx, req, NameSpaceGlobal, overwrite.Address)
	if err != nil {
		return nil, err
	}

	wallet, err := pmgr.walletStorage.updateWallet(ctx, req, &overwrite)
	if err != nil {
		return nil, err
	}

	err = pmgr.walletStorage.updateAlias(ctx, req, overwrite.Address, oldWallet.NameSpaces, overwrite.NameSpaces)
	if err != nil {
		return nil, err
	}

	return &logical.Response{
		Data: walletResponseData(wallet, false),
	}, nil
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
