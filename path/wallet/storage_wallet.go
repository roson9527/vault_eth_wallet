package wallet

import (
	"context"
	"errors"
	"fmt"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/roson9527/vault_eth_wallet/modules"
	"github.com/roson9527/vault_eth_wallet/path/base"
)

const (
	PatternWallet = "%s/wallet/%s"
)

type walletStorage struct {
	aliasStorage
}

func newWalletStorage() *walletStorage {
	return &walletStorage{}
}

func (as *walletStorage) readWallet(ctx context.Context, req *logical.Request, namespace, address string) (*modules.Wallet, error) {
	path := fmt.Sprintf(PatternWallet, NameSpaceGlobal, address)
	entry, err := req.Storage.Get(ctx, path)

	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, errors.New("wallet not found")
	}

	var wallet modules.Wallet
	err = entry.DecodeJSON(&wallet)

	if wallet.Address == "" {
		return nil, fmt.Errorf("failed to deserialize wallet at %s", path)
	}

	if namespace == NameSpaceGlobal {
		return &wallet, nil
	}

	if wallet.NameSpaces == nil {
		return nil, fmt.Errorf("wallet %s does not have namespace %s", wallet.Address, namespace)
	}

	for _, ns := range wallet.NameSpaces {
		if ns == namespace {
			return &wallet, nil
		}
	}

	return nil, fmt.Errorf("not support namespace %s", namespace)
}

func (as *walletStorage) updateWallet(ctx context.Context, req *logical.Request, overwrite *modules.Wallet) (*modules.Wallet, error) {
	wallet, err := as.readWallet(ctx, req, NameSpaceGlobal, overwrite.Address)
	if err != nil {
		return nil, err
	}

	wallet.UpdateTime = overwrite.UpdateTime
	wallet.NameSpaces = overwrite.NameSpaces
	if overwrite.Network != "" {
		wallet.Network = overwrite.Network
	}

	path := fmt.Sprintf(PatternWallet, NameSpaceGlobal, overwrite.Address)
	entry, err := logical.StorageEntryJSON(path, wallet)
	if err != nil {
		return nil, err
	}

	err = req.Storage.Put(ctx, entry)
	if err != nil {
		return nil, err
	}

	return overwrite, nil
}

func (as *walletStorage) createWallet(ctx context.Context, req *logical.Request, overwrite *modules.Wallet) (*modules.Wallet, error) {
	var walletEty *modules.Wallet
	var err error

	if overwrite.PrivateKey != "" {
		walletEty = overwrite
	} else {
		walletEty, err = base.GenerateKey()
		if err != nil {
			return nil, err
		}

		if overwrite.NameSpaces != nil && len(overwrite.NameSpaces) > 0 {
			walletEty.NameSpaces = overwrite.NameSpaces
		}
		walletEty.Network = networkETH
	}

	insertPath := fmt.Sprintf(PatternWallet, NameSpaceGlobal, walletEty.Address)

	entry, err := logical.StorageEntryJSON(insertPath, walletEty)
	if err != nil {
		return nil, err
	}

	err = req.Storage.Put(ctx, entry)
	if err != nil {
		return nil, err
	}

	return walletEty, nil
}

func (as *walletStorage) listWallet(ctx context.Context, req *logical.Request, namespace string) ([]string, error) {
	if namespace == NameSpaceGlobal {
		return req.Storage.List(ctx, fmt.Sprintf(PatternWallet, NameSpaceGlobal, ""))
	}

	return req.Storage.List(ctx, fmt.Sprintf(PatternAlias, namespace, ""))
}

func (as *walletStorage) deleteWallet(ctx context.Context, req *logical.Request, address string) error {
	return req.Storage.Delete(ctx, fmt.Sprintf(PatternWallet, NameSpaceGlobal, address))
}
