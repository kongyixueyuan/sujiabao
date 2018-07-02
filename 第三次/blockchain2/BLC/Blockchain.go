package BLC

import (
	"github.com/boltdb/bolt"
	"log"
	"fmt"
	"math/big"
	"time"
	"os"
)


const dbName = "myblockchain.db"
const blockTableName = "blocktable"

type Blockchain struct {
	Tip []byte
	DB  *bolt.DB
}

// iterator
func (blockchain *Blockchain) Iterator() *BlockchainIterator {

	return &BlockchainIterator{blockchain.Tip, blockchain.DB}
}

func DBExists() bool {
	if _, err := os.Stat(dbName); os.IsNotExist(err) {
		return false
	}
	return true
}


func (blc *Blockchain) PrintAllchain() {

	blockchainIterator := blc.Iterator()

	for {
		block := blockchainIterator.Next()
		fmt.Printf("PrevBlockHash：%x\n", block.PrevBlockHash)
		fmt.Printf("Height：%d\n", block.Height)
		fmt.Printf("Data：%s\n", block.Data)
		fmt.Printf("Timestamp：%s\n", time.Unix(block.Timestamp, 0).Format("2006-01-02 03:04:05 PM"))
		fmt.Printf("Hash：%x\n", block.Hash)
		fmt.Printf("Nonce：%d\n", block.Nonce)

		fmt.Println()

		var hashInt big.Int
		hashInt.SetBytes(block.PrevBlockHash)
		if big.NewInt(0).Cmp(&hashInt) == 0 {
			break;
		}
	}

}


func (blc *Blockchain) AddBlockToBlockchain(data string) {

	err := blc.DB.Update(func(tx *bolt.Tx) error {

		b := tx.Bucket([]byte(blockTableName))
		if b != nil {

			blockBytes := b.Get(blc.Tip)
			block := DeSerianlize(blockBytes)

			newBlock := NewBlock(data, block.Height+1, block.Hash)
			err := b.Put(newBlock.Hash, newBlock.Serialize())
			if err != nil {
				log.Panic(err)
			}
			err = b.Put([]byte("l"), newBlock.Hash)
			if err != nil {
				log.Panic(err)
			}
			blc.Tip = newBlock.Hash
		}

		return nil
	})

	if err != nil {
		log.Panic(err)
	}
}

func CreateBlockchainWithGenesisBlock(data string) {

	// 判断数据库是否存在
	if DBExists() {
		fmt.Println("数据库已存在.......")
		os.Exit(1)
	}


	fmt.Println("正在创建创世区块.......")

	// 创建或者打开数据库
	db, err := bolt.Open(dbName, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Update(func(tx *bolt.Tx) error {

		// 创建数据库表
		b, err := tx.CreateBucket([]byte(blockTableName))

		if err != nil {
			log.Panic(err)
		}

		if b != nil {
			// 创建创世区块
			genesisBlock := CreateGenesisBlock(data)
			// 将创世区块存储到表中
			err := b.Put(genesisBlock.Hash, genesisBlock.Serialize())
			if err != nil {
				log.Panic(err)
			}
			err = b.Put([]byte("l"), genesisBlock.Hash)
			if err != nil {
				log.Panic(err)
			}
		}

		return nil
	})

}


func BlockchainObject() *Blockchain {

	db, err := bolt.Open(dbName, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}

	var tip []byte

	err = db.View(func(tx *bolt.Tx) error {

		b := tx.Bucket([]byte(blockTableName))

		if b != nil {
			tip = b.Get([]byte("l"))

		}


		return nil
	})

	return &Blockchain{tip,db}
}
