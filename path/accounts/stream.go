package accounts

import (
	"context"
	"fmt"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/roson9527/vault_eth_wallet/modules"
	"github.com/roson9527/vault_eth_wallet/path/base"
)

func create(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	name := data.Get(fieldName).(string)
	_ = name

	accountEty, err := base.GenerateKey()
	if err != nil {
		return nil, err
	}

	entry, err := logical.StorageEntryJSON(req.Path, accountEty)
	if err != nil {
		return nil, err
	}

	err = req.Storage.Put(ctx, entry)
	if err != nil {
		return nil, err
	}

	return &logical.Response{
		Data: map[string]interface{}{
			fieldAddress: accountEty.Address,
		},
	}, nil
}

func read(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	name := data.Get(fieldName).(string)
	acnt, err := ReadByName(ctx, req, name)
	if err != nil {
		return nil, err
	}

	return &logical.Response{
		Data: map[string]interface{}{
			fieldAddress: acnt.Address,
		},
	}, nil
}

func ReadByName(ctx context.Context, req *logical.Request, name string) (*modules.Account, error) {
	path := fmt.Sprintf("accounts/%s", name)
	entry, err := req.Storage.Get(ctx, path)
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, nil
	}

	var account modules.Account
	err = entry.DecodeJSON(&account)

	if account.Address == "" {
		return nil, fmt.Errorf("failed to deserialize account at %s", path)
	}

	return &account, nil
}

func delete(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	name := data.Get("name").(string)
	account, err := ReadByName(ctx, req, name)
	if err != nil {
		return nil, fmt.Errorf(errReadAccountByName)
	}
	if account == nil {
		return nil, nil
	}
	if err := req.Storage.Delete(ctx, req.Path); err != nil {
		return nil, err
	}

	return nil, nil
}
