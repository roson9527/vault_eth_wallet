package socialid

import (
	"context"
	"errors"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/roson9527/vault_eth_wallet/modules"
	"github.com/roson9527/vault_eth_wallet/path/doc"
	"time"
)

func (cb *callback) delete(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	namespace := data.Get(doc.FieldNameSpace).(string)
	if namespace != doc.NameSpaceGlobal {
		return nil, errors.New("only global namespace can be deleted")
	}

	app := data.Get(doc.FieldApp).(string)
	user := data.Get(doc.FieldUser).(string)
	// 获取目标钱包
	oldSocialId, err := cb.Storage.Social.Read(ctx, req, doc.NameSpaceGlobal, app, user)
	if err != nil {
		return nil, err
	}

	err = cb.Storage.Social.Delete(ctx, req, app, user)
	if err != nil {
		return nil, err
	}

	err = cb.Storage.Alias.Update(ctx, req, app2AType(app), user, "", oldSocialId.NameSpaces, []string{})
	if err != nil {
		return nil, err
	}

	return &logical.Response{
		Data: map[string]any{
			doc.FieldUser: user,
		},
	}, nil
}

func (cb *callback) update(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	namespace := data.Get(doc.FieldNameSpace).(string)
	if namespace != doc.NameSpaceGlobal {
		return nil, errors.New("only global namespace can be updated")
	}

	app := data.Get(doc.FieldApp).(string)
	user := data.Get(doc.FieldUser).(string)
	// 获取目标钱包
	oldSocialId, err := cb.Storage.Social.Read(ctx, req, doc.NameSpaceGlobal, app, user)
	if err != nil {
		return nil, err
	}

	nameSpaces := data.Get(doc.FieldNameSpaces).([]string)
	wallet, err := cb.Storage.Social.Update(ctx, req, &modules.SocialID{
		NameSpaces: nameSpaces,
		UpdateTime: time.Now().Unix(),
	})
	if err != nil {
		return nil, err
	}

	err = cb.Storage.Alias.Update(ctx, req, app2AType(app), user, oldSocialId.NameSpaces, nameSpaces)
	if err != nil {
		return nil, err
	}

	return &logical.Response{
		Data: socialIDResponseData(wallet, false),
	}, nil
}

func (cb *callback) read(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	namespace := data.Get(doc.FieldNameSpace).(string)
	user := data.Get(doc.FieldUser).(string)
	app := data.Get(doc.FieldApp).(string)
	// 获取目标钱包
	socialId, err := cb.Storage.Social.Read(ctx, req, namespace, app, user)
	if err != nil {
		return nil, err
	}
	return &logical.Response{
		Data: socialIDResponseData(socialId, false),
	}, nil
}
