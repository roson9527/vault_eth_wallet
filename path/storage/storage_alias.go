package storage

import (
	"context"
	"fmt"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/roson9527/vault_eth_wallet/modules"
)

const (
	// PatternAlias is the pattern for the alias storage path
	/* PatternAlias = "%s/alias/%s/%s"
	 * [namespace]/alias/[type]/[address]
	 * e.g. global/alias/discord/0x1234567890
	 * e.g. global/alias/wallet/0x1234567890
	 */
	PatternAlias = "%s/alias/%s/%s"
)

type aliasStorage struct {
}

func newAliasStorage() *aliasStorage {
	return &aliasStorage{}
}

func (as *aliasStorage) List(ctx context.Context, req *logical.Request, namespace, aType string) ([]string, error) {
	// 获取namespace下所有对应vType的alias
	return req.Storage.List(ctx, fmt.Sprintf(PatternAlias, namespace, aType, ""))
}

func (as *aliasStorage) Update(ctx context.Context, req *logical.Request, aType, address string, oldNS, newNS []string) error {
	// 更新所有的alias指向
	// 1、删除在新的namespace中不存在的alias
	waitDel := make([]string, 0)
	for _, ns := range oldNS {
		if !contains(newNS, ns) {
			waitDel = append(waitDel, ns)
		}
	}
	for _, ns := range waitDel {
		path := fmt.Sprintf(PatternAlias, ns, aType, address)
		err := req.Storage.Delete(ctx, path)
		if err != nil {
			return err
		}
	}

	// 2、添加在旧的namespace中不存在的alias
	waitAdd := make([]string, 0)
	for _, ns := range newNS {
		if !contains(oldNS, ns) {
			waitAdd = append(waitAdd, ns)
		}
	}

	for _, ns := range waitAdd {
		path := fmt.Sprintf(PatternAlias, ns, aType, address)
		alias := modules.Alias{
			Source: address,
		}
		entry, err := logical.StorageEntryJSON(path, &alias)
		if err != nil {
			return err
		}
		err = req.Storage.Put(ctx, entry)
		if err != nil {
			return err
		}
	}
	return nil
}

func contains(nsSlice []string, ns string) bool {
	for _, n := range nsSlice {
		if n == ns {
			return true
		}
	}
	return false
}
