package BLC

import (
	"github.com/boltdb/bolt"
	"log"
	"fmt"
	"math/big"
	"time"
	"os"
	"strconv"
	"encoding/hex"
	"crypto/ecdsa"
	"bytes"
)

const dbName = "myblockchain_%s.db"
const blockTableName = "blockstable"

type SJB_Blockchain struct {
	SJB_Tip []byte
	SJB_DB  *bolt.DB
}

func (blockchain *SJB_Blockchain) SJB_Iterator() *SJB_BlockchainIterator {

	return &SJB_BlockchainIterator{blockchain.SJB_Tip, blockchain.SJB_DB}
}

func SJB_DBExists(dbName string) bool {
	if _, err := os.Stat(dbName); os.IsNotExist(err) {
		return false
	}

	return true
}

func (blc *SJB_Blockchain) SJB_Printchain() {

	blockchainIterator := blc.SJB_Iterator()

	for {
		block := blockchainIterator.SJB_Next()

		fmt.Printf("Height：%d\n", block.SJB_Height)
		fmt.Printf("PrevBlockHash：%x\n", block.SJB_PrevBlockHash)
		fmt.Printf("Timestamp： %s\n", time.Unix(block.SJB_Timestamp, 0).Format("2006-01-02 03:04:05 PM"))
		fmt.Printf("Hash：%x\n", block.SJB_Hash)
		fmt.Printf("Nonce：%d\n", block.SJB_Nonce)
		fmt.Println("Txs:")
		for _, tx := range block.SJB_Txs {

			fmt.Printf("%x\n", tx.SJB_TxHash)
			fmt.Println("Vins:")
			for _, in := range tx.SJB_Vins {
				fmt.Printf("%x\n", in.SJB_TxHash)
				fmt.Printf("%d\n", in.SJB_Vout)
				fmt.Printf("%s\n", in.SJB_PublicKey)
			}

			fmt.Println("Vouts:")
			for _, out := range tx.SJB_Vouts {
				fmt.Println(out.SJB_Value)
				fmt.Println(out.SJB_Ripemd160Hash)
			}
		}


		var hashInt big.Int
		hashInt.SetBytes(block.SJB_PrevBlockHash)
		if big.NewInt(0).Cmp(&hashInt) == 0 {
			break;
		}
	}

}

func (blc *SJB_Blockchain) SJB_AddBlockToBlockchain(txs []*SJB_Transaction) {

	err := blc.SJB_DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blockTableName))
		if b != nil {
			blockBytes := b.Get(blc.SJB_Tip)
			block := SJB_DeSerianlize(blockBytes)

			newBlock := SJB_NewBlock(txs, block.SJB_Height+1, block.SJB_Hash)
			err := b.Put(newBlock.SJB_Hash, newBlock.SJB_Serialize())
			if err != nil {
				log.Panic(err)
			}
			err = b.Put([]byte("l"), newBlock.SJB_Hash)
			if err != nil {
				log.Panic(err)
			}

			blc.SJB_Tip = newBlock.SJB_Hash
		}

		return nil
	})

	if err != nil {
		log.Panic(err)
	}
}

func SJB_CreateBlockchainWithGenesisBlock(address string,nodeId string) *SJB_Blockchain {

	dbname := fmt.Sprintf(dbName,nodeId)

	if SJB_DBExists(dbname) {
		fmt.Println("数据已经存在")
		os.Exit(1)
	}

	fmt.Println("创建创世区块。。。。。。")

	db, err := bolt.Open(dbname, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}

	var genesisHash []byte
	err = db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucket([]byte(blockTableName))

		if err != nil {
			log.Panic(err)
		}

		if b != nil {
			txCoinbase := SJB_NewCoinbaseTransaction(address)

			genesisBlock := SJB_CreateGenesisBlock([]*SJB_Transaction{txCoinbase})
			err := b.Put(genesisBlock.SJB_Hash, genesisBlock.SJB_Serialize())
			if err != nil {
				log.Panic(err)
			}
			err = b.Put([]byte("l"), genesisBlock.SJB_Hash)
			if err != nil {
				log.Panic(err)
			}

			genesisHash = genesisBlock.SJB_Hash
		}

		return nil
	})

	return &SJB_Blockchain{genesisHash, db}

}
func SJB_BlockchainObject(nodeId string) *SJB_Blockchain {

	dbname := fmt.Sprintf(dbName,nodeId)
	if SJB_DBExists(dbname) == false {
		fmt.Println("数据库不存在....")
		os.Exit(1)
	}

	db, err := bolt.Open(dbname, 0600, nil)
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

	return &SJB_Blockchain{tip, db}
}



func (blockchain *SJB_Blockchain) SJB_MineNewBlock(from []string, to []string, amount []string,nodeId string) {

	utxoSet := &SJB_UTXOSet{blockchain}

	//生成交易
	var txs []*SJB_Transaction

	for index,address := range from {
		value, _ := strconv.Atoi(amount[index])
		tx := SJB_NewSimpleTransaction(address, to[index], int64(value), utxoSet,txs,nodeId)
		txs = append(txs, tx)
	}
	//奖励
	tx := SJB_NewCoinbaseTransaction(from[0])
	txs = append(txs,tx)


	//验证交易
	_txs := []*SJB_Transaction{}

	for _,tx := range txs  {

		if blockchain.SJB_VerifyTransaction(tx,_txs) != true {
			log.Panic("ERROR: Invalid transaction")
		}

		_txs = append(_txs,tx)
	}

	//------------------------------------------------------------------------------------//

	//生成区块，保存数据库上链
	var block *SJB_Block
	blockchain.SJB_DB.View(func(tx *bolt.Tx) error {

		b := tx.Bucket([]byte(blockTableName))
		if b != nil {

			hash := b.Get([]byte("l"))

			blockBytes := b.Get(hash)

			block = SJB_DeSerianlize(blockBytes)

		}

		return nil
	})

	block = SJB_NewBlock(txs, block.SJB_Height+1, block.SJB_Hash)

	blockchain.SJB_DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blockTableName))
		if b != nil {

			b.Put(block.SJB_Hash, block.SJB_Serialize())

			b.Put([]byte("l"), block.SJB_Hash)

			blockchain.SJB_Tip = block.SJB_Hash

		}
		return nil
	})

}

func (blockchain *SJB_Blockchain) SJB_GetBalance(address string) int64 {

	utxos := blockchain.SJB_UnUTXOs(address,[]*SJB_Transaction{})

	var amount int64

	for _, utxo := range utxos {
		amount = amount + utxo.SJB_Output.SJB_Value
	}

	return amount
}

func (bclockchain *SJB_Blockchain) SJB_SignTransaction(tx *SJB_Transaction,privKey ecdsa.PrivateKey,txs []*SJB_Transaction)  {

	if tx.SJB_IsCoinbaseTransaction() {
		return
	}

	prevTXs := make(map[string]SJB_Transaction)

	for _, vin := range tx.SJB_Vins {
		prevTX, err := bclockchain.SJB_FindTransaction(vin.SJB_TxHash,txs)
		if err != nil {
			log.Panic(err)
		}
		prevTXs[hex.EncodeToString(prevTX.SJB_TxHash)] = prevTX
	}

	tx.SJB_Sign(privKey, prevTXs)
}


func (blockchain *SJB_Blockchain) SJB_VerifyTransaction(tx *SJB_Transaction,txs []*SJB_Transaction) bool {


	prevTXs := make(map[string]SJB_Transaction)

	for _, vin := range tx.SJB_Vins {
		prevTX, err := blockchain.SJB_FindTransaction(vin.SJB_TxHash,txs)
		if err != nil {
			log.Panic(err)
		}
		prevTXs[hex.EncodeToString(prevTX.SJB_TxHash)] = prevTX
	}

	return tx.SJB_Verify(prevTXs)
}


func (blockchain *SJB_Blockchain) SJB_FindTransaction(ID []byte,txs []*SJB_Transaction) (SJB_Transaction,error) {

	for _,tx := range txs  {
		if bytes.Compare(tx.SJB_TxHash, ID) == 0 {
			return *tx, nil
		}
	}

	bci := blockchain.SJB_Iterator()

	for {
		block := bci.SJB_Next()
		for _, tx := range block.SJB_Txs {
			if bytes.Compare(tx.SJB_TxHash, ID) == 0 {
				return *tx, nil
			}
		}
		var hashInt big.Int
		hashInt.SetBytes(block.SJB_PrevBlockHash)

		if big.NewInt(0).Cmp(&hashInt) == 0 {
			break;
		}
	}

	return SJB_Transaction{},nil
}

func (blockchain *SJB_Blockchain) SJB_UnUTXOs(address string,txs []*SJB_Transaction) []*SJB_UTXO {

	var unUTXOs []*SJB_UTXO
	spentTXOutputs := make(map[string][]int)

	for _,tx := range txs {
		//记录txs中已使用UTXO
		if tx.SJB_IsCoinbaseTransaction() == false {
			for _, in := range tx.SJB_Vins {
				publicKeyHash := SJB_Base58Decode([]byte(address))
				ripemd160Hash := publicKeyHash[1:len(publicKeyHash) - 4]

				if in.SJB_UnLockRipemd160Hash(ripemd160Hash) {
					key := hex.EncodeToString(in.SJB_TxHash)
					spentTXOutputs[key] = append(spentTXOutputs[key], in.SJB_Vout)
				}
			}
		}
	}


	for _,tx := range txs {
		Work1:
		for index,out := range tx.SJB_Vouts {
			if out.SJB_UnLockScriptPubKeyWithAddress(address) {
				fmt.Printf("address:%s  spend %d\n",address,spentTXOutputs)
				if len(spentTXOutputs) == 0 {
					utxo := &SJB_UTXO{tx.SJB_TxHash, index, out}
					unUTXOs = append(unUTXOs, utxo)
				} else {
					for hash,indexArray := range spentTXOutputs {
						txHashStr := hex.EncodeToString(tx.SJB_TxHash)

						if hash == txHashStr {
							var isUnSpentUTXO bool
							for _,outIndex := range indexArray {
								if index == outIndex {
									isUnSpentUTXO = true
									continue Work1
								}
								if isUnSpentUTXO == false {
									utxo := &SJB_UTXO{tx.SJB_TxHash, index, out}
									unUTXOs = append(unUTXOs, utxo)
								}
							}
						} else {
							utxo := &SJB_UTXO{tx.SJB_TxHash, index, out}
							unUTXOs = append(unUTXOs, utxo)
						}
					}
				}
			}
		}

	}

	blockIterator := blockchain.SJB_Iterator()

	for {
		block := blockIterator.SJB_Next()
		for i := len(block.SJB_Txs) - 1; i >= 0 ; i-- {
			tx := block.SJB_Txs[i]
			if tx.SJB_IsCoinbaseTransaction() == false {
				for _, in := range tx.SJB_Vins {
					publicKeyHash := SJB_Base58Decode([]byte(address))
					ripemd160Hash := publicKeyHash[1:len(publicKeyHash) - 4]
					if in.SJB_UnLockRipemd160Hash(ripemd160Hash) {
						key := hex.EncodeToString(in.SJB_TxHash)
						spentTXOutputs[key] = append(spentTXOutputs[key], in.SJB_Vout)
					}
				}
			}

			work:
			for index, out := range tx.SJB_Vouts {
				if out.SJB_UnLockScriptPubKeyWithAddress(address) {
					if spentTXOutputs != nil {
						if len(spentTXOutputs) != 0 {
							var isSpentUTXO bool
							for txHash, indexArray := range spentTXOutputs {
								for _, i := range indexArray {
									if index == i && txHash == hex.EncodeToString(tx.SJB_TxHash) {
										isSpentUTXO = true
										continue work
									}
								}
							}

							if isSpentUTXO == false {
								utxo := &SJB_UTXO{tx.SJB_TxHash, index, out}
								unUTXOs = append(unUTXOs, utxo)
							}
						} else {
							utxo := &SJB_UTXO{tx.SJB_TxHash, index, out}
							unUTXOs = append(unUTXOs, utxo)
						}
					}
				}
			}
		}

		//创世区块
		var hashInt big.Int
		hashInt.SetBytes(block.SJB_PrevBlockHash)
		if hashInt.Cmp(big.NewInt(0)) == 0 {
			break;
		}

	}

	return unUTXOs
}

func (blockchain *SJB_Blockchain) SJB_FindSpendableUTXOS(from string, amount int,txs []*SJB_Transaction) (int64, map[string][]int) {

	utxos := blockchain.SJB_UnUTXOs(from,txs)
	spendableUTXO := make(map[string][]int)

	var value int64
	for _, utxo := range utxos {
		value = value + utxo.SJB_Output.SJB_Value
		hash := hex.EncodeToString(utxo.SJB_TxHash)
		spendableUTXO[hash] = append(spendableUTXO[hash], utxo.SJB_Index)
		if value >= int64(amount) {
			break
		}
	}

	if value < int64(amount) {
		fmt.Printf("%s's account is not enough\n", from)
		os.Exit(1)
	}

	return value, spendableUTXO
}

func (blc *SJB_Blockchain) SJB_FindUTXOMap() map[string]*SJB_TXOutputs  {

	blcIterator := blc.SJB_Iterator()

	spentableUTXOsMap := make(map[string][]*SJB_TXInput)
	utxoMaps := make(map[string]*SJB_TXOutputs)

	for {
		//遍历所有区块
		block := blcIterator.SJB_Next()

		for i := len(block.SJB_Txs) - 1; i >= 0 ;i-- {
			//遍历一个区块中的所有交易

			txOutputs := &SJB_TXOutputs{[]*SJB_UTXO{}}
			tx := block.SJB_Txs[i]

			//如果不是coinbase交易，则把input保存到已使用UTXO字典
			if tx.SJB_IsCoinbaseTransaction() == false {
				for _,txInput := range tx.SJB_Vins {
					txHash := hex.EncodeToString(txInput.SJB_TxHash)
					spentableUTXOsMap[txHash] = append(spentableUTXOsMap[txHash],txInput)
				}
			}

			txHash := hex.EncodeToString(tx.SJB_TxHash)

			UTXOLoop:
			for index,out := range tx.SJB_Vouts  {
				txInputs := spentableUTXOsMap[txHash]
				if len(txInputs) > 0 {
					isSpent := false
					for _,in := range  txInputs {
						outPublicKey := out.SJB_Ripemd160Hash
						inPublicKey := in.SJB_PublicKey

						if bytes.Compare(outPublicKey,SJB_Ripemd160Hash(inPublicKey)) == 0{
							if index == in.SJB_Vout {
								isSpent = true
								continue UTXOLoop
							}
						}
					}

					if isSpent == false {
						utxo := &SJB_UTXO{tx.SJB_TxHash,index,out}
						txOutputs.SJB_UTXOS = append(txOutputs.SJB_UTXOS,utxo)
					}

				} else {
					//找不到已使用记录
					utxo := &SJB_UTXO{tx.SJB_TxHash,index,out}
					txOutputs.SJB_UTXOS = append(txOutputs.SJB_UTXOS,utxo)
				}
			}
			utxoMaps[txHash] = txOutputs
		}

		var hashInt big.Int
		hashInt.SetBytes(block.SJB_PrevBlockHash)

		if hashInt.Cmp(big.NewInt(0)) == 0 {
			break;
		}
	}

	return utxoMaps
}

func (bc *SJB_Blockchain) GetBestHeight() int64 {

	blockIterator := bc.SJB_Iterator()
	block := blockIterator.SJB_Next()
	bestheight := block.SJB_Height

	return bestheight
}

func (bc *SJB_Blockchain) GetBlockHashes() [][]byte {

	blockIterator := bc.SJB_Iterator()

	var blockHashs [][]byte

	for {
		block := blockIterator.SJB_Next()

		blockHashs = append(blockHashs,block.SJB_Hash)

		var hashInt big.Int
		hashInt.SetBytes(block.SJB_PrevBlockHash)

		if hashInt.Cmp(big.NewInt(0)) == 0 {
			break;
		}
	}

	return blockHashs
}

func (bc *SJB_Blockchain) SJB_GetBlock(blockHash []byte) ([]byte ,error) {

	var blockBytes []byte

	err := bc.SJB_DB.View(func(tx *bolt.Tx) error {

		b := tx.Bucket([]byte(blockTableName))

		if b != nil {

			blockBytes = b.Get(blockHash)

		}

		return nil
	})

	return blockBytes,err
}

func (bc *SJB_Blockchain) SJB_AddBlock(block *SJB_Block)  {

	err := bc.SJB_DB.Update(func(tx *bolt.Tx) error {

		b := tx.Bucket([]byte(blockTableName))

		if b != nil {

			blockExist := b.Get(block.SJB_Hash)

			if blockExist != nil {
				// 如果存在，不需要做任何过多的处理
				return nil
			}

			err := b.Put(block.SJB_Hash,block.SJB_Serialize())

			if err != nil {
				log.Panic(err)
			}

			// 最新的区块链的Hash
			blockHash := b.Get([]byte("l"))

			blockBytes := b.Get(blockHash)

			blockInDB := SJB_DeSerianlize(blockBytes)

			if blockInDB.SJB_Height < block.SJB_Height {

				b.Put([]byte("l"),block.SJB_Hash)
				bc.SJB_Tip = block.SJB_Hash
			}
		}

		return nil
	})

	if err != nil {
		log.Panic(err)
	}
}