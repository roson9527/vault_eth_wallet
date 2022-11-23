package wallet

import (
	"context"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/roson9527/vault_eth_wallet/path/doc"
)

func (h *handler) export(pattern string) *framework.Path {
	return &framework.Path{
		Pattern: pattern,
		// 字段
		Fields: map[string]*framework.FieldSchema{
			doc.FieldNameSpace: {Type: framework.TypeString, Required: true},
			doc.FieldAddress:   {Type: framework.TypeString, Required: true},
		},
		// 执行的位置，有read，listWallet，createWallet，update
		Operations: map[logical.Operation]framework.OperationHandler{
			logical.ReadOperation: &framework.PathOperation{
				Callback: h.callback.export,
			},
		},
		HelpSynopsis:    doc.PathReadSyn,
		HelpDescription: doc.PathReadSyn,
	}
}

func (cb *callback) export(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	namespace := data.Get(doc.FieldNameSpace).(string)
	address := data.Get(doc.FieldAddress).(string)
	// 获取目标钱包
	wallet, err := cb.Storage.Wallet.Read(ctx, req, namespace, address)
	if err != nil {
		return nil, err
	}
	return &logical.Response{
		Data: walletResponseData(wallet, true),
	}, nil
}
