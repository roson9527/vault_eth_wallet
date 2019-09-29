package accounts

import (
	"context"
	"fmt"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/roson9527/vault_eth_wallet/modules"
)

// 简单更新数据
func addressWrite(ctx context.Context, req *logical.Request, addr, name string) error {
	addrEty := modules.Address{
		Name: name,
	}

	storageKey := fmt.Sprintf("%s/%s", patternAddressStr, addr)
	entry, err := logical.StorageEntryJSON(storageKey, addrEty)
	if err != nil {
		return err
	}

	return req.Storage.Put(ctx, entry)
}

func addressDelete(ctx context.Context, req *logical.Request, addr string) error {
	storageKey := fmt.Sprintf("%s/%s", patternAddressStr, addr)
	return req.Storage.Delete(ctx, storageKey)
}
