package base

import (
	"crypto/ecdsa"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/roson9527/vault_eth_wallet/modules"
	"github.com/roson9527/vault_eth_wallet/path/doc"
	"github.com/roson9527/vault_eth_wallet/utils"
	"time"
)

var CryptoETH = &cryptoETH{}

type cryptoETH struct {
}

func (c *cryptoETH) GenerateKey() (*modules.WalletExtra, error) {
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

	address := crypto.PubkeyToAddress(privateKey.PublicKey).Hex()

	return &modules.WalletExtra{
		AddressAlias: map[string]string{
			doc.ChainETH: address,
		},
		PrivateKey: privateKeyString,
		NameSpaces: make([]string, 0),
		Wallet: modules.Wallet{
			Address:    address,
			PublicKey:  publicKeyString,
			UpdateTime: time.Now().Unix(),
		},
		CryptoType: doc.CryptoSECP256K1,
	}, nil
}

func (c *cryptoETH) PrivateToWallet(pri string) (*modules.WalletExtra, error) {
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

	address := crypto.PubkeyToAddress(privateKey.PublicKey).Hex()

	return &modules.WalletExtra{
		AddressAlias: map[string]string{
			doc.ChainETH: address,
		},
		PrivateKey: privateKeyString,
		Wallet: modules.Wallet{
			Address:    address,
			PublicKey:  publicKeyString,
			UpdateTime: time.Now().Unix(),
		},
		CryptoType: doc.CryptoSECP256K1,
	}, nil
}
