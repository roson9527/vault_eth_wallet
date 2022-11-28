package socialid

import (
	"fmt"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/roson9527/vault_eth_wallet/path/base"
	"github.com/roson9527/vault_eth_wallet/path/doc"
	"github.com/roson9527/vault_eth_wallet/path/storage"
)

func (h *handler) action() []*framework.Path {
	pattern := fmt.Sprintf(storage.PatternSocialID, framework.GenericNameRegex(doc.FieldNameSpace),
		framework.GenericNameRegex(doc.FieldApp), framework.GenericNameRegex(doc.FieldAddress))

	return []*framework.Path{
		{
			Pattern: pattern,
			// 字段
			Fields: map[string]*framework.FieldSchema{
				doc.FieldNameSpace:  {Type: framework.TypeString},
				doc.FieldUser:       {Type: framework.TypeString},
				doc.FieldApp:        {Type: framework.TypeString},
				doc.FieldNameSpaces: {Type: framework.TypeCommaStringSlice, Default: []string{}},
			},
			ExistenceCheck: base.PathExistenceCheck,
			// 执行的位置，有read，listWallet，createWallet，update
			Operations: map[logical.Operation]framework.OperationHandler{
				logical.ReadOperation: &framework.PathOperation{
					Callback: h.callback.read,
				},
				logical.UpdateOperation: &framework.PathOperation{
					Callback: h.callback.update,
				},
				logical.DeleteOperation: &framework.PathOperation{
					Callback: h.callback.delete,
				},
			},
			HelpSynopsis:    doc.PathReadSyn,
			HelpDescription: doc.PathReadDesc,
		},
	}
}
