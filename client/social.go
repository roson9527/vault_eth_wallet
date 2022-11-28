package client

import (
	"fmt"
	"github.com/mitchellh/mapstructure"
	"github.com/roson9527/vault_eth_wallet/modules"
	"github.com/roson9527/vault_eth_wallet/path/doc"
	"github.com/roson9527/vault_eth_wallet/path/storage"
)

type social struct {
	*core
}

func (c *social) List(project, app string) ([]string, error) {
	sec, err := c.Meta.Logical().List(c.conf.SecretPath + fmt.Sprintf(storage.PatternSocialID, project, app, ""))
	if err != nil {
		return nil, err
	}

	// 没有数据
	if sec == nil {
		return []string{}, nil
	}

	out := make([]string, 0)
	for _, v := range sec.Data[doc.FieldKeys].([]any) {
		out = append(out, fmt.Sprintf("%s", v))
	}

	return out, nil
}

func (c *social) Read(project, app, user string) (*modules.SocialID, error) {
	sec, err := c.Meta.Logical().Read(c.conf.SecretPath + fmt.Sprintf(storage.PatternSocialID, project, app, user))
	if err != nil {
		return nil, err
	}

	// 没有数据
	if sec == nil {
		return nil, nil
	}
	var out modules.SocialID
	err = mapstructure.WeakDecode(sec.Data, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}
