package wallet

import (
	"context"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/roson9527/vault_eth_wallet/path/doc"
)

func (h *handler) list(pattern string) *framework.Path {
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
				Callback: h.callback.list,
			},
		},
		HelpSynopsis:    doc.PathListSyn,
		HelpDescription: doc.PathListDesc,
	}
}

func (cb *callback) list(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	namespace := data.Get(doc.FieldNameSpace).(string)

	var out []string
	var err error
	if namespace == doc.NameSpaceGlobal {
		out, err = cb.Storage.Wallet.List(ctx, req)
	} else {
		out, err = cb.Storage.Alias.List(ctx, req, namespace, doc.AliasWallet)
	}
	if err != nil {
		return nil, err
	}

	hclog.Default().Info("wallet:list", "namespace", namespace, "length", len(out))

	return logical.ListResponse(out), nil
}
