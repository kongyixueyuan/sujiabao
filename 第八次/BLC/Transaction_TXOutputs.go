package BLC

import (
	"bytes"
	"encoding/gob"
	"log"
)

type SJB_TXOutputs struct {
	SJB_UTXOS []*SJB_UTXO
}


// 将区块序列化成字节数组
func (txOutputs *SJB_TXOutputs) SJB_Serialize() []byte {

	var result bytes.Buffer

	encoder := gob.NewEncoder(&result)

	err := encoder.Encode(txOutputs)
	if err != nil {
		log.Panic(err)
	}

	return result.Bytes()
}

// 反序列化
func SJB_DeserializeTXOutputs(txOutputsBytes []byte) *SJB_TXOutputs {

	var txOutputs SJB_TXOutputs

	decoder := gob.NewDecoder(bytes.NewReader(txOutputsBytes))
	err := decoder.Decode(&txOutputs)
	if err != nil {
		log.Panic(err)
	}

	return &txOutputs
}