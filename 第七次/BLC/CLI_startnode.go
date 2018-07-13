package BLC

import (
	"fmt"
	"os"
)

func (cli *SJB_CLI) startNode(nodeID string,minerAddress string)  {
	if minerAddress == "" || SJB_IsValidForAdress([]byte(minerAddress))  {
		fmt.Printf("启动服务器:localhost:%s\n",nodeID)
		SJB_startServer(nodeID,minerAddress)
	} else {
		fmt.Println("无效的挖矿收益地址")
		os.Exit(0)
	}
}