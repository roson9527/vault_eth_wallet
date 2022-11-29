package socialid

import (
	"context"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/roson9527/vault_eth_wallet/path/doc"
)

func (cb *callback) read(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	ns := data.Get(doc.FieldNameSpace).(string)
	app := data.Get(doc.FieldApp).(string)
	user := data.Get(doc.FieldUser).(string)

	// 获取目标钱包
	socialId, err := cb.Storage.read(ctx, req, ns, app, user) //.Social.Read(ctx, req, namespace, app, user)
	if err != nil {
		return nil, err
	}
	return &logical.Response{
		Data: socialIDResponseData(socialId, false),
	}, nil
}

func (cb *callback) list(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	namespace := data.Get(doc.FieldNameSpace).(string)
	app := data.Get(doc.FieldApp).(string)

	list, err := cb.Storage.list(ctx, req, namespace, app)
	if err != nil {
		return nil, err
	}

	return logical.ListResponse(list), nil
}
