package storage

import (
	"context"
	"errors"
	"fmt"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/roson9527/vault_eth_wallet/modules"
	"github.com/roson9527/vault_eth_wallet/path/doc"
)

const (
	// PatternSocialID is the pattern for the storage path of a social id
	/*

		[global | namespace]/social/[app]/[user]

		e.g. [global | namespace]/social/discord/123456789
	*/
	PatternSocialID = "%s/social/%s/%s"
)

type socialIDStorage struct {
}

func newSocialIDStorage() *socialIDStorage {
	return &socialIDStorage{}
}

func (as *socialIDStorage) Read(ctx context.Context, req *logical.Request, namespace, app, user string) (*modules.SocialID, error) {
	path := fmt.Sprintf(PatternSocialID, doc.NameSpaceGlobal, app, user)
	entry, err := req.Storage.Get(ctx, path)

	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, errors.New(fmt.Sprintf("[%s] social id not found", path))
	}

	var ety modules.SocialID
	err = entry.DecodeJSON(&ety)

	if ety.User == "" {
		return nil, fmt.Errorf("failed to deserialize user at %s", path)
	}

	if namespace == doc.NameSpaceGlobal {
		return &ety, nil
	}

	if ety.NameSpaces == nil {
		return nil, fmt.Errorf("discord %s does not have namespace %s", ety.User, namespace)
	}

	for _, ns := range ety.NameSpaces {
		if ns == namespace {
			return &ety, nil
		}
	}

	return nil, fmt.Errorf("not support namespace %s", namespace)
}

func (as *socialIDStorage) Update(ctx context.Context, req *logical.Request, payload *modules.SocialID) (*modules.SocialID, error) {
	oldSocialId, err := as.Read(ctx, req, doc.NameSpaceGlobal, payload.App, payload.User)
	if err != nil {
		return nil, err
	}

	oldSocialId.UpdateTime = payload.UpdateTime
	oldSocialId.NameSpaces = payload.NameSpaces

	path := fmt.Sprintf(PatternSocialID, doc.NameSpaceGlobal, payload.App, payload.User)
	entry, err := logical.StorageEntryJSON(path, oldSocialId)
	if err != nil {
		return nil, err
	}

	err = req.Storage.Put(ctx, entry)
	if err != nil {
		return nil, err
	}

	return payload, nil
}

func (as *socialIDStorage) Create(ctx context.Context, req *logical.Request, user string, payload *modules.SocialID) (*modules.SocialID, error) {
	if payload.User == "" {
		return nil, errors.New("socialID.user is empty")
	}
	if payload.App == "" {
		return nil, errors.New("socialID.app is empty")
	}

	insertPath := fmt.Sprintf(PatternSocialID, doc.NameSpaceGlobal, payload.App, user)
	entry, err := logical.StorageEntryJSON(insertPath, payload)
	if err != nil {
		return nil, err
	}

	err = req.Storage.Put(ctx, entry)
	if err != nil {
		return nil, err
	}

	return payload, nil
}

func (as *socialIDStorage) List(ctx context.Context, req *logical.Request, app string) ([]string, error) {
	return req.Storage.List(ctx, fmt.Sprintf(PatternSocialID, doc.NameSpaceGlobal, app, ""))
}

func (as *socialIDStorage) Delete(ctx context.Context, req *logical.Request, app, user string) error {
	return req.Storage.Delete(ctx, fmt.Sprintf(PatternSocialID, doc.NameSpaceGlobal, app, user))
}
