package socialid

import (
	"fmt"
	"github.com/hashicorp/vault/api"
	"github.com/roson9527/vault_eth_wallet/path/doc"
	"github.com/roson9527/vault_eth_wallet/path/storage"
	"testing"
)

func createLocalClient() *api.Client {
	cli, err := api.NewClient(nil)
	if err != nil {
		panic(err)
	}
	cli.SetAddress("http://127.0.0.1:8200")
	cli.SetToken("root")
	return cli
}

func TestSocialIdPush(t *testing.T) {
	cli := createLocalClient()
	sec, err := cli.Logical().Write("mock_web3/"+fmt.Sprintf(storage.PatternSocialID, doc.NameSpaceGlobal, doc.AppDiscord, "test1")+"/new", map[string]any{
		"namespaces": []string{"test_b"},
	})
	fmt.Println(sec, err)
}
