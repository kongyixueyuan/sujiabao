package BLC

import "fmt"

// 先用它去查询余额
func (cli *SJB_CLI) SJB_getBalance(address string)  {

	blockchain := SJB_BlockchainObject()

	defer blockchain.SJB_DB.Close()

	utxoSet := &SJB_UTXOSet{blockchain}
	amount := utxoSet.SJB_GetBalance(address)

	fmt.Printf("address % stotal %d Token\n",address,amount)


}
