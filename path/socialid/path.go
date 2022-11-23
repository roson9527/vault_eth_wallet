package socialid

import (
	"fmt"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/roson9527/vault_eth_wallet/path/doc"
	"github.com/roson9527/vault_eth_wallet/path/storage"
)

func Path() []*framework.Path {
	pMgr := NewPathMgr()
	out := make([]*framework.Path, 0)
	for _, app := range []string{doc.AppDiscord, doc.AppTwitter} {
		out = append(out, pMgr.Builder(app)...)
	}
	return out
}

type PathMgr struct {
	pathSocialID
}

func (pmgr *PathMgr) Builder(app string) []*framework.Path {
	out := make([]*framework.Path, 0)
	out = append(out, pmgr.pushPath(fmt.Sprintf(storage.PatternSocialID, doc.NameSpaceGlobal, app, doc.PathSubNew)))
	out = append(out, pmgr.listPath(fmt.Sprintf(storage.PatternSocialID, framework.GenericNameRegex(doc.FieldNameSpace),
		app, "?")))
	out = append(out, pmgr.actionPath(fmt.Sprintf(storage.PatternSocialID, framework.GenericNameRegex(doc.FieldNameSpace),
		app, framework.GenericNameRegex(doc.FieldUser))))
	out = append(out, pmgr.exportPath(fmt.Sprintf(storage.PatternSocialID, framework.GenericNameRegex(doc.FieldNameSpace),
		app, framework.GenericNameRegex(doc.FieldUser)+doc.PathSubExport)))
	return out
}

func NewPathMgr() *PathMgr {
	storageIns := storage.NewCore()
	return &PathMgr{
		pathSocialID{storageIns},
	}
}
