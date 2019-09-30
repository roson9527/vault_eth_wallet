package addresses

import (
	"context"
	"fmt"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/roson9527/vault_eth_wallet/modules"
)

func read(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	addr := data.Get(fieldAddress).(string)

	address, err := readByAddr(ctx, req, addr)
	if err != nil {
		return nil, err
	}

	return &logical.Response{
		Data: map[string]interface{}{
			fieldAccountName: address.Name,
		},
	}, nil
}

func readByAddr(ctx context.Context, req *logical.Request, addr string) (*modules.Address, error) {
	path := fmt.Sprintf("%s/%s", patternStr, addr)
	entry, err := req.Storage.Get(ctx, path)
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, nil
	}

	var address modules.Address
	err = entry.DecodeJSON(&address)
	if address.Name == "" {
		return nil, fmt.Errorf(errDeserializeAccount, path)
	}

	return &address, nil
}
