package BLC

import "bytes"

type TXOutput struct {
	Value int64
	Ripemd160Hash []byte
}


func (txOutput *TXOutput)  Lock(address string)  {

publicKeyHash := Base58Decode([]byte(address))

txOutput.Ripemd160Hash = publicKeyHash[1:len(publicKeyHash) - 4]
}


func NewTXOutput(value int64,address string) *TXOutput {

	txOutput := &TXOutput{value,nil}
	txOutput.Lock(address)

	return txOutput
}

func (txOutput *TXOutput) UnLockScriptPubKeyWithAddress(address string) bool {

	publickeyHash := Base58Decode([]byte(address))
	hash_160 := publickeyHash[1:len(publickeyHash)-4]

	if(bytes.Compare(hash_160,txOutput.Ripemd160Hash) == 0){
		return true
	}

	return false
}


