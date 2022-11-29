package socialid

import (
	"context"
	"errors"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/mitchellh/mapstructure"
	"github.com/roson9527/vault_eth_wallet/modules"
	"github.com/roson9527/vault_eth_wallet/path/doc"
	"time"
)

func (cb *callback) export(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	ns := data.Get(doc.FieldNameSpace).(string)
	app := data.Get(doc.FieldApp).(string)
	user := data.Get(doc.FieldUser).(string)
	// 获取目标钱包
	socialId, err := cb.Storage.read(ctx, req, ns, app, user)
	if err != nil {
		return nil, err
	}
	return &logical.Response{
		Data: socialIDResponseData(socialId, true),
	}, nil
}

func (cb *callback) put(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	var payload modules.SocialID
	var err error

	err = mapstructure.Decode(data.Get(doc.FieldSocialID).(map[string]any), &payload)
	if err != nil {
		return nil, err
	}

	app := data.Get(doc.FieldApp).(string)
	user := data.Get(doc.FieldUser).(string)
	if app != payload.App && len(app) > 0 {
		return nil, errors.New("app not match")
	}
	payload.UpdateTime = time.Now().Unix()
	// 获取所有的钱包
	err = cb.Storage.put(ctx, req, app, user, &payload)
	if err != nil {
		return nil, err
	}

	return &logical.Response{
		Data: socialIDResponseData(&payload, false),
	}, nil
}

func (cb *callback) delete(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	app := data.Get(doc.FieldApp).(string)
	user := data.Get(doc.FieldUser).(string)
	// 获取目标钱包
	oldSocialId, err := cb.Storage.read(ctx, req, doc.NameSpaceGlobal, app, user)
	if err != nil {
		return nil, err
	}

	err = cb.Storage.Social.Delete(ctx, req, app, user)
	if err != nil {
		return nil, err
	}

	err = cb.Storage.Alias.Delete(ctx, req, aliasType(app), user, oldSocialId.NameSpaces)
	if err != nil {
		return nil, err
	}

	return &logical.Response{
		Data: map[string]any{
			doc.FieldUser: user,
		},
	}, nil
}
