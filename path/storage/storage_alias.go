package storage

import (
	"context"
	"fmt"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/roson9527/vault_eth_wallet/modules"
	"github.com/roson9527/vault_eth_wallet/path/doc"
)

const (
	// PatternAlias is the pattern for the alias storage path
	/* PatternAlias = "%s/alias/%s/%s"
	 * [namespace]/alias/[type]/[address]
	 * e.g. type = social/discord global/alias/discord/0x1234567890
	 * e.g. type = wallet/eth global/alias/wallet/eth/0x1234567890
	 */
	PatternAlias = "%s/alias/%s/%s"
	nilPath      = ""
)

var (
	zeroStringSlice   = make([]string, 0)
	createStringSlice = []string{doc.NameSpaceGlobal}
)

type aliasStorage struct {
}

func newAliasStorage() *aliasStorage {
	return &aliasStorage{}
}

func (as *aliasStorage) ReadSrcPath(ctx context.Context, req *logical.Request, ns, aType, address string) (string, error) {
	path := fmt.Sprintf(PatternAlias, ns, aType, address)
	hclog.Default().Info("alias:read", "path", path)
	ety, err := req.Storage.Get(ctx, path)
	if err != nil || ety == nil {
		return "", err
	}
	var alias modules.Alias
	err = ety.DecodeJSON(&alias)
	if err != nil {
		return "", err
	}

	hclog.Default().Info("alias:read:alias", "alias.src", alias.Source)

	return alias.Source, nil
}

func (as *aliasStorage) List(ctx context.Context, req *logical.Request, namespace, aType string) ([]string, error) {
	// 获取namespace下所有对应vType的alias
	hclog.Default().Info("alias:list", "path", fmt.Sprintf(PatternAlias, namespace, aType, ""))
	return req.Storage.List(ctx, fmt.Sprintf(PatternAlias, namespace, aType, ""))
}

func (as *aliasStorage) update(ctx context.Context, req *logical.Request, aType, address, srcPath string, oldNS, newNS []string) error {
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
		hclog.Default().Info("alias:delete", "path", path)
		err := req.Storage.Delete(ctx, path)
		if err != nil {
			return err
		}
	}

	// WARN 因为多了alias, 所以在单独新增alias时，这里无法进行正常的刷新, 需要外部对参数进行判断
	// 2、添加在旧的namespace中不存在的alias
	waitAdd := make([]string, 0)
	for _, ns := range newNS {
		if !contains(oldNS, ns) && ns != "" {
			waitAdd = append(waitAdd, ns)
		}
	}

	for _, ns := range waitAdd {
		path := fmt.Sprintf(PatternAlias, ns, aType, address)
		// 如果是全局的alias, 则需要进行判断是否已经添加过
		if ns == doc.NameSpaceGlobal {
			ety, err := req.Storage.Get(ctx, path)
			if err != nil {
				return err
			}
			if ety != nil {
				continue
			}
		}

		alias := modules.Alias{
			Source: srcPath,
		}
		hclog.Default().Info("alias:add", "path", path, "src", srcPath)
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

func (as *aliasStorage) Update(ctx context.Context, req *logical.Request, aType, address, srcPath string, oldNS, newNS []string) error {
	newNS = as.fixNs(newNS)
	return as.update(ctx, req, aType, address, srcPath, oldNS, newNS)
}

func (as *aliasStorage) Delete(ctx context.Context, req *logical.Request, aType, address string, oldNS []string) error {
	oldNS = as.fixNs(oldNS)
	return as.update(ctx, req, aType, address, nilPath, oldNS, zeroStringSlice)
}

func (as *aliasStorage) Create(ctx context.Context, req *logical.Request, aType, address, srcPath string, newNS []string) error {
	newNS = as.fixNs(newNS)
	return as.update(ctx, req, aType, address, srcPath, zeroStringSlice, newNS)
}

func (as *aliasStorage) fixNs(ns []string) []string {
	flag := false
	for _, n := range ns {
		if n == doc.NameSpaceGlobal {
			flag = true
			break
		}
	}
	if !flag {
		ns = append(ns, doc.NameSpaceGlobal)
	}
	return ns
}

func contains(nsSlice []string, ns string) bool {
	for _, n := range nsSlice {
		if n == ns {
			return true
		}
	}
	return false
}
