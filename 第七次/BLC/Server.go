package BLC

import (
	"fmt"
	"net"
	"log"
	"io/ioutil"
)

func SJB_startServer(nodeID string,minerAddress string)  {

	bc := SJB_BlockchainObject(nodeID)
	defer bc.SJB_DB.Close()

	nodeAddress = fmt.Sprintf("localhost:%s",nodeID)
	ln,err := net.Listen(PROTOCOL,nodeAddress)
	if err != nil {
		log.Panic(err)
	}
	defer ln.Close()

	if SJB_IsValidForAdress([]byte(minerAddress)) {
		MyminerAddress = minerAddress
	}
	if nodeAddress != knowNodes[0]{
		 SJB_sendVersion(knowNodes[0],bc)
	}

	//accept other node message
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Panic(err)
		}
		go SJB_handleConnection(conn,bc)
	}
}


func SJB_handleConnection(conn net.Conn,bc *SJB_Blockchain) {
	request, err := ioutil.ReadAll(conn)
	if err != nil {
		log.Panic(err)
	}

	fmt.Printf("Receive a Message:%s\n",request[:COMMANDLENGTH])

	command := SJB_BytesToCommand(request[:COMMANDLENGTH])

	switch command {
		case COMMAND_VERSION:
			SJB_handleVersion(request, bc)
		case COMMAND_ADDR:
			SJB_handleAddr(request, bc)
		case COMMAND_BLOCK:
			SJB_handleBlock(request, bc)
		case COMMAND_GETBLOCKS:
			SJB_handleGetblocks(request, bc)
		case COMMAND_GETDATA:
			SJB_handleGetData(request, bc)
		case COMMAND_INV:
			SJB_handleInv(request, bc)
		case COMMAND_TX:
			SJB_handleTx(request, bc)
		default:
			fmt.Println("Unknown command!")
	}

	conn.Close()
}







