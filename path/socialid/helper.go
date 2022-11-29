package socialid

import (
	"fmt"
	"github.com/mitchellh/mapstructure"
	"github.com/roson9527/vault_eth_wallet/modules"
	"github.com/roson9527/vault_eth_wallet/path/base"
	"github.com/roson9527/vault_eth_wallet/path/doc"
)

func socialIDResponseData(socialId *modules.SocialID, extra bool) map[string]any {
	socialId.GAuth = nil
	if socialId.TwoFASecret != "" {
		auth, err := base.GAuth.Generate(socialId.TwoFASecret)
		if err == nil {
			socialId.GAuth = auth
		}
	}
	if !extra {
		out := map[string]any{
			doc.FieldUser:       socialId.User,
			doc.FieldUpdateTime: socialId.UpdateTime,
			doc.FieldApp:        socialId.App,
		}
		if socialId.GAuth != nil {
			out[doc.FieldGAuth] = socialId.GAuth
		}
		return out
	}
	out := make(map[string]any)
	_ = mapstructure.Decode(socialId, &out)
	return out
}

func aliasType(app string) string {
	return fmt.Sprintf(doc.AliasSocial, app)
}
