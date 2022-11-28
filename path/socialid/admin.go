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
		{
			Pattern: pattern + doc.PathSubExport,
			// 字段
			Fields: map[string]*framework.FieldSchema{
				doc.FieldNameSpace: {Type: framework.TypeString},
				doc.FieldUser:      {Type: framework.TypeString},
				doc.FieldApp:       {Type: framework.TypeString},
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
		{
			Pattern: fmt.Sprintf(storage.PatternSocialID,
				framework.GenericNameRegex(doc.FieldNameSpace), app, "?"),
			// 字段
			Fields: map[string]*framework.FieldSchema{
				doc.FieldNameSpace: {
					Type:        framework.TypeString,
					Description: "namespace",
					Required:    true,
				},
				doc.FieldApp: {
					Type:        framework.TypeString,
					Description: "app",
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
		},
		{
			Pattern: pattern,
			// 字段
			Fields: map[string]*framework.FieldSchema{
				doc.FieldNameSpace: {Type: framework.TypeString, Required: true},
				doc.FieldUser:      {Type: framework.TypeString, Required: true},
				doc.FieldSocialID:  {Type: framework.TypeMap, Default: map[string]any{}},
			},
			// 执行的位置，有read，listWallet，createWallet，update
			Operations: map[logical.Operation]framework.OperationHandler{
				logical.UpdateOperation: &framework.PathOperation{
					Callback: h.callback.push,
				},
				logical.CreateOperation: &framework.PathOperation{
					Callback: h.callback.push,
				},
			},
			HelpSynopsis:    doc.PathCreateSyn,
			HelpDescription: doc.PathCreateDesc,
		},
	}
}
