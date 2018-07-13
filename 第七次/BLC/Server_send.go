package BLC

import (
	"fmt"
	"io"
	"bytes"
	"log"
	"net"
)


func SJB_sendVersion(toAddress string,bc *SJB_Blockchain)  {
	bestHeight := bc.GetBestHeight()
	payload := SJB_GobEncode(Version{NODE_VERSION, bestHeight, nodeAddress})
	request := append(SJB_CommandToBytes(COMMAND_VERSION), payload...)
	SJB_sendData(toAddress,request)
}

func SJB_sendGetBlocks(toAddress string)  {
	payload := SJB_GobEncode(GetBlocks{nodeAddress})
	request := append(SJB_CommandToBytes(COMMAND_GETBLOCKS), payload...)
	SJB_sendData(toAddress,request)
}

func SJB_sendInv(toAddress string, messagetype string, hashes [][]byte) {
	payload := SJB_GobEncode(Inv{nodeAddress,messagetype,hashes})
	request := append(SJB_CommandToBytes(COMMAND_INV), payload...)
	SJB_sendData(toAddress,request)
}

func SJB_sendGetData(toAddress string, messagetype string ,blockHash []byte) {
	payload := SJB_GobEncode(GetData{nodeAddress,messagetype,blockHash})
	request := append(SJB_CommandToBytes(COMMAND_GETDATA), payload...)
	SJB_sendData(toAddress,request)
}

func SJB_sendBlock(toAddress string, block []byte)  {
	payload := SJB_GobEncode(BlockData{nodeAddress,block})
	request := append(SJB_CommandToBytes(COMMAND_BLOCK), payload...)
	SJB_sendData(toAddress,request)
}

func SJB_sendTx(toAddress string, txn *SJB_Transaction){
	data := TxData{nodeAddress, txn.SJB_Serialize()}
	payload := SJB_GobEncode(data)
	request := append(SJB_CommandToBytes(COMMAND_TX), payload...)
	SJB_sendData(toAddress,request)
}

func SJB_sendData(to string,data []byte)  {
	fmt.Printf("向%s发送数据......",to)
	conn, err := net.Dial("tcp", to)
	if err != nil {
		log.Panic(err)
	}
	defer conn.Close()

	_, err = io.Copy(conn, bytes.NewReader(data))
	if err != nil {
		log.Panic(err)
	}
}