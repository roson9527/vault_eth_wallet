package socialid

import (
	"context"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/roson9527/vault_eth_wallet/modules"
	"github.com/roson9527/vault_eth_wallet/path/doc"
)

func (s *storageEx) put(ctx context.Context, req *logical.Request, app, user string, socialID *modules.SocialID) error {
	// 获取目标钱包
	oldWallet, err := s.Social.Read(ctx, req, app, user)
	if err != nil {
		oldWallet = &modules.SocialID{
			NameSpaces: make([]string, 0),
		}
	}

	err = s.Social.Put(ctx, req, app, user, socialID)
	if err != nil {
		return err
	}

	// 更新Alias
	srcPath := s.Social.SrcPath(app, user)
	err = s.Alias.Update(ctx, req, aliasType(app), user, srcPath, oldWallet.NameSpaces, socialID.NameSpaces)
	return err
}

func (s *storageEx) delete(ctx context.Context, req *logical.Request, app, address string) error {
	// 获取目标钱包
	oldWallet, err := s.Social.Read(ctx, req, app, address)
	if err != nil {
		return err
	}

	err = s.Social.Delete(ctx, req, app, address)
	if err != nil {
		return err
	}

	err = s.Alias.Delete(ctx, req, aliasType(app), address, oldWallet.NameSpaces)
	return err
}

func (s *storageEx) read(ctx context.Context, req *logical.Request, ns, app, address string) (*modules.SocialID, error) {
	path, _ := s.Alias.ReadSrcPath(ctx, req, ns, aliasType(app), address)
	w, err := s.Social.ReadFromPath(ctx, req, path)
	if err != nil {
		return nil, err
	}
	return w, nil
}

func (s *storageEx) list(ctx context.Context, req *logical.Request, ns, app string) ([]string, error) {
	if ns == doc.NameSpaceGlobal && app == "" {
		return s.Social.List(ctx, req, app)
	}
	return s.Alias.List(ctx, req, ns, aliasType(app))
}
