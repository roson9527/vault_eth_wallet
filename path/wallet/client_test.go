package wallet

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/hashicorp/vault/api"
	"testing"
)

const mockPluginName = "mock_web3"

func createLocalClient() *api.Client {
	cli, err := api.NewClient(nil)
	if err != nil {
		panic(err)
	}
	cli.SetAddress("http://127.0.0.1:8200")
	cli.SetToken("root")
	return cli
}

func buildPluginIns(cli *api.Client) {
	//cli := createLocalClient()
	_, err := cli.Sys().MountConfig("mock_web3")
	if err == nil {
		return
	}
	err = cli.Sys().Mount(mockPluginName, &api.MountInput{
		Type:        "web3_wallet",
		Description: "web3 wallet",
	})
	if err != nil {
		panic(err)
	}
}

func TestClient_BuildPluginIns(t *testing.T) {
	cli := createLocalClient()
	buildPluginIns(cli)
}

func TestClient_CreateAndRead(t *testing.T) {
	cli := createLocalClient()
	sec, err := cli.Logical().Write("mock_web3/"+fmt.Sprintf(patternWallet, nameSpaceGlobal, "new"), map[string]any{
		"namespaces": []string{"test_b"},
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(sec, err)
	address := sec.Data["address"].(string)
	sec, err = cli.Logical().Read("mock_web3/" + fmt.Sprintf(patternWallet, nameSpaceGlobal, address))
	if err != nil {
		panic(err)
	}
	fmt.Println(sec, err)
	sec, err = cli.Logical().Read("mock_web3/" + fmt.Sprintf(patternWallet, nameSpaceGlobal, address) + "/export")
	if err != nil {
		panic(err)
	}
	fmt.Println(sec, err)
}

func TestClient_List(t *testing.T) {
	cli := createLocalClient()
	sec, err := cli.Logical().List("mock_web3/" + fmt.Sprintf(patternWallet, nameSpaceGlobal, ""))
	if err != nil {
		panic(err)
	}
	fmt.Println(sec, err)
}

func TestClient_SignTx(t *testing.T) {
	binary := "0x02f9038201830202398085061b8dd5ba83037e1294a69babef1ca67a37ffaf7a485dfff3382056e78c825c02b901441cff79cd0000000000000000000000000f6a2783548b4846d4b8334d43623a7d9d035810000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000c4588116d8000000000000000000000000000000000000000000000000092d60f88871110c000000000000000000000000000000000000000000002a1a0f328008300000000000000000000000000000000000000000000000039ab1a46f630ab59067dae10000000000000000000000000000000000000000000000000000b4d369988ee900000000000000000000000000000000000000000000000000000000636b597cb2000000000000000000000000000000000000000000000000000000000104bf00000000000000000000000000000000000000000000000000000000f901cdf87a94c02aaa39b223fe8d0a0e5c4f27ead9083c756cc2f863a09b732dcd2fed2b720cc896f075566c04bfed55bdaacc1458ffbc4b3e1d9bd084a030bd84b96629f958113934633d3bd1b64c3d259a85c57ceac65da8c5ec9bf3a7a00facc5ed2299a6cf172c64296b0494ffceff3b01c4184cedbdc1ea2fb5f36ec6f8dd94e1d92f1de49caec73514f696fea2a7d5441498e5f8c6a00000000000000000000000000000000000000000000000000000000000000000a00000000000000000000000000000000000000000000000000000000000000004a00000000000000000000000000000000000000000000000000000000000000002a09b637a02e6f8cc8aa1e3935c0b27bde663b11428c7707039634076a3fb8a0c48a00000000000000000000000000000000000000000000000000000000000000141a00000000000000000000000000000000000000000000000000000000000000142f85994bbbbca6a901c926f240b89eacb641d8aec7aeafdf842a007f90970db204e2e62b5c0e6d0d237d824119af9f10f7f3873f41a53d57893d8a0f7c84b5d1f3a0563cd20b346c00dcaaaf749870e80d340ce8dc6213863709a60d6940f6a2783548b4846d4b8334d43623a7d9d035810c080a074abeb983f56edde0729427739cd1b7b805675b6db8528cf7f58c970b35c0087a063ceba8e7aff7394e0eecde78f6021a37f1613d3d7090e9a41953acbf882b0cd"
	_ = binary
	cli := createLocalClient()
	sec, err := cli.Logical().List("mock_web3/" + fmt.Sprintf(patternWallet, nameSpaceGlobal, ""))
	if err != nil {
		panic(err)
	}
	address := sec.Data["keys"].([]any)[0].(string)
	fmt.Println(address)

	data := make(map[string]any)
	data["tx_binary"] = binary
	sec, err = cli.Logical().Write("mock_web3/"+fmt.Sprintf(patternWallet, "test_b", address)+"/sign_tx", data)
	if err != nil {
		panic(err)
	}
	fmt.Println(sec, err)
	var signedTx types.Transaction
	_ = signedTx.UnmarshalBinary(hexutil.MustDecode(sec.Data["tx_binary"].(string)))
	fmt.Println(getSender(&signedTx))
}

func TestTransactionTx(t *testing.T) {
	ethCli, _ := ethclient.Dial("https://rpc.ankr.com/eth")
	tx, _, _ := ethCli.TransactionByHash(context.Background(), common.HexToHash("0xa195c8b6332cf24bf2262e21e4b0919348c5748901ecda81e168b7bf6a440b38"))
	binary, _ := tx.MarshalBinary()
	fmt.Println(hexutil.Encode(binary))
	jbin, _ := tx.MarshalJSON()
	fmt.Println(string(jbin))
}

func getSender(tx *types.Transaction) (common.Address, error) {
	signer := types.LatestSignerForChainID(tx.ChainId())
	return signer.Sender(tx)
}
