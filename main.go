package main

import (
	"os"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/sdk/plugin"
)

func main() {
	var err error
	defer func() {
		if err != nil {
			logger := hclog.New(&hclog.LoggerOptions{})
			logger.Error("plugin shutting down", "error", err)
		}
		os.Exit(1)
	}()

	apiClientMeta := &api.PluginAPIClientMeta{}
	flags := apiClientMeta.FlagSet()
	err = flags.Parse(os.Args[1:])
	if err != nil {
		return
	}

	tlsConfig := apiClientMeta.GetTLSConfig()
	tlsProviderFunc := api.VaultPluginTLSProvider(tlsConfig)

	err = plugin.Serve(&plugin.ServeOpts{
		BackendFactoryFunc: Factory,
		TLSProviderFunc:    tlsProviderFunc,
	})
}
