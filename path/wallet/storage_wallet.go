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
	patternWallet = "%s/wallet/%s"
)

type walletStorage struct {
}

func newWalletStorage() *walletStorage {
	return &walletStorage{}
}

func (as *walletStorage) readWallet(ctx context.Context, req *logical.Request, namespace, address string) (*modules.Wallet, error) {
	path := fmt.Sprintf(patternWallet, nameSpaceGlobal, address)
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

	if namespace == nameSpaceGlobal {
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
	wallet, err := as.readWallet(ctx, req, nameSpaceGlobal, overwrite.Address)
	if err != nil {
		return nil, err
	}

	wallet.UpdateTime = overwrite.UpdateTime
	wallet.NameSpaces = overwrite.NameSpaces

	path := fmt.Sprintf(patternWallet, nameSpaceGlobal, overwrite.Address)
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
	}

	insertPath := fmt.Sprintf(patternWallet, nameSpaceGlobal, walletEty.Address)
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

func (as *walletStorage) listWallet(ctx context.Context, req *logical.Request, namespace string) (map[string]*modules.Wallet, []string, error) {
	entries, err := req.Storage.List(ctx, fmt.Sprintf(patternWallet, nameSpaceGlobal, ""))
	if err != nil {
		return nil, nil, err
	}

	out := make([]string, 0)
	wallets := make(map[string]*modules.Wallet)

	for _, entry := range entries {
		wallet, subErr := as.readWallet(ctx, req, namespace, entry)
		if subErr != nil {
			continue
		}
		wallets[entry] = wallet
		out = append(out, entry)

	}

	return wallets, out, nil
}

func (as *walletStorage) deleteWallet(ctx context.Context, req *logical.Request, address string) error {
	return req.Storage.Delete(ctx, fmt.Sprintf(patternWallet, nameSpaceGlobal, address))
}
