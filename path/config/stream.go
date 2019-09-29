package config

import (
	"context"
	"fmt"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/roson9527/vault_eth_wallet/modules"
)

func write(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	rpcURL := data.Get(fieldRpcUrl).(string)
	chainId := data.Get(fieldChainId).(string)

	conf := modules.Config{
		RPC:     rpcURL,
		ChainID: chainId,
	}

	entry, err := logical.StorageEntryJSON(storageKey, conf)
	if err != nil {
		return nil, err
	}

	if err := req.Storage.Put(ctx, entry); err != nil {
		return nil, err
	}
	// Return the secret
	return &logical.Response{
		Data: map[string]interface{}{
			fieldChainId: conf.ChainID,
			fieldRpcUrl:  conf.RPC,
		},
	}, nil
}

func read(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	conf, err := readFromStorage(ctx, req.Storage)
	if err != nil {
		return nil, err
	}

	if conf == nil {
		return nil, nil
	}

	// Return the secret
	return &logical.Response{
		Data: map[string]interface{}{
			fieldChainId: conf.ChainID,
			fieldRpcUrl:  conf.RPC,
		},
	}, nil
}

func readFromStorage(ctx context.Context, s logical.Storage) (*modules.Config, error) {
	entry, err := s.Get(ctx, storageKey)
	if err != nil {
		return nil, err
	}

	if entry == nil {
		return nil, fmt.Errorf(errNotConfiguredEthBackend)
	}

	var result modules.Config
	if err = entry.DecodeJSON(&result); err != nil {
		return nil, err //fmt.Errorf("error reading configuration: %s", err)
	}

	return &result, nil
}
