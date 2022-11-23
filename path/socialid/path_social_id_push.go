package socialid

import (
	"context"
	"errors"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/mitchellh/mapstructure"
	"github.com/roson9527/vault_eth_wallet/modules"
	"github.com/roson9527/vault_eth_wallet/path/doc"
)

func (h *handler) push(pattern string) *framework.Path {
	return &framework.Path{
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
				Callback: h.callback.create,
			},
			logical.CreateOperation: &framework.PathOperation{
				Callback: h.callback.create,
			},
		},
		HelpSynopsis:    doc.PathCreateSyn,
		HelpDescription: doc.PathCreateDesc,
	}
}

func (cb *callback) create(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	var payload modules.SocialID
	var err error

	err = mapstructure.Decode(data.Get(doc.FieldSocialID).(map[string]any), &payload)
	if err != nil {
		return nil, err
	}

	app := data.Get(doc.FieldApp).(string)
	user := data.Get(doc.FieldUser).(string)
	if app != payload.App {
		return nil, errors.New("app not match")
	}
	// 获取所有的钱包
	ret, err := cb.Storage.Social.Create(ctx, req, user, &payload)
	if err != nil {
		return nil, err
	}

	err = cb.Storage.Alias.Update(ctx, req, app2AType(app), user, []string{}, ret.NameSpaces) // 更新别名
	if err != nil {
		return nil, err
	}

	return &logical.Response{
		Data: socialIDResponseData(ret, false),
	}, nil
}
