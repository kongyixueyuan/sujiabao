package BLC

import (
	"time"
	"fmt"
	"bytes"
	"encoding/gob"
	"log"
)

type SJB_Block struct {

	SJB_Height int64
	SJB_PrevBlockHash []byte
	SJB_Txs []*SJB_Transaction
	SJB_Timestamp int64
	SJB_Hash []byte
	SJB_Nonce int64
}



func (block *SJB_Block) SJB_HashTransactions() []byte  {


	var txHashes [][]byte

	for _, tx := range block.SJB_Txs {
		txHashes = append(txHashes, tx.SJB_Serialize())
	}
	newMerkleTree  := SJB_NewMerkleTree(txHashes)

	return newMerkleTree.SJB_RootNode.SJB_Data
}


func (block *SJB_Block)SJB_Serialize() []byte{
	var result  bytes.Buffer
	encoder := gob.NewEncoder(&result)
	err := encoder.Encode(block)
	if err != nil {
		log.Panic(err)
	}

	return result.Bytes()
}

func SJB_DeSerianlize(blockBytes []byte) *SJB_Block{

	var block SJB_Block

	decoder := gob.NewDecoder(bytes.NewReader(blockBytes))
	err := decoder.Decode(&block)
	if err != nil {
		log.Panic(err)
	}

	return &block
}

func SJB_NewBlock(txs []*SJB_Transaction,height int64,prevBlockHash []byte) *SJB_Block {

	block := &SJB_Block{height,prevBlockHash,txs,time.Now().Unix(),nil,0}

	pow := SJB_NewProofOfWork(block)
	hash,nonce := pow.SJB_Run()

	block.SJB_Hash = hash[:]
	block.SJB_Nonce = nonce

	fmt.Println("******************* new block **********************")

	return block

}


func SJB_CreateGenesisBlock(txs []*SJB_Transaction) *SJB_Block {
	return SJB_NewBlock(txs,1, []byte{0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0})
}

