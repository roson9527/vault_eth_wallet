package policy

import (
	"fmt"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/roson9527/vault_eth_wallet/path/doc"
	"github.com/roson9527/vault_eth_wallet/path/storage"
)

func Path() []*framework.Path {
	pMgr := NewPathMgr()
	return []*framework.Path{
		pMgr.policyPath(fmt.Sprintf(storage.PatternPolicy, framework.GenericNameRegex(doc.FieldNameSpace))),
	}
}

type PathMgr struct {
	pathPolicy
}

func NewPathMgr() *PathMgr {
	storageIns := storage.NewCore()
	return &PathMgr{
		pathPolicy{storageIns},
	}
}
