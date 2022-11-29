package eth

import (
	"fmt"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/roson9527/vault_eth_wallet/path/doc"
	"github.com/roson9527/vault_eth_wallet/path/storage"
)

func (h *handler) admin() []*framework.Path {
	pattern := fmt.Sprintf(storage.PatternPolicy,
		framework.GenericNameRegex(doc.FieldNameSpace),
		doc.CryptoSECP256K1,
		framework.GenericNameRegex(doc.FieldChain))

	return []*framework.Path{
		{
			Pattern: pattern,
			Fields: map[string]*framework.FieldSchema{
				doc.FieldNameSpace: {
					Type:        framework.TypeString,
					Description: "Namespace",
					Required:    true,
				},
				doc.FieldCryptoType: {
					Type:        framework.TypeString,
					Description: "crypto type",
					Default:     doc.CryptoSECP256K1,
				},
				doc.FieldChain: {
					Type:        framework.TypeString,
					Description: "Chain",
					Required:    true,
				},
				doc.FieldPolicyHCL: {Type: framework.TypeString, Default: ""},
				doc.FieldPolicy:    {Type: framework.TypeMap, Default: map[string]any{}},
			},

			Operations: map[logical.Operation]framework.OperationHandler{
				logical.UpdateOperation: &framework.PathOperation{
					Callback: h.callback.put,
				},
				logical.CreateOperation: &framework.PathOperation{
					Callback: h.callback.put,
				},
				logical.ReadOperation: &framework.PathOperation{
					Callback: h.callback.read,
				},
			},
			HelpSynopsis:    "",
			HelpDescription: "",
		},
	}
}
