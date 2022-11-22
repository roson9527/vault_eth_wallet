package wallet

import (
	"context"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/roson9527/vault_eth_wallet/path/doc"
)

func (pmgr *pathWallet) listWalletPath(pattern string) *framework.Path {
	return &framework.Path{
		Pattern: pattern,
		// 字段
		Fields: map[string]*framework.FieldSchema{
			doc.FieldNameSpace: {
				Type:        framework.TypeString,
				Description: "Namespace",
				Required:    true,
			},
		},
		// 执行的位置，有read，listWallet，createWallet，update
		Operations: map[logical.Operation]framework.OperationHandler{
			logical.ListOperation: &framework.PathOperation{
				Callback: pmgr.listCallBack,
			},
		},
		HelpSynopsis:    doc.PathListSyn,
		HelpDescription: doc.PathListDesc,
	}
}

func (pmgr *pathWallet) listCallBack(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	namespace := data.Get(doc.FieldNameSpace).(string)

	// 获取所有的钱包
	out, err := pmgr.Storage.Wallet.List(ctx, req, namespace)
	hclog.Default().Info("listWallet", "namespace", namespace, "length", len(out))
	if err != nil {
		return nil, err
	}

	return logical.ListResponse(out), nil
}
