package BLC

import "bytes"

type TXInput struct {
	TxHash      []byte
	Vout      int
	Signature []byte // 数字签名
	PublicKey    []byte
}


func (txInput *TXInput) UnLockRipemd160Hash(ripemd160Hash []byte) bool {

	publicKey := Ripemd160Hash(txInput.PublicKey)

	if(bytes.Compare(publicKey,ripemd160Hash) == 0){
		return true
	}
	return false
}