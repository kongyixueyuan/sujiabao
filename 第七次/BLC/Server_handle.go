package BLC

import (
	"bytes"
	"log"
	"encoding/gob"
	"fmt"
	"encoding/hex"
	"github.com/boltdb/bolt"
)

func SJB_handleVersion(request []byte,bc *SJB_Blockchain)  {

	var buff bytes.Buffer
	var payload Version

	dataBytes := request[COMMANDLENGTH:]
	buff.Write(dataBytes)
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}

	bestHeight := bc.GetBestHeight()
	remoteBestHeight := payload.SJB_BestHeight

	if bestHeight > remoteBestHeight {
		SJB_sendVersion(payload.SJB_AddrFrom,bc)
	} else if bestHeight < remoteBestHeight {
		SJB_sendGetBlocks(payload.SJB_AddrFrom)
	}

	if !SJB_nodeIsKnow(payload.SJB_AddrFrom) {
		knowNodes = append(knowNodes, payload.SJB_AddrFrom)
	}
}


func SJB_handleAddr(request []byte,bc *SJB_Blockchain)  {
	var buff bytes.Buffer
	var payload addr

	buff.Write(request[COMMANDLENGTH:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}

	knowNodes = append(knowNodes, payload.AddrList...)
	fmt.Printf("There are %d known nodes now!\n", len(knowNodes))
	for _, node := range knowNodes {
		SJB_sendGetBlocks(node)
	}
}

func SJB_handleGetblocks(request []byte,bc *SJB_Blockchain)  {
	var buff bytes.Buffer
	var payload GetBlocks

	dataBytes := request[COMMANDLENGTH:]

	buff.Write(dataBytes)
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}
	//send block hash
	blocks := bc.GetBlockHashes()
	SJB_sendInv(payload.SJB_AddrFrom, BLOCK_TYPE, blocks)
}

func SJB_handleGetData(request []byte,bc *SJB_Blockchain)  {

	var buff bytes.Buffer
	var payload GetData

	dataBytes := request[COMMANDLENGTH:]
	buff.Write(dataBytes)
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}

	if payload.SJB_Type == BLOCK_TYPE {
		block, err := bc.SJB_GetBlock([]byte(payload.Hash))
		if err != nil {
			return
		}
		SJB_sendBlock(payload.SJB_AddrFrom, block)
	}

	if payload.SJB_Type == "tx" {
		txID := hex.EncodeToString(payload.Hash)
		TxData := mempool[txID]
		SJB_sendTx(payload.SJB_AddrFrom, &TxData)
	}
}

func SJB_handleBlock(request []byte,bc *SJB_Blockchain)  {

	var buff bytes.Buffer
	var payload BlockData

	dataBytes := request[COMMANDLENGTH:]

	buff.Write(dataBytes)
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}

	blockBytes := payload.Block
	block := SJB_DeSerianlize(blockBytes)

	fmt.Println("Recevied a new block!")
	bc.SJB_AddBlock(block)
	fmt.Printf("Added block %x\n", block.SJB_Hash)

	if len(transactionArray) > 0 {
		blockHash := transactionArray[0]
		SJB_sendGetData(payload.SJB_AddrFrom, "block", blockHash)
		transactionArray = transactionArray[1:]
	} else {
		fmt.Println("数据库重置......")
		UTXOSet := &SJB_UTXOSet{bc}
		UTXOSet.SJB_ResetUTXOSet()
	}
}

func SJB_handleTx(request []byte,bc *SJB_Blockchain)  {

	var buff bytes.Buffer
	var payload TxData

	dataBytes := request[COMMANDLENGTH:]

	buff.Write(dataBytes)
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}

	txData := payload.SJB_TransactonData
	tx := SJB_DeserializeTransaction(txData)
	mempool[hex.EncodeToString(tx.SJB_TxHash)] = tx

	if nodeAddress == knowNodes[0] {
		for _, node := range knowNodes {
			if node != nodeAddress && node != payload.SJB_AddrFrom {
				SJB_sendInv(node, "tx", [][]byte{tx.SJB_TxHash})
			}
		}
	} else {
		if len(mempool) >= 2 && len(MyminerAddress) > 0 {
		MineTransactions:
			var txs []*SJB_Transaction
			for id := range mempool {
				tx := mempool[id]
				if bc.SJB_VerifyTransaction(&tx,txs) {
					txs = append(txs, &tx)
				}
			}

			if len(txs) == 0 {
				fmt.Println("All transactions are invalid! Waiting for new ones...")
				return
			}

			cbTx := SJB_NewCoinbaseTransaction(MyminerAddress)
			txs = append(txs, cbTx)

			var block *SJB_Block
			bc.SJB_DB.View(func(tx *bolt.Tx) error {
				b := tx.Bucket([]byte(blockTableName))
				if b != nil {
					hash := b.Get([]byte("l"))
					blockBytes := b.Get(hash)
					block = SJB_DeSerianlize(blockBytes)
				}
				return nil
			})

			newBlock := SJB_NewBlock(txs,block.SJB_Height,block.SJB_PrevBlockHash)
			UTXOSet := SJB_UTXOSet{bc}
			UTXOSet.SJB_ResetUTXOSet()

			fmt.Println("New block is mined!")

			for _, tx := range txs {
				txID := hex.EncodeToString(tx.SJB_TxHash)
				delete(mempool, txID)
			}

			for _, node := range knowNodes {
				if node != nodeAddress {
					SJB_sendInv(node, "block", [][]byte{newBlock.SJB_Hash})
				}
			}

			if len(mempool) > 0 {
				goto MineTransactions
			}
		}
	}
}


func SJB_handleInv(request []byte,bc *SJB_Blockchain)  {

	var buff bytes.Buffer
	var payload Inv

	dataBytes := request[COMMANDLENGTH:]

	buff.Write(dataBytes)
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}

	if payload.SJB_Type == BLOCK_TYPE {
		blockHash := payload.SJB_Items[0]
		SJB_sendGetData(payload.SJB_AddrFrom, BLOCK_TYPE , blockHash)

		if len(payload.SJB_Items) >= 1 {
			transactionArray = payload.SJB_Items[1:]
		}
	}

	if payload.SJB_Type == TX_TYPE {
		txHash := payload.SJB_Items[0]

		if mempool[hex.EncodeToString(txHash)].SJB_TxHash == nil {
			SJB_sendGetData(payload.SJB_AddrFrom, TX_TYPE, txHash)
		}
	}

}