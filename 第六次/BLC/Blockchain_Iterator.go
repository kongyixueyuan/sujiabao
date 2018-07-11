package BLC

import (
	"github.com/boltdb/bolt"
	"log"
)

type SJB_BlockchainIterator struct {
	SJB_CurrentHash []byte
	SJB_DB  *bolt.DB
}



func (blockchainIterator *SJB_BlockchainIterator)SJB_Next() *SJB_Block{

	var block *SJB_Block
	err := blockchainIterator.SJB_DB.View(func(tx *bolt.Tx) error{
		b := tx.Bucket([]byte(blockTableName))
		if b!= nil{
			currentBlockBytes := b.Get(blockchainIterator.SJB_CurrentHash)
			block = SJB_DeSerianlize(currentBlockBytes)
			blockchainIterator.SJB_CurrentHash = block.SJB_PrevBlockHash
		}
		return nil
	})

	if err != nil {

		log.Panic(err)
	}

	return block
}