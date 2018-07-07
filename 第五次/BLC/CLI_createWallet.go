package BLC

func (cli *CLI) createWallet()  {

	wallets,_ := NewWallets()

	wallets.CreateNewWallet()
}
