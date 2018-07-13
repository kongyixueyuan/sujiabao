package BLC

import (
	"bytes"
	"log"
	"encoding/gob"
	"crypto/sha256"
	"encoding/hex"
	"crypto/ecdsa"
	"crypto/rand"
	"math/big"
	"crypto/elliptic"
	"time"
	"fmt"
	"os"
)

// UTXO
type SJB_Transaction struct {

	SJB_TxHash []byte
	SJB_Vins []*SJB_TXInput
	SJB_Vouts []*SJB_TXOutput
}


func (tx *SJB_Transaction) SJB_IsCoinbaseTransaction() bool {
	if(len(tx.SJB_Vins[0].SJB_TxHash) == 0 && tx.SJB_Vins[0].SJB_Vout == -1){
		return true
	}

	return false
}


func SJB_NewCoinbaseTransaction(address string) *SJB_Transaction {

	fmt.Print("reward for %s 10 token",address);
	txInput := &SJB_TXInput{[]byte{}, -1, nil,[]byte{}}
	txOutput := SJB_NewTXOutput(10,address)
	txCoinbase := &SJB_Transaction{[]byte{}, []*SJB_TXInput{txInput}, []*SJB_TXOutput{txOutput}}
	txCoinbase.SJB_HashTransaction()
	return txCoinbase
}

func (tx *SJB_Transaction) SJB_HashTransaction()  {

	result := tx.SJB_Serialize()
	resultBytes := bytes.Join([][]byte{SJB_IntToHex(time.Now().Unix()),result},[]byte{})

	hash := sha256.Sum256(resultBytes)

	tx.SJB_TxHash = hash[:]
}



func SJB_NewSimpleTransaction(from string, to string,amount int64,utxoSet *SJB_UTXOSet,txs []*SJB_Transaction,nodeId string) *SJB_Transaction {

	wallets,_ := SJB_NewWallets(nodeId)
	wallet := wallets.SJB_WalletsMap[from]
	if wallet == nil{
		fmt.Println("没有发送者私钥")
		os.Exit(1)
	}

	money,spendableUTXODic := utxoSet.SJB_FindSpendableUTXOS(from,amount,txs)

	var txIntputs []*SJB_TXInput
	var txOutputs []*SJB_TXOutput

	for txHash,indexArray := range spendableUTXODic  {

		txHashBytes,_  := hex.DecodeString(txHash)
		for _,index := range indexArray  {
			txInput := &SJB_TXInput{txHashBytes,index,nil,wallet.SJB_PublicKey}
			txIntputs = append(txIntputs,txInput)
		}

	}

	txOutput := SJB_NewTXOutput(int64(amount),to)
	txOutputs = append(txOutputs,txOutput)
	//找零
	txOutput = SJB_NewTXOutput(int64(money) - int64(amount),from)
	txOutputs = append(txOutputs,txOutput)

	tx := &SJB_Transaction{[]byte{},txIntputs,txOutputs}

	tx.SJB_HashTransaction()

	utxoSet.SJB_Blockchain.SJB_SignTransaction(tx, wallet.SJB_PrivateKey,txs)
	return tx

}


func (tx *SJB_Transaction) SJB_Hash() []byte {

	txCopy := tx

	txCopy.SJB_TxHash = []byte{}

	hash := sha256.Sum256(txCopy.SJB_Serialize())
	return hash[:]
}


func (tx *SJB_Transaction) SJB_Serialize() []byte {
	var encoded bytes.Buffer

	enc := gob.NewEncoder(&encoded)
	err := enc.Encode(tx)
	if err != nil {
		log.Panic(err)
	}

	return encoded.Bytes()
}

func SJB_DeserializeTransaction(txdata []byte) SJB_Transaction {

	var tx SJB_Transaction

	decoder := gob.NewDecoder(bytes.NewReader(txdata))
	err := decoder.Decode(&tx)
	if err != nil {
		log.Panic(err)
	}

	return tx
}

func (tx *SJB_Transaction) SJB_TrimmedCopy() SJB_Transaction {
	var inputs []*SJB_TXInput
	var outputs []*SJB_TXOutput

	for _, vin := range tx.SJB_Vins {
		inputs = append(inputs, &SJB_TXInput{vin.SJB_TxHash, vin.SJB_Vout, nil, nil})
	}

	for _, vout := range tx.SJB_Vouts {
		outputs = append(outputs, &SJB_TXOutput{vout.SJB_Value, vout.SJB_Ripemd160Hash})
	}

	txCopy := SJB_Transaction{tx.SJB_TxHash, inputs, outputs}

	return txCopy
}



func (tx *SJB_Transaction) SJB_Sign(privKey ecdsa.PrivateKey, prevTXs map[string]SJB_Transaction) {

	if tx.SJB_IsCoinbaseTransaction() {
		return
	}


	for _, vin := range tx.SJB_Vins {
		prevTx := prevTXs[hex.EncodeToString(vin.SJB_TxHash)]
		if prevTx.SJB_TxHash == nil {
			log.Panic("ERROR: Previous Transaction can not find")
		}
	}


	txCopy := tx.SJB_TrimmedCopy()

	for inID, vin := range txCopy.SJB_Vins {
		txCopy.SJB_Vins[inID].SJB_Signature = nil
		txCopy.SJB_Vins[inID].SJB_PublicKey = prevTXs[hex.EncodeToString(vin.SJB_TxHash)].SJB_Vouts[vin.SJB_Vout].SJB_Ripemd160Hash
		txCopy.SJB_TxHash = txCopy.SJB_Hash()

		r, s, err := ecdsa.Sign(rand.Reader, &privKey, txCopy.SJB_TxHash)
		if err != nil {
			log.Panic(err)
		}
		signature := append(r.Bytes(), s.Bytes()...)

		tx.SJB_Vins[inID].SJB_Signature = signature
	}
	return
}


func (tx *SJB_Transaction) SJB_Verify(prevTXs map[string]SJB_Transaction) bool {

	if tx.SJB_IsCoinbaseTransaction() {
		return true
	}

	for _, vin := range tx.SJB_Vins {
		prevTx := prevTXs[hex.EncodeToString(vin.SJB_TxHash)]
		if prevTx.SJB_TxHash == nil {
			log.Panic("ERROR: Previous Transaction can not find")
		}
	}

	txCopy := tx.SJB_TrimmedCopy()

	curve := elliptic.P256()

	for inID, vin := range tx.SJB_Vins {
		txCopy.SJB_Vins[inID].SJB_Signature = nil
		txCopy.SJB_Vins[inID].SJB_PublicKey = prevTXs[hex.EncodeToString(vin.SJB_TxHash)].SJB_Vouts[vin.SJB_Vout].SJB_Ripemd160Hash
		txCopy.SJB_TxHash = txCopy.SJB_Hash()
		//txCopy.SJB_Vins[inID].PublicKey = nil

		r := big.Int{}
		s := big.Int{}
		sigLen := len(vin.SJB_Signature)
		r.SetBytes(vin.SJB_Signature[:(sigLen / 2)])
		s.SetBytes(vin.SJB_Signature[(sigLen / 2):])

		x := big.Int{}
		y := big.Int{}
		keyLen := len(vin.SJB_PublicKey)
		x.SetBytes(vin.SJB_PublicKey[:(keyLen / 2)])
		y.SetBytes(vin.SJB_PublicKey[(keyLen / 2):])

		rawPubKey := ecdsa.PublicKey{curve, &x, &y}
		if ecdsa.Verify(&rawPubKey, txCopy.SJB_TxHash, &r, &s) == false {
			return false
		}
	}

	return true
}



