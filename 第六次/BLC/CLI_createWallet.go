package BLC

func (cli *SJB_CLI) SJB_createWallet()  {

	wallets,_ := SJB_NewWallets()

	wallets.SJB_CreateNewWallet()
}
