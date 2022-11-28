package socialid

import (
	"context"
	"errors"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/mitchellh/mapstructure"
	"github.com/roson9527/vault_eth_wallet/modules"
	"github.com/roson9527/vault_eth_wallet/path/doc"
)

func (cb *callback) export(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	namespace := data.Get(doc.FieldNameSpace).(string)
	app := data.Get(doc.FieldApp).(string)
	user := data.Get(doc.FieldUser).(string)
	// 获取目标钱包
	socialId, err := cb.Storage.Social.Read(ctx, req, namespace, app, user)
	if err != nil {
		return nil, err
	}
	return &logical.Response{
		Data: socialIDResponseData(socialId, true),
	}, nil
}

func (cb *callback) list(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	namespace := data.Get(doc.FieldNameSpace).(string)
	app := data.Get(doc.FieldApp).(string)

	var out []string
	var err error
	if namespace == doc.NameSpaceGlobal {
		out, err = cb.Storage.Social.List(ctx, req, app)
	} else {
		out, err = cb.Storage.Alias.List(ctx, req, namespace, app)
	}
	if err != nil {
		return nil, err
	}

	hclog.Default().Info(app+":list", "namespace", namespace, "length", len(out))

	return logical.ListResponse(out), nil
}

func (cb *callback) push(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
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
