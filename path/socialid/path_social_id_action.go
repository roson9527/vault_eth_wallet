package socialid

import (
	"context"
	"errors"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/roson9527/vault_eth_wallet/modules"
	"github.com/roson9527/vault_eth_wallet/path/base"
	"github.com/roson9527/vault_eth_wallet/path/doc"
	"time"
)

func (pmgr *pathSocialID) actionPath(pattern string) *framework.Path {
	return &framework.Path{
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
				Callback: pmgr.readCallBack,
			},
			logical.UpdateOperation: &framework.PathOperation{
				Callback: pmgr.updateCallBack,
			},
			logical.DeleteOperation: &framework.PathOperation{
				Callback: pmgr.deleteCallBack,
			},
		},
		HelpSynopsis:    doc.PathReadSyn,
		HelpDescription: doc.PathReadDesc,
	}
}

func (pmgr *pathSocialID) deleteCallBack(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	namespace := data.Get(doc.FieldNameSpace).(string)
	if namespace != doc.NameSpaceGlobal {
		return nil, errors.New("only global namespace can be deleted")
	}

	app := data.Get(doc.FieldApp).(string)
	user := data.Get(doc.FieldUser).(string)
	// 获取目标钱包
	oldSocialId, err := pmgr.Storage.Social.Read(ctx, req, doc.NameSpaceGlobal, app, user)
	if err != nil {
		return nil, err
	}

	err = pmgr.Storage.Social.Delete(ctx, req, app, user)
	if err != nil {
		return nil, err
	}

	err = pmgr.Storage.Alias.Update(ctx, req, app2AType(app), user, oldSocialId.NameSpaces, []string{})
	if err != nil {
		return nil, err
	}

	return &logical.Response{
		Data: map[string]any{
			doc.FieldUser: user,
		},
	}, nil
}

func (pmgr *pathSocialID) updateCallBack(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	namespace := data.Get(doc.FieldNameSpace).(string)
	if namespace != doc.NameSpaceGlobal {
		return nil, errors.New("only global namespace can be updated")
	}

	app := data.Get(doc.FieldApp).(string)
	user := data.Get(doc.FieldUser).(string)
	// 获取目标钱包
	oldSocialId, err := pmgr.Storage.Social.Read(ctx, req, doc.NameSpaceGlobal, app, user)
	if err != nil {
		return nil, err
	}

	nameSpaces := data.Get(doc.FieldNameSpaces).([]string)
	wallet, err := pmgr.Storage.Social.Update(ctx, req, &modules.SocialID{
		NameSpaces: nameSpaces,
		UpdateTime: time.Now().Unix(),
	})
	if err != nil {
		return nil, err
	}

	err = pmgr.Storage.Alias.Update(ctx, req, app2AType(app), user, oldSocialId.NameSpaces, nameSpaces)
	if err != nil {
		return nil, err
	}

	return &logical.Response{
		Data: socialIDResponseData(wallet, false),
	}, nil
}

func (pmgr *pathSocialID) readCallBack(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	namespace := data.Get(doc.FieldNameSpace).(string)
	user := data.Get(doc.FieldUser).(string)
	app := data.Get(doc.FieldApp).(string)
	// 获取目标钱包
	socialId, err := pmgr.Storage.Social.Read(ctx, req, namespace, app, user)
	if err != nil {
		return nil, err
	}
	return &logical.Response{
		Data: socialIDResponseData(socialId, false),
	}, nil
}
