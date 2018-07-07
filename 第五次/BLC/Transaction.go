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

	txInput := &TXInput{[]byte{}, -1, nil,[]byte{}}
	txOutput := NewTXOutput(10,address)
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

	wallets,_ := NewWallets()
	wallet := wallets.WalletsMap[from]
	money,spendableUTXODic := blockchain.FindSpendableUTXOS(from,amount,txs)

	var txIntputs []*TXInput
	var txOutputs []*TXOutput

	for txHash,indexArray := range spendableUTXODic  {

		txHashBytes,_  := hex.DecodeString(txHash)
		for _,index := range indexArray  {
			txInput := &TXInput{txHashBytes,index,nil,wallet.PublicKey}
			txIntputs = append(txIntputs,txInput)
		}

	}

	txOutput := NewTXOutput(int64(amount),to)
	txOutputs = append(txOutputs,txOutput)

	txOutput = NewTXOutput(int64(money) - int64(amount),from)
	txOutputs = append(txOutputs,txOutput)

	tx := &Transaction{[]byte{},txIntputs,txOutputs}

	tx.HashTransaction()

	blockchain.SignTransaction(tx, wallet.PrivateKey)
	return tx

}


func (tx *Transaction) Hash() []byte {

	txCopy := tx

	txCopy.TxHash = []byte{}

	hash := sha256.Sum256(txCopy.Serialize())
	return hash[:]
}


func (tx *Transaction) Serialize() []byte {
	var encoded bytes.Buffer

	enc := gob.NewEncoder(&encoded)
	err := enc.Encode(tx)
	if err != nil {
		log.Panic(err)
	}

	return encoded.Bytes()
}


func (tx *Transaction) Sign(privKey ecdsa.PrivateKey, prevTXs map[string]Transaction) {

	if tx.IsCoinbaseTransaction() {
		return
	}


	for _, vin := range tx.Vins {
		if prevTXs[hex.EncodeToString(vin.TxHash)].TxHash == nil {
			log.Panic("ERROR: Previous transaction is not correct")
		}
	}


	txCopy := tx.TrimmedCopy()

	for inID, vin := range txCopy.Vins {
		prevTx := prevTXs[hex.EncodeToString(vin.TxHash)]
		txCopy.Vins[inID].Signature = nil
		txCopy.Vins[inID].PublicKey = prevTx.Vouts[vin.Vout].Ripemd160Hash
		txCopy.TxHash = txCopy.Hash()
		txCopy.Vins[inID].PublicKey = nil

		r, s, err := ecdsa.Sign(rand.Reader, &privKey, txCopy.TxHash)
		if err != nil {
			log.Panic(err)
		}
		signature := append(r.Bytes(), s.Bytes()...)

		tx.Vins[inID].Signature = signature
	}
}

func (tx *Transaction) TrimmedCopy() Transaction {
	var inputs []*TXInput
	var outputs []*TXOutput

	for _, vin := range tx.Vins {
		inputs = append(inputs, &TXInput{vin.TxHash, vin.Vout, nil, nil})
	}

	for _, vout := range tx.Vouts {
		outputs = append(outputs, &TXOutput{vout.Value, vout.Ripemd160Hash})
	}

	txCopy := Transaction{tx.TxHash, inputs, outputs}

	return txCopy
}


func (tx *Transaction) Verify(prevTXs map[string]Transaction) bool {
	if tx.IsCoinbaseTransaction() {
		return true
	}

	for _, vin := range tx.Vins {
		if prevTXs[hex.EncodeToString(vin.TxHash)].TxHash == nil {
			log.Panic("ERROR: Previous Transaction is wrong")
		}
	}

	txCopy := tx.TrimmedCopy()

	curve := elliptic.P256()

	for inID, vin := range tx.Vins {

		prevTx := prevTXs[hex.EncodeToString(vin.TxHash)]
		txCopy.Vins[inID].Signature = nil
		txCopy.Vins[inID].PublicKey = prevTx.Vouts[vin.Vout].Ripemd160Hash
		txCopy.TxHash = txCopy.Hash()
		txCopy.Vins[inID].PublicKey = nil

		a := big.Int{}
		b := big.Int{}
		sigLen := len(vin.Signature)
		a.SetBytes(vin.Signature[:(sigLen/2)])
		b.SetBytes(vin.Signature[(sigLen/2):])

		x := big.Int{}
		y := big.Int{}
		keyLen := len(vin.PublicKey)
		x.SetBytes(vin.PublicKey[:(keyLen/2)])
		y.SetBytes(vin.PublicKey[(keyLen/2):])

		rawPubKey := ecdsa.PublicKey{curve, &x, &y}
		if ecdsa.Verify(&rawPubKey, txCopy.TxHash, &a, &b) == false {
			return false
		}
	}

	return true
}



