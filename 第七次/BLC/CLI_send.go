package BLC

import (
	"strconv"
	"fmt"
)

func (cli *SJB_CLI) SJB_send(from []string,to []string,amount []string,nodeId string,minenow bool) {

	blockchain := SJB_BlockchainObject(nodeId)
	defer blockchain.SJB_DB.Close()

	if (minenow) {
		blockchain.SJB_MineNewBlock(from, to, amount,nodeId)
		utxoSet := &SJB_UTXOSet{blockchain}
		utxoSet.SJB_Update()
	}else{
		nodeAddress = fmt.Sprintf("localhost:%s",nodeId)
		if nodeAddress != knowNodes[0]{
			utxoSet := &SJB_UTXOSet{blockchain}
			var txs []*SJB_Transaction

			for index,addressfrom := range from {
				value, _ := strconv.Atoi(amount[index])
				tx := SJB_NewSimpleTransaction(addressfrom, to[index], int64(value), utxoSet,txs,nodeId)
				SJB_sendTx(knowNodes[0], tx)
				txs = append(txs, tx)
			}
			fmt.Printf("%d个交易发送",len(txs))
		}
	}
}