package client

import "github.com/hashicorp/vault/api"

type Client struct {
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
	return &Client{
		Meta: cli,
		conf: conf,
	}
}
