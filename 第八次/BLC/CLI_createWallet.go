package BLC

import "fmt"

func (cli *SJB_CLI) SJB_createWallet(nodeId string)  {

	wallets,_ := SJB_NewWallets(nodeId)

	wallets.SJB_CreateNewWallet(nodeId)

	fmt.Printf("已有%d个钱包",len(wallets.SJB_WalletsMap))
}
