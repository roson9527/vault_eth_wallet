package text

import (
	"fmt"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/roson9527/vault_eth_wallet/path/doc"
	"github.com/roson9527/vault_eth_wallet/path/storage"
)

func (h *handler) action() []*framework.Path {
	return []*framework.Path{
		// read
		{
			Pattern: fmt.Sprintf(storage.PatternWalletByChain, framework.GenericNameRegex(doc.FieldNameSpace),
				doc.CryptoTEXT, framework.GenericNameRegex(doc.FieldChain),
				framework.GenericNameRegex(doc.FieldAddress)),
			Fields: map[string]*framework.FieldSchema{
				doc.FieldNameSpace: {Type: framework.TypeString, Required: true},
				doc.FieldChain:     {Type: framework.TypeString, Required: true},
				doc.FieldAddress:   {Type: framework.TypeString, Required: true},
			},
			Operations: map[logical.Operation]framework.OperationHandler{
				logical.ReadOperation: &framework.PathOperation{
					Callback: h.callback.read,
				},
			},
		},
	}
}
