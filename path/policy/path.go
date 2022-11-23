package policy

import (
	"fmt"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/roson9527/vault_eth_wallet/path/doc"
	"github.com/roson9527/vault_eth_wallet/path/storage"
)

func Path() []*framework.Path {
	pMgr := newPathMgr()
	return []*framework.Path{
		pMgr.handler.policy(fmt.Sprintf(storage.PatternPolicy, framework.GenericNameRegex(doc.FieldNameSpace))),
	}
}

type pathMgr struct {
	handler handler
}

func newPathMgr() *pathMgr {
	storageIns := storage.NewCore()
	return &pathMgr{
		handler{
			callback{
				Storage: storageIns,
			},
		},
	}
}
