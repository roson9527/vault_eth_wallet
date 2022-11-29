package text

import (
	"context"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/mitchellh/mapstructure"
	"github.com/roson9527/vault_eth_wallet/modules"
	"github.com/roson9527/vault_eth_wallet/path/doc"
	"time"
)

func (cb *callback) put(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	var payload modules.WalletExtra
	err := mapstructure.Decode(data.Raw, &payload)
	if err != nil {
		return nil, err
	}

	if payload.Address != data.Get(doc.FieldAddress).(string) {
		return nil, ErrAddressNotMatch
	}

	payload.CryptoType = doc.CryptoTEXT
	payload.UpdateTime = time.Now().Unix()

	w, err := cb.Storage.put(ctx, req, &payload)
	if err != nil {
		return nil, err
	}

	return &logical.Response{
		Data: walletResponseData(w, "", false),
	}, nil
}

func (cb *callback) delete(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	address := data.Get(doc.FieldAddress).(string)
	err := cb.Storage.delete(ctx, req, address)
	if err != nil {
		return nil, err
	}

	return &logical.Response{
		Data: map[string]any{
			doc.FieldAddress: address,
		},
	}, nil
}

func (cb *callback) export(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	wallet, err := cb.readRaw(ctx, req, data)
	if err != nil {
		return nil, err
	}
	chain := data.Get(doc.FieldChain).(string)
	return &logical.Response{
		Data: walletResponseData(wallet, chain, true),
	}, nil
}

func (cb *callback) readRaw(ctx context.Context, req *logical.Request, data *framework.FieldData) (*modules.WalletExtra, error) {
	namespace := data.Get(doc.FieldNameSpace).(string)
	chain := data.Get(doc.FieldChain).(string)
	address := data.Get(doc.FieldAddress).(string)

	return cb.Storage.read(ctx, req, namespace, chain, address)
}

func (cb *callback) list(ctx context.Context, req *logical.Request, _ *framework.FieldData) (*logical.Response, error) {
	wallets, err := cb.Storage.list(ctx, req, doc.NameSpaceGlobal, doc.CryptoTEXT)
	if err != nil {
		return nil, err
	}
	return &logical.Response{
		Data: map[string]any{
			doc.FieldKeys: wallets,
		},
	}, nil
}

func (cb *callback) listAlias(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	namespace := data.Get(doc.FieldNameSpace).(string)
	chain := data.Get(doc.FieldChain).(string)
	alias, err := cb.Storage.list(ctx, req, namespace, chain)
	if err != nil {
		return nil, err
	}
	return &logical.Response{
		Data: map[string]any{
			doc.FieldKeys: alias,
		},
	}, nil
}
