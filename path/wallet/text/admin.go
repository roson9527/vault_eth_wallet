package text

import (
	"fmt"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/roson9527/vault_eth_wallet/path/doc"
	"github.com/roson9527/vault_eth_wallet/path/storage"
)

func (h *handler) admin() []*framework.Path {
	return []*framework.Path{
		{
			// update + delete
			Pattern: fmt.Sprintf(storage.PatternWallet, doc.NameSpaceGlobal, doc.CryptoTEXT,
				framework.GenericNameRegex(doc.FieldAddress)),
			Fields: map[string]*framework.FieldSchema{
				doc.FieldAddress:      {Type: framework.TypeString, Required: true},
				doc.FieldNameSpaces:   {Type: framework.TypeCommaStringSlice, Default: []string{}},
				doc.FieldAddressAlias: {Type: framework.TypeKVPairs, Default: map[string]string{}},
				doc.FieldExtra:        {Type: framework.TypeMap, Default: map[string]any{}},
			},
			Operations: map[logical.Operation]framework.OperationHandler{
				logical.UpdateOperation: &framework.PathOperation{
					Callback: h.callback.update,
				},
				logical.CreateOperation: &framework.PathOperation{
					Callback: h.callback.update,
				},
				logical.DeleteOperation: &framework.PathOperation{
					Callback: h.callback.delete,
				},
			},
		},
		{
			// export
			Pattern: fmt.Sprintf(storage.PatternWalletByChain, framework.GenericNameRegex(doc.FieldNameSpace),
				doc.CryptoTEXT, framework.GenericNameRegex(doc.FieldChain),
				framework.GenericNameRegex(doc.FieldAddress)) + doc.PathSubExport,
			Fields: map[string]*framework.FieldSchema{
				doc.FieldNameSpace: {Type: framework.TypeString, Required: true},
				doc.FieldChain:     {Type: framework.TypeString, Required: true},
				doc.FieldAddress:   {Type: framework.TypeString, Required: true},
			},
			Operations: map[logical.Operation]framework.OperationHandler{
				logical.ReadOperation: &framework.PathOperation{
					Callback: h.callback.export,
				},
			},
		},
		{
			// list
			Pattern: fmt.Sprintf(storage.PatternWalletAdmin, doc.NameSpaceGlobal, doc.CryptoTEXT+"/?"),
			Fields:  map[string]*framework.FieldSchema{
				// 不需要任何参数
			},
			Operations: map[logical.Operation]framework.OperationHandler{
				logical.ListOperation: &framework.PathOperation{
					Callback: h.callback.list,
				},
			},
		},
		{
			// list alias
			Pattern: fmt.Sprintf(storage.PatternWalletByChain, framework.GenericNameRegex(doc.FieldNameSpace),
				doc.CryptoTEXT, framework.GenericNameRegex(doc.FieldChain), "?"),
			Fields: map[string]*framework.FieldSchema{
				// 不需要任何参数
				doc.FieldNameSpace: {Type: framework.TypeString, Required: true},
				doc.FieldChain:     {Type: framework.TypeString, Required: true},
			},
			Operations: map[logical.Operation]framework.OperationHandler{
				logical.ListOperation: &framework.PathOperation{
					Callback: h.callback.listAlias,
				},
			},
		},
	}
}
