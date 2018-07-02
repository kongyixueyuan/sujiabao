package BLC

import (
	"bytes"
	"log"
	"encoding/gob"
	"crypto/sha256"
	"encoding/hex"
)

// UTXO
type Transaction struct {

	TxHash []byte
	Vins []*TXInput
	Vouts []*TXOutput
}


func (tx *Transaction) IsCoinbaseTransaction() bool {
	if(len(tx.Vins[0].TxHash) == 0 && tx.Vins[0].Vout == -1){
		return true
	}

	return false
}


func NewCoinbaseTransaction(address string) *Transaction {

	txInput := &TXInput{[]byte{}, -1,"Genesis Data"}
	txOutput := &TXOutput{10, address}
	txCoinbase := &Transaction{[]byte{}, []*TXInput{txInput}, []*TXOutput{txOutput}}
	txCoinbase.HashTransaction()
	return txCoinbase
}

func (tx *Transaction) HashTransaction()  {

	var result bytes.Buffer

	encoder := gob.NewEncoder(&result)

	err := encoder.Encode(tx)
	if err != nil {
		log.Panic(err)
	}

	hash := sha256.Sum256(result.Bytes())

	tx.TxHash = hash[:]
}



func NewSimpleTransaction(from string, to string,amount int,blockchain *Blockchain,txs []*Transaction) *Transaction {

	money,spendableUTXODic := blockchain.FindSpendableUTXOS(from,amount,txs)

	var txIntputs []*TXInput
	var txOutputs []*TXOutput

	for txHash,indexArray := range spendableUTXODic  {

		txHashBytes,_  := hex.DecodeString(txHash)
		for _,index := range indexArray  {
			txInput := &TXInput{txHashBytes,index,from}
			txIntputs = append(txIntputs,txInput)
		}

	}

	txOutput := &TXOutput{int64(amount),to}
	txOutputs = append(txOutputs,txOutput)

	txOutput = &TXOutput{int64(money) - int64(amount),from}
	txOutputs = append(txOutputs,txOutput)

	tx := &Transaction{[]byte{},txIntputs,txOutputs}

	tx.HashTransaction()

	return tx

}




