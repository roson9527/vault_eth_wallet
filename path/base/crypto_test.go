package base

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/roson9527/vault_eth_wallet/modules"
	"golang.org/x/crypto/sha3"
	"math/big"
	"strconv"
	"testing"
)

func TestSignatureTx(t *testing.T) {
	testAccount, _ := GenerateKey()
	params := modules.SignParams{
		Nonce:     4,
		ToAddress: nil,
		Amount:    big.NewInt(1),
		GasLimit:  1,
		GasPrice:  big.NewInt(1),
		Data:      nil,
		ChainId:   big.NewInt(1),
	}
	ret, err := SignatureTx(testAccount, &params)
	if err != nil {
		t.Error(err)
	}
	t.Log(ret.TransactionHash, ret.Signed)
}

func TestSign(t *testing.T) {
	priKey := "18dfe3522881fb925801abd62f3c418f87a95813d44c84faf6caac9e9394f712"
	data := "hello"
	//byteData:=[]byte(data)//common.FromHex(data)
	//fmt.Println(byteData)
	privateKey, _ := crypto.HexToECDSA(priKey)

	h, _ := TextAndHash([]byte(data))
	signed, err := crypto.Sign(h, privateKey)
	fmt.Println(err)
	fmt.Println(common.Bytes2Hex(signed))
}

func TextAndHash(data []byte) ([]byte, string) {
	msg := fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(data), string(data))
	hasher := sha3.NewLegacyKeccak256()
	hasher.Write([]byte(msg))
	return hasher.Sum(nil), msg
}

func TestLen(t *testing.T) {
	str := "I really did make this message"
	hexData := []byte(str)
	fmt.Println(len(str), str)
	fmt.Println(len(hexData), hexData)
}

func TestKeccak256(t *testing.T) {
	str := "I really did make this message"
	fmt.Println(ethHashData([]byte(str)))
	fmt.Println(ethHashString(str))
}

func ethHashData(data []byte) []byte {
	validationMsg := "\x19Ethereum Signed Message:\n" + strconv.Itoa(len(data))
	return crypto.Keccak256([]byte(validationMsg), data)
}

func ethHashString(data string) []byte {
	validationMsg := "\x19Ethereum Signed Message:\n" + strconv.Itoa(len(data)) + data
	return crypto.Keccak256([]byte(validationMsg))
}
