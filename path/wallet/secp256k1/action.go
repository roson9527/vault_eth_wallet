package secp256k1

import (
	"fmt"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/roson9527/vault_eth_wallet/path/doc"
	"github.com/roson9527/vault_eth_wallet/path/storage"
)

func (h *handler) action() []*framework.Path {
	pattern := fmt.Sprintf(storage.PatternWalletByChain, framework.GenericNameRegex(doc.FieldNameSpace),
		doc.CryptoSECP256K1, framework.GenericNameRegex(doc.FieldChain), framework.GenericNameRegex(doc.FieldAddress))

	return []*framework.Path{
		// read
		{
			Pattern: pattern,
			Fields: map[string]*framework.FieldSchema{
				doc.FieldNameSpace: {Type: framework.TypeString, Required: true},
				doc.FieldAddress:   {Type: framework.TypeString, Required: true},
				doc.FieldChain:     {Type: framework.TypeString, Required: true},
			},
			Operations: map[logical.Operation]framework.OperationHandler{
				logical.ReadOperation: &framework.PathOperation{
					Callback: h.callback.read,
				},
			},
		},
		// sign_tx
		{
			Pattern: pattern + doc.PathSubSignTx,
			Fields: map[string]*framework.FieldSchema{
				doc.FieldNameSpace: {Type: framework.TypeString, Required: true},
				doc.FieldAddress:   {Type: framework.TypeString, Required: true},
				doc.FieldChain:     {Type: framework.TypeString, Required: true},
				doc.FieldTxBinary:  {Type: framework.TypeString, Required: true},
			},
			Operations: map[logical.Operation]framework.OperationHandler{
				logical.UpdateOperation: &framework.PathOperation{
					Callback: h.callback.signTx,
				},
			},
		},
		// list
		{
			Pattern: fmt.Sprintf(storage.PatternWallet, framework.GenericNameRegex(doc.FieldNameSpace),
				doc.CryptoSECP256K1, framework.GenericNameRegex(doc.FieldChain)+"/?"),
			Fields: map[string]*framework.FieldSchema{
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
