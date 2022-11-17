package wallet

import (
	"fmt"
	"github.com/hashicorp/vault/sdk/framework"
)

func Path() []*framework.Path {
	pMgr := NewPathMgr()
	return []*framework.Path{
		pMgr.listWalletPath(fmt.Sprintf(PatternWallet, framework.GenericNameRegex(fieldNameSpace), "?")),
		pMgr.createWalletPath(fmt.Sprintf(PatternWallet, NameSpaceGlobal, "new")),
		pMgr.walletExportPath(fmt.Sprintf(PatternWallet, framework.GenericNameRegex(fieldNameSpace), framework.GenericNameRegex(fieldAddress)) + "/export"),
		pMgr.readWalletPath(fmt.Sprintf(PatternWallet, framework.GenericNameRegex(fieldNameSpace), framework.GenericNameRegex(fieldAddress))),
		pMgr.walletSignTxPath(
			fmt.Sprintf(PatternWallet,
				framework.GenericNameRegex(fieldNameSpace),
				framework.GenericNameRegex(fieldAddress)+"/sign_tx")),
		pMgr.policyPath(fmt.Sprintf(patternPolicy, framework.GenericNameRegex(fieldNameSpace))),
	}
}

type PathMgr struct {
	pathWallet
	pathPolicy
}

func NewPathMgr() *PathMgr {
	storageIns := newStorage()
	return &PathMgr{
		pathWallet{storageIns},
		pathPolicy{storageIns},
	}
}
