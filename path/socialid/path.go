package socialid

import (
	"fmt"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/roson9527/vault_eth_wallet/path/doc"
	"github.com/roson9527/vault_eth_wallet/path/storage"
)

func Path() []*framework.Path {
	pMgr := newPathMgr()
	out := make([]*framework.Path, 0)
	for _, app := range []string{doc.AppDiscord, doc.AppTwitter} {
		out = append(out, pMgr.builder(app)...)
	}
	return out
}

type pathMgr struct {
	handler handler
}

func (pmgr *pathMgr) builder(app string) []*framework.Path {
	out := make([]*framework.Path, 0)
	out = append(out, pmgr.handler.push(fmt.Sprintf(storage.PatternSocialID, doc.NameSpaceGlobal, app,
		framework.GenericNameRegex(doc.FieldUser))+"/"+doc.PathSubNew))
	out = append(out, pmgr.handler.list(fmt.Sprintf(storage.PatternSocialID,
		framework.GenericNameRegex(doc.FieldNameSpace), app, "?")))
	out = append(out, pmgr.handler.action(fmt.Sprintf(storage.PatternSocialID,
		framework.GenericNameRegex(doc.FieldNameSpace), app, framework.GenericNameRegex(doc.FieldUser))))
	out = append(out, pmgr.handler.export(fmt.Sprintf(storage.PatternSocialID,
		framework.GenericNameRegex(doc.FieldNameSpace), app,
		framework.GenericNameRegex(doc.FieldUser)+doc.PathSubExport)))
	return out
}

func newPathMgr() *pathMgr {
	storageIns := storage.NewCore()
	return &pathMgr{
		handler: handler{
			callback: callback{
				Storage: storageIns,
			},
		},
	}
}
