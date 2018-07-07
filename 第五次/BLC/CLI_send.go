package BLC

import (
	"fmt"
	"os"
)

func (cli *CLI) send(from []string,to []string,amount []string)  {


	if DBExists() == false {
		fmt.Println("请先创建区块链")
		os.Exit(1)
	}

	blockchain := BlockchainObject()
	defer blockchain.DB.Close()

	blockchain.MineNewBlock(from,to,amount)

}