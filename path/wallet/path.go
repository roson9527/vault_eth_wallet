package wallet

import (
	"fmt"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/roson9527/vault_eth_wallet/path/doc"
	"github.com/roson9527/vault_eth_wallet/path/storage"
)

func Path() []*framework.Path {
	pMgr := NewPathMgr()
	return []*framework.Path{
		pMgr.listWalletPath(fmt.Sprintf(storage.PatternWallet, framework.GenericNameRegex(doc.FieldNameSpace), "?")),
		pMgr.createWalletPath(fmt.Sprintf(storage.PatternWallet, doc.NameSpaceGlobal, "new")),
		pMgr.walletExportPath(fmt.Sprintf(storage.PatternWallet, framework.GenericNameRegex(doc.FieldNameSpace), framework.GenericNameRegex(doc.FieldAddress)) + "/export"),
		pMgr.walletActionPath(fmt.Sprintf(storage.PatternWallet, framework.GenericNameRegex(doc.FieldNameSpace), framework.GenericNameRegex(doc.FieldAddress))),
		pMgr.walletSignTxPath(
			fmt.Sprintf(storage.PatternWallet,
				framework.GenericNameRegex(doc.FieldNameSpace),
				framework.GenericNameRegex(doc.FieldAddress)+"/sign_tx")),
	}
}

type PathMgr struct {
	pathWallet
}

func NewPathMgr() *PathMgr {
	storageIns := storage.NewCore()
	return &PathMgr{
		pathWallet{storageIns},
	}
}
