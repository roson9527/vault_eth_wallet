package storage

import (
	"context"
	"errors"
	"fmt"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/roson9527/vault_eth_wallet/modules"
	"github.com/roson9527/vault_eth_wallet/path/doc"
)

const (
	// PatternWalletByChain pattern for wallet
	/*

		PatternWalletByChain = [global|namespace]/wallet/[cryptoType]/[chain]/[address]
		PatternWallet = [global|namespace]/wallet/[cryptoType]/[def_address]
		PatternWalletAdmin = "global/wallet/[cryptoType]"
	*/
	PatternWalletByChain = "%s/wallet/%s/%s/%s"
	PatternWallet        = "%s/wallet/%s/%s"
	PatternWalletAdmin   = "%s/wallet/%s"
)

type walletStorage struct {
}

func newWalletStorage() *walletStorage {
	return &walletStorage{}
}

func (as *walletStorage) ReadFromPath(ctx context.Context, req *logical.Request, path string) (*modules.WalletExtra, error) {
	if path == "" {
		return nil, errors.New("path is nil")
	}

	entry, err := req.Storage.Get(ctx, path)

	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, ErrWalletNotFound
	}

	var wallet modules.WalletExtra
	err = entry.DecodeJSON(&wallet)
	if err != nil {
		return nil, err
	}
	return &wallet, nil
}

func (as *walletStorage) Read(ctx context.Context, req *logical.Request, namespace, cryptoType, address string) (*modules.WalletExtra, error) {
	path := fmt.Sprintf(PatternWallet, doc.NameSpaceGlobal, cryptoType, address)
	wallet, err := as.ReadFromPath(ctx, req, path)

	if err != nil {
		return nil, err
	}

	if namespace == doc.NameSpaceGlobal {
		return wallet, nil
	}

	if doc.CryptoSECP256K1 != wallet.CryptoType {
		return nil, fmt.Errorf("failed to deserialize wallet at %s", path)
	}

	if wallet.NameSpaces == nil {
		return nil, fmt.Errorf("wallet does not have namespace %s", namespace)
	}

	for _, ns := range wallet.NameSpaces {
		if ns == namespace {
			return wallet, nil
		}
	}

	return nil, fmt.Errorf("not support namespace %s", namespace)
}

func (as *walletStorage) Put(ctx context.Context, req *logical.Request, wallet *modules.WalletExtra) error {
	if wallet == nil {
		return errors.New("wallet is nil")
	}

	if wallet.Address == "" {
		return errors.New("address is nil")
	}

	path := fmt.Sprintf(PatternWallet, doc.NameSpaceGlobal, wallet.CryptoType, wallet.Address)
	hclog.Default().Info("wallet:put", "path", path)
	entry, err := logical.StorageEntryJSON(path, wallet)
	if err != nil {
		return err
	}
	return req.Storage.Put(ctx, entry)
}

func (as *walletStorage) Update(ctx context.Context, req *logical.Request, payload *modules.WalletExtra) (*modules.WalletExtra, error) {
	wallet, err := as.Read(ctx, req, doc.NameSpaceGlobal, payload.CryptoType, payload.Address)
	if err != nil {
		return nil, err
	}

	wallet.UpdateTime = payload.UpdateTime
	wallet.NameSpaces = payload.NameSpaces
	wallet.AddressAlias = payload.AddressAlias
	wallet.Extra = payload.Extra

	path := fmt.Sprintf(PatternWallet, doc.NameSpaceGlobal, wallet.CryptoType, payload.Address)
	entry, err := logical.StorageEntryJSON(path, wallet)
	if err != nil {
		return nil, err
	}

	err = req.Storage.Put(ctx, entry)
	if err != nil {
		return nil, err
	}

	return wallet, nil
}

func (as *walletStorage) Create(ctx context.Context, req *logical.Request, payload *modules.WalletExtra) (*modules.WalletExtra, error) {
	if payload.CryptoType == "" {
		return nil, ErrCryptoTypeIsNil
	}
	insertPath := fmt.Sprintf(PatternWallet, doc.NameSpaceGlobal, payload.CryptoType, payload.Address)
	entry, err := logical.StorageEntryJSON(insertPath, payload)
	if err != nil {
		return nil, err
	}

	err = req.Storage.Put(ctx, entry)
	if err != nil {
		return nil, err
	}

	return payload, nil
}

func (as *walletStorage) List(ctx context.Context, req *logical.Request, cryptoType string) ([]string, error) {
	return req.Storage.List(ctx, fmt.Sprintf(PatternWallet, doc.NameSpaceGlobal, cryptoType, ""))
}

func (as *walletStorage) Delete(ctx context.Context, req *logical.Request, cryptoType, address string) error {
	return req.Storage.Delete(ctx, fmt.Sprintf(PatternWallet, doc.NameSpaceGlobal, cryptoType, address))
}
