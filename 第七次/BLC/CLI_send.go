package BLC

import (
	"strconv"
	"os"
	"fmt"
	"encoding/hex"
)

func (cli *SJB_CLI) SJB_send(from []string,to []string,amount []string,nodeId string,minenow bool) {

	blockchain := SJB_BlockchainObject(nodeId)
	defer blockchain.SJB_DB.Close()
	utxoSet := &SJB_UTXOSet{blockchain}

	value, _ := strconv.Atoi(amount[0])

	if  value <= 0{
		fmt.Println("amount is wrong" )
		os.Exit(1)
	}

	if (minenow) {
		blockchain.SJB_MineNewBlock(from, to, amount,nodeId)
		utxoSet.SJB_Update()
	}else{
		nodeAddress = fmt.Sprintf("localhost:%s",nodeId)
		var txs []*SJB_Transaction
		for id := range mempool {
			tx := mempool[id]
			txs = append(txs, &tx)
		}
		if nodeAddress != knowNodes[0] {
			tx := SJB_NewSimpleTransaction(from[0], to[0], int64(value), utxoSet, txs, nodeId)
			mempool[hex.EncodeToString(tx.SJB_TxHash)] = *tx
			SJB_sendTx(knowNodes[0], tx)
		}else{
			println("主节点挖矿命令需要 + -mine")
		}
	}
}