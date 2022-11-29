package socialid

import (
	"fmt"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/roson9527/vault_eth_wallet/path/doc"
	"github.com/roson9527/vault_eth_wallet/path/storage"
)

func (h *handler) admin(app string) []*framework.Path {
	pattern := fmt.Sprintf(storage.PatternSocialID,
		framework.GenericNameRegex(doc.FieldNameSpace), app,
		framework.GenericNameRegex(doc.FieldUser))
	return []*framework.Path{
		// EXPORT
		{
			Pattern: pattern + doc.PathSubExport,
			// 字段
			Fields: map[string]*framework.FieldSchema{
				doc.FieldNameSpace: {Type: framework.TypeString},
				doc.FieldUser:      {Type: framework.TypeString},
				doc.FieldApp:       {Type: framework.TypeString, Required: true, Default: app},
			},
			// 执行的位置，有read，listWallet，createWallet，update
			Operations: map[logical.Operation]framework.OperationHandler{
				logical.ReadOperation: &framework.PathOperation{
					Callback: h.callback.export,
				},
			},
			HelpSynopsis:    doc.PathReadSyn,
			HelpDescription: doc.PathReadSyn,
		},

		// PUT + DELETE
		{
			Pattern: fmt.Sprintf(storage.PatternSocialID,
				doc.NameSpaceGlobal, app, framework.GenericNameRegex(doc.FieldUser)),
			// 字段
			Fields: map[string]*framework.FieldSchema{
				doc.FieldNameSpace: {Type: framework.TypeString, Required: true, Default: doc.NameSpaceGlobal},
				doc.FieldApp:       {Type: framework.TypeString, Required: true, Default: app},
				doc.FieldUser:      {Type: framework.TypeString, Required: true},
				doc.FieldSocialID:  {Type: framework.TypeMap, Default: map[string]any{}},
			},
			// 执行的位置，有read，listWallet，createWallet，update
			Operations: map[logical.Operation]framework.OperationHandler{
				logical.ReadOperation: &framework.PathOperation{
					Callback: h.callback.read,
				},
				logical.UpdateOperation: &framework.PathOperation{
					Callback: h.callback.put,
				},
				logical.CreateOperation: &framework.PathOperation{
					Callback: h.callback.put,
				},
				logical.DeleteOperation: &framework.PathOperation{
					Callback: h.callback.delete,
				},
			},
			HelpSynopsis:    doc.PathCreateSyn,
			HelpDescription: doc.PathCreateDesc,
		},
	}
}
