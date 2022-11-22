package wallet

import (
	"context"
	"errors"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/roson9527/vault_eth_wallet/modules"
	"github.com/roson9527/vault_eth_wallet/path/base"
	"github.com/roson9527/vault_eth_wallet/path/doc"
	"time"
)

func (pmgr *pathWallet) walletActionPath(pattern string) *framework.Path {
	return &framework.Path{
		Pattern: pattern,
		// 字段
		Fields: map[string]*framework.FieldSchema{
			doc.FieldNetwork:    {Type: framework.TypeString, Default: ""},
			doc.FieldNameSpace:  {Type: framework.TypeString},
			doc.FieldAddress:    {Type: framework.TypeString},
			doc.FieldNameSpaces: {Type: framework.TypeCommaStringSlice, Default: []string{}},
			doc.FieldExtra:      {Type: framework.TypeMap, Default: map[string]any{}},
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
		HelpSynopsis:    doc.PathReadSyn,
		HelpDescription: doc.PathReadDesc,
	}
}

func (pmgr *pathWallet) deleteCallBack(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	namespace := data.Get(doc.FieldNameSpace).(string)
	if namespace != doc.NameSpaceGlobal {
		return nil, errors.New("only global namespace can be deleted")
	}

	address := data.Get(doc.FieldAddress).(string)
	// 获取目标钱包
	oldWallet, err := pmgr.Storage.Wallet.Read(ctx, req, doc.NameSpaceGlobal, address)
	if err != nil {
		return nil, err
	}

	err = pmgr.Storage.Wallet.Delete(ctx, req, address)
	if err != nil {
		return nil, err
	}

	err = pmgr.Storage.Alias.Update(ctx, req, address, oldWallet.NameSpaces, []string{})
	if err != nil {
		return nil, err
	}

	return &logical.Response{
		Data: map[string]any{
			doc.FieldAddress: address,
		},
	}, nil
}

func (pmgr *pathWallet) updateCallBack(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	namespace := data.Get(doc.FieldNameSpace).(string)
	if namespace != doc.NameSpaceGlobal {
		return nil, errors.New("only global namespace can be updated")
	}

	overwrite := modules.Wallet{
		Address:    data.Get(doc.FieldAddress).(string),
		NameSpaces: data.Get(doc.FieldNameSpaces).([]string),
		Network:    data.Get(doc.FieldNetwork).(string),
		UpdateTime: time.Now().Unix(),
	}

	if err := overwrite.DecodeExtra(data.Get(doc.FieldExtra).(map[string]any)); err != nil {
		return nil, err
	}

	// 获取目标钱包
	oldWallet, err := pmgr.Storage.Wallet.Read(ctx, req, doc.NameSpaceGlobal, overwrite.Address)
	if err != nil {
		return nil, err
	}

	wallet, err := pmgr.Storage.Wallet.Update(ctx, req, &overwrite)
	if err != nil {
		return nil, err
	}

	err = pmgr.Storage.Alias.Update(ctx, req, overwrite.Address, oldWallet.NameSpaces, overwrite.NameSpaces)
	if err != nil {
		return nil, err
	}

	return &logical.Response{
		Data: walletResponseData(wallet, false),
	}, nil
}

func (pmgr *pathWallet) readCallBack(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	namespace := data.Get(doc.FieldNameSpace).(string)
	address := data.Get(doc.FieldAddress).(string)
	// 获取目标钱包
	wallet, err := pmgr.Storage.Wallet.Read(ctx, req, namespace, address)
	if err != nil {
		return nil, err
	}
	return &logical.Response{
		Data: walletResponseData(wallet, false),
	}, nil
}
