package secp256k1

import (
	"context"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/roson9527/vault_eth_wallet/modules"
	"github.com/roson9527/vault_eth_wallet/path/base"
	"github.com/roson9527/vault_eth_wallet/path/doc"
	"time"
)

func (cb *callback) create(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	payload := &modules.WalletExtra{}
	var err error

	privateKey := data.Get(doc.FieldPrivateKey).(string)
	if privateKey != "" {
		payload, err = base.CryptoETH.PrivateToWallet(privateKey)
		if err != nil {
			return nil, err
		}
	} else {
		payload, err = base.CryptoETH.GenerateKey()
		if err != nil {
			return nil, err
		}
	}
	payload.NameSpaces = data.Get(doc.FieldNameSpaces).([]string)

	if err = payload.Extra.Decode(data.Get(doc.FieldExtra).(map[string]any)); err != nil {
		return nil, err
	}

	// 获取所有的钱包
	wallet, err := cb.Storage.create(ctx, req, payload)
	if err != nil {
		return nil, err
	}

	return &logical.Response{
		Data: walletResponseData(wallet, doc.ChainETH, false),
	}, nil
}

func (cb *callback) update(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	payload := modules.WalletExtra{
		Wallet: modules.Wallet{
			Address:    data.Get(doc.FieldAddress).(string),
			UpdateTime: time.Now().Unix(),
		},
		AddressAlias: data.Get(doc.FieldAddressAlias).(map[string]string),
		NameSpaces:   data.Get(doc.FieldNameSpaces).([]string),
		CryptoType:   doc.CryptoSECP256K1,
	}

	if err := payload.Extra.Decode(data.Get(doc.FieldExtra).(map[string]any)); err != nil {
		return nil, err
	}

	w, err := cb.Storage.update(ctx, req, &payload)
	if err != nil {
		return nil, err
	}

	return &logical.Response{
		Data: walletResponseData(w, doc.ChainETH, false),
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

//func (cb *callback) exportRaw(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
//	address := data.Get(doc.FieldAddress).(string)
//	wallet, err := cb.Storage.readRaw(ctx, req, address)
//	if err != nil {
//		return nil, err
//	}
//	return &logical.Response{
//		Data: walletResponseData(wallet, doc.ChainETH, true),
//	}, nil
//}

func (cb *callback) export(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	wallet, err := cb.readRaw(ctx, req, data)
	if err != nil {
		return nil, err
	}
	return &logical.Response{
		Data: walletResponseData(wallet, doc.ChainETH, true),
	}, nil
}

func (cb *callback) readRaw(ctx context.Context, req *logical.Request, data *framework.FieldData) (*modules.WalletExtra, error) {
	namespace := data.Get(doc.FieldNameSpace).(string)
	chain := data.Get(doc.FieldChain).(string)
	address := data.Get(doc.FieldAddress).(string)

	return cb.Storage.read(ctx, req, namespace, chain, address)
}

func (cb *callback) list(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	_ = data
	wallets, err := cb.Storage.Wallet.List(ctx, req, doc.CryptoSECP256K1)
	if err != nil {
		return nil, err
	}
	return &logical.Response{
		Data: map[string]any{
			doc.FieldKeys: wallets,
		},
	}, nil
}
