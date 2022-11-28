package client

import (
	"github.com/hashicorp/vault/api"
	"strings"
)

type Client struct {
	Wallet *wallet
	Policy *policy
	Social *social
}

type core struct {
	Meta *api.Client
	conf *Config
}

func NewClient(conf *Config) *Client {
	cli, err := api.NewClient(nil)
	if err != nil {
		panic(err)
	}
	err = cli.SetAddress(conf.Address)
	if err != nil {
		panic(err)
	}
	cli.SetToken(conf.Token)

	if !strings.HasSuffix(conf.SecretPath, "/") {
		conf.SecretPath += "/"
	}

	c := &core{
		Meta: cli,
		conf: conf,
	}

	return &Client{
		Wallet: newWallet(c),
		Policy: newPolicy(c),
		Social: &social{c},
	}
}
