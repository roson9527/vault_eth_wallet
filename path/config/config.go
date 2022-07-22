package config

import (
	"context"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/roson9527/vault_eth_wallet/modules"
)

// 已经测试通过
// 可以进行简单的读取和写入
func Path() []*framework.Path {
	return []*framework.Path{
		pathConfig("config"),
	}
}

func Read(ctx context.Context, req *logical.Request) (*modules.Config, error) {
	conf, err := readFromStorage(ctx, req.Storage)
	if err != nil {
		return nil, err
	}

	return conf, nil
}
