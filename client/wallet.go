package client

type wallet struct {
	ETH  *walletETH
	Text *walletText
}

func newWallet(c *core) *wallet {
	return &wallet{
		ETH:  newWalletETH(c),
		Text: newWalletText(c),
	}
}
