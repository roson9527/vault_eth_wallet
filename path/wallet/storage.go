package wallet

type storage struct {
	walletStorage *walletStorage
	policyStorage *policyStorage
}

func newStorage() *storage {
	return &storage{
		walletStorage: newWalletStorage(),
		policyStorage: newPolicyStorage(),
	}
}
