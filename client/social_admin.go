package client

import (
	"fmt"
	"github.com/mitchellh/mapstructure"
	"github.com/roson9527/vault_eth_wallet/modules"
	"github.com/roson9527/vault_eth_wallet/path/doc"
	"github.com/roson9527/vault_eth_wallet/path/storage"
)

func (c *social) Delete(app, user string) error {
	_, err := c.Meta.Logical().Delete(c.conf.SecretPath + fmt.Sprintf(storage.PatternSocialID, doc.NameSpaceGlobal, app, user))
	return err
}

func (c *social) Put(user string, src *modules.SocialID) (*modules.SocialID, error) {
	var payload = make(map[string]any)
	data := make(map[string]any)
	err := mapstructure.WeakDecode(src, &data)
	if err != nil {
		return nil, err
	}

	payload[doc.FieldSocialID] = data

	sec, err := c.Meta.Logical().Write(c.conf.SecretPath+
		fmt.Sprintf(storage.PatternSocialID, doc.NameSpaceGlobal, src.App, user), payload)
	if err != nil {
		panic(err)
	}

	var out modules.SocialID
	err = mapstructure.WeakDecode(sec.Data, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *social) Export(namespace, app, user string) (*modules.SocialID, error) {
	sec, err := c.Meta.Logical().Read(c.conf.SecretPath + fmt.Sprintf(storage.PatternSocialID, namespace, app, user) +
		doc.PathSubExport)
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
