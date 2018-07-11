package BLC

import (
	"fmt"
	"os"
)

func (cli *SJB_CLI) SJB_printchain()  {

	if SJB_DBExists() == false {
		fmt.Println("请先创建区块链")
		os.Exit(1)
	}

	blockchain := SJB_BlockchainObject()
	defer blockchain.SJB_DB.Close()

	blockchain.SJB_Printchain()

}