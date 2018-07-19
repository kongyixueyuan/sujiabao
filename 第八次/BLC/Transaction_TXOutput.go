package BLC

import "bytes"

type SJB_TXOutput struct {
	SJB_Value int64
	SJB_Ripemd160Hash []byte
}


func (txOutput *SJB_TXOutput)  SJB_Lock(address string)  {

	publicKeyHash := SJB_Base58Decode([]byte(address))
	txOutput.SJB_Ripemd160Hash = publicKeyHash[1:len(publicKeyHash) - 4]

}


func SJB_NewTXOutput(value int64,address string) *SJB_TXOutput {

	txOutput := &SJB_TXOutput{value,nil}
	txOutput.SJB_Lock(address)

	return txOutput
}

func (txOutput *SJB_TXOutput) SJB_UnLockScriptPubKeyWithAddress(address string) bool {

	publickeyHash := SJB_Base58Decode([]byte(address))
	hash_160 := publickeyHash[1:len(publickeyHash)-4]

	if(bytes.Compare(hash_160,txOutput.SJB_Ripemd160Hash) == 0){
		return true
	}

	return false
}


