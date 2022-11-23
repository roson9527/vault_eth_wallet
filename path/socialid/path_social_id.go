package socialid

import (
	"fmt"
	"github.com/mitchellh/mapstructure"
	"github.com/roson9527/vault_eth_wallet/modules"
	"github.com/roson9527/vault_eth_wallet/path/doc"
	"github.com/roson9527/vault_eth_wallet/path/storage"
)

type pathSocialID struct {
	Storage *storage.Core
}

func socialIDResponseData(socialId *modules.SocialID, extra bool) map[string]any {
	if !extra {
		return map[string]any{
			doc.FieldUser:       socialId.User,
			doc.FieldUpdateTime: socialId.UpdateTime,
			doc.FieldApp:        socialId.App,
		}
	}
	out := make(map[string]any)
	_ = mapstructure.Decode(socialId, &out)
	return out
}

func app2AType(app string) string {
	return fmt.Sprintf(doc.AliasSocial, app)
}
