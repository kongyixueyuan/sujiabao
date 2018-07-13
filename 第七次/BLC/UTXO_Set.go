package BLC

import (
	"github.com/boltdb/bolt"
	"log"
	"encoding/hex"
	"fmt"
	"bytes"
)



const utxoTableName  = "utxoTableName"

type SJB_UTXOSet struct {
	SJB_Blockchain *SJB_Blockchain
}

func (utxoSet *SJB_UTXOSet) SJB_ResetUTXOSet()  {

	err := utxoSet.SJB_Blockchain.SJB_DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(utxoTableName))
		if b != nil {
			err := tx.DeleteBucket([]byte(utxoTableName))
			if err!= nil {
				log.Panic(err)
			}

		}

		b ,_ = tx.CreateBucket([]byte(utxoTableName))
		if b != nil {
			txOutputsMap := utxoSet.SJB_Blockchain.SJB_FindUTXOMap()
			for keyHash,outs := range txOutputsMap {
				txHash,_ := hex.DecodeString(keyHash)
				b.Put(txHash,outs.SJB_Serialize())
			}
		}
		return nil
	})

	if err != nil {
		log.Panic(err)
	}

}

func (utxoSet *SJB_UTXOSet) SJB_findUTXOForAddress(address string) []*SJB_UTXO{

	var utxos []*SJB_UTXO
	utxoSet.SJB_Blockchain.SJB_DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(utxoTableName))
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			txOutputs := SJB_DeserializeTXOutputs(v)
			for _,utxo := range txOutputs.SJB_UTXOS  {
				if utxo.SJB_Output.SJB_UnLockScriptPubKeyWithAddress(address) {
					utxos = append(utxos,utxo)
				}
			}
		}
		return nil
	})

	return utxos
}




func (utxoSet *SJB_UTXOSet) SJB_GetBalance(address string) int64 {

	UTXOS := utxoSet.SJB_findUTXOForAddress(address)
	fmt.Println(len(UTXOS))
	var amount int64
	for _,utxo := range UTXOS  {
		amount += utxo.SJB_Output.SJB_Value
	}

	return amount
}



func (utxoSet *SJB_UTXOSet) SJB_FindUnPackageSpendableUTXOS(from string, txs []*SJB_Transaction) []*SJB_UTXO {

	var unUTXOs []*SJB_UTXO

	spentTXOutputs := make(map[string][]int)

	for _,tx := range txs {

		if tx.SJB_IsCoinbaseTransaction() == false {
			for _, in := range tx.SJB_Vins {
				publicKeyHash := SJB_Base58Decode([]byte(from))

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

			if out.SJB_UnLockScriptPubKeyWithAddress(from) {

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

	return unUTXOs

}

func (utxoSet *SJB_UTXOSet) SJB_FindSpendableUTXOS(from string,amount int64,txs []*SJB_Transaction) (int64,map[string][]int)  {

	unPackageUTXOS := utxoSet.SJB_FindUnPackageSpendableUTXOS(from,txs)

	spentableUTXO := make(map[string][]int)

	var money int64 = 0

	for _,UTXO := range unPackageUTXOS {
		money += UTXO.SJB_Output.SJB_Value;
		txHash := hex.EncodeToString(UTXO.SJB_TxHash)
		spentableUTXO[txHash] = append(spentableUTXO[txHash],UTXO.SJB_Index)
		if money >= amount{
			return  money,spentableUTXO
		}
	}


	utxoSet.SJB_Blockchain.SJB_DB.View(func(tx *bolt.Tx) error {

		b := tx.Bucket([]byte(utxoTableName))

		if b != nil {
			c := b.Cursor()
			UTXOBREAK:
			for k, v := c.First(); k != nil; k, v = c.Next() {

				txOutputs := SJB_DeserializeTXOutputs(v)

				for _,utxo := range txOutputs.SJB_UTXOS {

					money += utxo.SJB_Output.SJB_Value
					txHash := hex.EncodeToString(utxo.SJB_TxHash)
					spentableUTXO[txHash] = append(spentableUTXO[txHash],utxo.SJB_Index)

					if money >= amount {
						 break UTXOBREAK;
					}
				}
			}

		}

		return nil
	})

	if money < amount{
		log.Panic("余额不足......")
	}


	return  money,spentableUTXO
}


func (utxoSet *SJB_UTXOSet) SJB_Update()  {

	block := utxoSet.SJB_Blockchain.SJB_Iterator().SJB_Next()

	ins := []*SJB_TXInput{}
	outsMap := make(map[string]*SJB_TXOutputs)

	for _,tx := range block.SJB_Txs {
		for _,in := range tx.SJB_Vins {
			ins = append(ins,in)
		}
	}

	for _,tx := range block.SJB_Txs  {

		utxos := []*SJB_UTXO{}

		for index,out := range tx.SJB_Vouts  {
			isSpent := false
			for _,in := range ins  {
				if in.SJB_Vout == index && bytes.Compare(tx.SJB_TxHash ,in.SJB_TxHash) == 0 && bytes.Compare(out.SJB_Ripemd160Hash,SJB_Ripemd160Hash(in.SJB_PublicKey)) == 0 {
					isSpent = true
					continue
				}
			}

			if isSpent == false {
				utxo := &SJB_UTXO{tx.SJB_TxHash,index,out}
				utxos = append(utxos,utxo)
			}

		}

		if len(utxos) > 0 {
			txHash := hex.EncodeToString(tx.SJB_TxHash)
			outsMap[txHash] = &SJB_TXOutputs{utxos}
		}

	}

	err := utxoSet.SJB_Blockchain.SJB_DB.Update(func(tx *bolt.Tx) error{

		b := tx.Bucket([]byte(utxoTableName))

		if b != nil {

			for _,in := range ins {

				txOutputsBytes := b.Get(in.SJB_TxHash)
				if len(txOutputsBytes) == 0 {
					continue
				}

				fmt.Println(txOutputsBytes)

				txOutputs := SJB_DeserializeTXOutputs(txOutputsBytes)

				fmt.Println(txOutputs)

				UTXOS := []*SJB_UTXO{}
				isNeedDelete := false
				for _,utxo := range txOutputs.SJB_UTXOS  {
					if in.SJB_Vout == utxo.SJB_Index && bytes.Compare(utxo.SJB_Output.SJB_Ripemd160Hash,SJB_Ripemd160Hash(in.SJB_PublicKey)) == 0 {
						isNeedDelete = true
					} else {
						UTXOS = append(UTXOS,utxo)
					}
				}



				if isNeedDelete {
					b.Delete(in.SJB_TxHash)
					if len(UTXOS) > 0 {

						preTXOutputs := outsMap[hex.EncodeToString(in.SJB_TxHash)]

						preTXOutputs.SJB_UTXOS = append(preTXOutputs.SJB_UTXOS,UTXOS...)

						outsMap[hex.EncodeToString(in.SJB_TxHash)] = preTXOutputs

					}
				}

			}

			for keyHash,outPuts := range outsMap  {
				keyHashBytes,_ := hex.DecodeString(keyHash)
				b.Put(keyHashBytes,outPuts.SJB_Serialize())
			}

		}

		return nil
	})

	if err != nil {
		log.Panic(err)
	}

}




