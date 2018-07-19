package BLC

import "bytes"

type SJB_TXInput struct {
	SJB_TxHash      []byte
	SJB_Vout      int
	SJB_Signature []byte // 数字签名
	SJB_PublicKey    []byte
}


func (txInput *SJB_TXInput) SJB_UnLockRipemd160Hash(ripemd160Hash []byte) bool {

	publicKey := SJB_Ripemd160Hash(txInput.SJB_PublicKey)

	if(bytes.Compare(publicKey,ripemd160Hash) == 0){
		return true
	}
	return false
}