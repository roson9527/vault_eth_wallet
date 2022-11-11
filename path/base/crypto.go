package base

import (
	"bytes"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
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

func SignatureTx(account *modules.Wallet, params *modules.SignParams) (*modules.SignResult, error) {
	privateKey, err := crypto.HexToECDSA(account.PrivateKey)
	if err != nil {
		return nil, err
	}
	// 防止没有进入回收被检索到
	defer utils.ZeroKey(privateKey)

	tx := newTx(params)
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(params.ChainId), privateKey)
	if err != nil {
		return nil, err
	}

	var signedTxBuff bytes.Buffer
	err = signedTx.EncodeRLP(&signedTxBuff)
	if err != nil {
		return nil, err
	}

	return &modules.SignResult{
		Signed:          hexutil.Encode(signedTxBuff.Bytes()),
		TransactionHash: signedTx.Hash().Hex(),
	}, nil
}

func Signature(account *modules.Wallet, params *modules.SignParams) (*modules.SignResult, error) {
	privateKey, err := crypto.HexToECDSA(account.PrivateKey)
	if err != nil {
		return nil, err
	}
	defer utils.ZeroKey(privateKey)
	dataBytes := params.Data
	if !params.IsHashData {
		hash := crypto.Keccak256Hash(dataBytes)
		dataBytes = hash.Bytes()
	}

	sign, err := crypto.Sign(dataBytes, privateKey)
	if err != nil {
		return nil, err
	}

	return &modules.SignResult{
		Signed:          hexutil.Encode(sign),
		TransactionHash: hexutil.Encode(dataBytes),
	}, nil
}

func newTx(params *modules.SignParams) *types.Transaction {
	if params.ToAddress == nil {
		return types.NewContractCreation(params.Nonce, params.Amount, params.GasLimit, params.GasPrice, params.Data)
	}
	return types.NewTransaction(
		params.Nonce, *params.ToAddress, params.Amount, params.GasLimit, params.GasPrice, params.Data)
}

func Verify(acct *modules.Wallet, dataByte []byte, signature string, isHash bool) (bool, error) {
	privateKey, err := crypto.HexToECDSA(acct.PrivateKey)
	if err != nil {
		return false, err
	}
	defer utils.ZeroKey(privateKey)

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return false, fmt.Errorf(errCastingPubToECDSA)
	}
	publicKeyBytes := crypto.FromECDSAPub(publicKeyECDSA)

	signatureBytes, err := hexutil.Decode(signature)
	if err != nil {
		return false, err
	}

	if !isHash {
		hash := crypto.Keccak256Hash(dataByte)
		dataByte = hash.Bytes()
	}

	signPubKey, err := crypto.Ecrecover(dataByte, signatureBytes)
	if err != nil {
		return false, err
	}

	matches := bytes.Equal(signPubKey, publicKeyBytes)
	if !matches {
		return false, fmt.Errorf(errSignCheckFailed)
	}

	return true, nil
}
