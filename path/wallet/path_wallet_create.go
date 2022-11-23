package wallet

import (
	"context"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/roson9527/vault_eth_wallet/modules"
	"github.com/roson9527/vault_eth_wallet/path/base"
	"github.com/roson9527/vault_eth_wallet/path/doc"
)

func (h *handler) create(pattern string) *framework.Path {
	return &framework.Path{
		Pattern: pattern,
		// 字段
		Fields: map[string]*framework.FieldSchema{
			doc.FieldNameSpaces: {Type: framework.TypeCommaStringSlice, Default: []string{}},
			doc.FieldPrivateKey: {Type: framework.TypeString},
			doc.FieldAddress:    {Type: framework.TypeString},
			doc.FieldNetwork:    {Type: framework.TypeString, Default: doc.NetworkETH},
			doc.FieldExtra:      {Type: framework.TypeMap, Default: map[string]any{}},
		},
		// 执行的位置，有read，listWallet，createWallet，update
		Operations: map[logical.Operation]framework.OperationHandler{
			logical.UpdateOperation: &framework.PathOperation{
				Callback: h.callback.create,
			},
			logical.CreateOperation: &framework.PathOperation{
				Callback: h.callback.create,
			},
		},
		HelpSynopsis:    doc.PathCreateSyn,
		HelpDescription: doc.PathCreateDesc,
	}
}

func (cb *callback) create(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	overwrite := &modules.Wallet{}
	var err error

	privateKey := data.Get(doc.FieldPrivateKey).(string)
	if privateKey != "" {
		overwrite, err = base.CryptoETH.PrivateToWallet(privateKey)
		if err != nil {
			return nil, err
		}
	}
	overwrite.NameSpaces = data.Get(doc.FieldNameSpaces).([]string)
	overwrite.Network = data.Get(doc.FieldNetwork).(string)
	if err = overwrite.Extra.Decode(data.Get(doc.FieldExtra).(map[string]any)); err != nil {
		return nil, err
	}

	// 获取所有的钱包
	wallet, err := cb.Storage.Wallet.Create(ctx, req, overwrite)
	if err != nil {
		return nil, err
	}

	err = cb.Storage.Alias.Update(ctx, req, doc.AliasWallet, wallet.Address, []string{}, wallet.NameSpaces) // 更新别名
	if err != nil {
		return nil, err
	}

	return &logical.Response{
		Data: walletResponseData(wallet, false),
	}, nil
}
