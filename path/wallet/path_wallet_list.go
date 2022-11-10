package wallet

import (
	"context"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

func (pmgr *pathWallet) listWalletPath(pattern string) *framework.Path {
	return &framework.Path{
		Pattern: pattern,
		// 字段
		Fields: map[string]*framework.FieldSchema{
			fieldNameSpace: {
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
		HelpSynopsis:    pathListSyn,
		HelpDescription: pathListDesc,
	}
}

func (pmgr *pathWallet) listCallBack(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	namespace := data.Get(fieldNameSpace).(string)

	// 获取所有的钱包
	_, out, err := pmgr.walletStorage.listWallet(ctx, req, namespace)
	hclog.Default().Info("listWallet", "namespace", namespace, "length", len(out))
	if err != nil {
		return nil, err
	}

	return logical.ListResponse(out), nil
}
