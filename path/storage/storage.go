package storage

type Core struct {
	Wallet *walletStorage
	Policy *policyStorage
	Alias  *aliasStorage
	Social *socialIDStorage
}

var core = newCore()

func newCore() *Core {
	return &Core{
		Alias:  newAliasStorage(),
		Wallet: newWalletStorage(),
		Policy: newPolicyStorage(),
		Social: newSocialIDStorage(),
	}
}

func NewCore() *Core {
	return core
}
