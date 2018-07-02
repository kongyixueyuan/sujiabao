package BLC

import (
	"bytes"
	"encoding/binary"
	"log"
	"encoding/json"
)


func IntToHex(num int64) []byte {
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, num)
	if err != nil {
		log.Panic(err)
	}

	return buff.Bytes()
}


func JSONToArray(jsonString string) []string {
	var sArr []string
	err := json.Unmarshal([]byte(jsonString),  &sArr);
	if  err != nil {
		log.Panic(err)
	}
	return sArr
}