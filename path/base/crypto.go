package base

import (
	"crypto/ecdsa"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/roson9527/vault_eth_wallet/modules"
	"github.com/roson9527/vault_eth_wallet/utils"
	"time"
)

func GenerateKey() (*modules.Wallet, error) {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		return nil, err
	}
	defer utils.ZeroKey(privateKey)

	privateKeyBytes := crypto.FromECDSA(privateKey)
	privateKeyString := hexutil.Encode(privateKeyBytes)[2:] // 移除0x

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf(errCastingPubToECDSA)
	}

	publicKeyBytes := crypto.FromECDSAPub(publicKeyECDSA)
	publicKeyString := hexutil.Encode(publicKeyBytes)[4:]

	return &modules.Wallet{
		Address:    crypto.PubkeyToAddress(privateKey.PublicKey).Hex(),
		PrivateKey: privateKeyString,
		PublicKey:  publicKeyString,
		UpdateTime: time.Now().Unix(),
		NameSpaces: make([]string, 0),
	}, nil
}

func PrivateToWallet(pri string) (*modules.Wallet, error) {
	privateKey, err := crypto.HexToECDSA(pri)
	if err != nil {
		return nil, errors.New("private key error")
	}
	defer utils.ZeroKey(privateKey)

	privateKeyBytes := crypto.FromECDSA(privateKey)
	privateKeyString := hexutil.Encode(privateKeyBytes)[2:] // 移除0x

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf(errCastingPubToECDSA)
	}

	publicKeyBytes := crypto.FromECDSAPub(publicKeyECDSA)
	publicKeyString := hexutil.Encode(publicKeyBytes)[4:]

	return &modules.Wallet{
		Address:    crypto.PubkeyToAddress(privateKey.PublicKey).Hex(),
		PrivateKey: privateKeyString,
		PublicKey:  publicKeyString,
		UpdateTime: time.Now().Unix(),
	}, nil
}
