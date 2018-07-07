package BLC

import (
	"time"
	"fmt"
	"bytes"
	"encoding/gob"
	"log"
	"crypto/sha256"
)

type Block struct {

	Height int64
	PrevBlockHash []byte
	Txs []*Transaction
	Timestamp int64
	Hash []byte
	Nonce int64
}



func (block *Block) HashTransactions() []byte  {


	var txHashes [][]byte
	var txHash [32]byte

	for _, tx := range block.Txs {
		txHashes = append(txHashes, tx.TxHash)
	}
	txHash = sha256.Sum256(bytes.Join(txHashes, []byte{}))

	return txHash[:]

}


func (block *Block)Serialize() []byte{
	var result  bytes.Buffer
	encoder := gob.NewEncoder(&result)
	err := encoder.Encode(block)
	if err != nil {
		log.Panic(err);
	}

	return result.Bytes();
}

func DeSerianlize(blockBytes []byte) *Block{

	var block Block

	decoder := gob.NewDecoder(bytes.NewReader(blockBytes))
	err := decoder.Decode(&block)
	if err != nil {
		log.Panic(err)
	}

	return &block
}

func NewBlock(txs []*Transaction,height int64,prevBlockHash []byte) *Block {

	block := &Block{height,prevBlockHash,txs,time.Now().Unix(),nil,0}

	pow := NewProofOfWork(block)
	hash,nonce := pow.Run()

	block.Hash = hash[:]
	block.Nonce = nonce

	fmt.Println("******************* new block **********************")

	return block

}


func CreateGenesisBlock(txs []*Transaction) *Block {

	return NewBlock(txs,1, []byte{0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0})
}

