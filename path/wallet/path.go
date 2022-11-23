package wallet

import (
	"fmt"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/roson9527/vault_eth_wallet/path/doc"
	"github.com/roson9527/vault_eth_wallet/path/storage"
)

func Path() []*framework.Path {
	pMgr := newPathMgr()
	return []*framework.Path{
		pMgr.handler.list(fmt.Sprintf(storage.PatternWallet, framework.GenericNameRegex(doc.FieldNameSpace), "?")),
		pMgr.handler.create(fmt.Sprintf(storage.PatternWallet, doc.NameSpaceGlobal, doc.PathSubNew)),
		pMgr.handler.export(fmt.Sprintf(storage.PatternWallet, framework.GenericNameRegex(doc.FieldNameSpace), framework.GenericNameRegex(doc.FieldAddress)) + doc.PathSubExport),
		pMgr.handler.action(fmt.Sprintf(storage.PatternWallet, framework.GenericNameRegex(doc.FieldNameSpace), framework.GenericNameRegex(doc.FieldAddress))),
		pMgr.handler.signTx(
			fmt.Sprintf(storage.PatternWallet,
				framework.GenericNameRegex(doc.FieldNameSpace),
				framework.GenericNameRegex(doc.FieldAddress)+doc.PathSubSignTx)),
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
