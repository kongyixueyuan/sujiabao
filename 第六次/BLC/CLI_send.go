package BLC

import (
	"fmt"
	"os"
)

func (cli *SJB_CLI) SJB_send(from []string,to []string,amount []string)  {


	if SJB_DBExists() == false {
		fmt.Println("block chain is not exist")
		os.Exit(1)
	}

	blockchain := SJB_BlockchainObject()
	defer blockchain.SJB_DB.Close()

	blockchain.SJB_MineNewBlock(from,to,amount)

	utxoSet := &SJB_UTXOSet{blockchain}

	//转账成功以后，需要更新一下
	utxoSet.SJB_Update()

}