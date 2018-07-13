package BLC

import "fmt"

// 打印所有的钱包地址
func (cli *SJB_CLI) SJB_addressLists(nodeId string)  {

	wallets,_ := SJB_NewWallets(nodeId)

	for address,_ := range wallets.SJB_WalletsMap {

		fmt.Println(address)
	}
}