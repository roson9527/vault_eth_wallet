package modules

import (
	"context"
	"fmt"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

type EthWalletBackend struct {
	*framework.Backend
	store map[string][]byte
}

func (b *EthWalletBackend) pathExistenceCheck(
	ctx context.Context, req *logical.Request, data *framework.FieldData) (bool, error) {

	out, err := req.Storage.Get(ctx, req.Path)
	if err != nil {
		return false, fmt.Errorf("existence check failed: %v", err)
	}

	return out != nil, nil
}
