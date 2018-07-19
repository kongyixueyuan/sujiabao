package BLC

import (
	"fmt"
	"os"
	"flag"
	"log"
)

type SJB_CLI struct {}


func SJB_printUsage()  {


	fmt.Println("---------使用说明--------")
	fmt.Println("\t createblockchain -address -- 交易数据.")
	fmt.Println("\t send -from FROM -to TO -amount AMOUNT -mine -- 交易明细")

	fmt.Println("\t printchain -- 输出区块信息")
	fmt.Println("\t getbalance -address -- 输出区块信息.")

	fmt.Println("\t createwallet -- 创建钱包")
	fmt.Println("\t addresslists -- 输出所有钱包地址")

	fmt.Println("\t resetUTXO -- 重置.")
	fmt.Println("\tstartnode -miner Address -- 挖矿")
	fmt.Println("---------------------")
}

func SJB_isValidArgs()  {
	if len(os.Args) < 2 {
		SJB_printUsage()
		os.Exit(1)
	}
}



func (cli *SJB_CLI) SJB_Run()  {

	SJB_isValidArgs()

	nodeID := os.Getenv("NODE_ID")
	if nodeID == ""{
		println("请先设置节点")
		os.Exit(1)
	}


	//blockchain
	createBlockchainCmd := flag.NewFlagSet("createblockchain",flag.ExitOnError)
	flagCreateBlockchainWithAddress := createBlockchainCmd.String("address","","creat genesis blockchain")

	//打印已有区块
	printChainCmd := flag.NewFlagSet("printchain",flag.ExitOnError)

	//balance
	getbalanceCmd := flag.NewFlagSet("getbalance",flag.ExitOnError)
	getbalanceWithAdress := getbalanceCmd.String("address","","account address")

	//交易
	sendBlockCmd := flag.NewFlagSet("send",flag.ExitOnError)
	flagFrom := sendBlockCmd.String("from","","transfer frome......")
	flagTo := sendBlockCmd.String("to","","transfer to......")
	flagAmount := sendBlockCmd.String("amount","","transfer amount......")
	flagMine := sendBlockCmd.Bool("mine",false,"start mine?")

	//wallet
	addresslistsCmd := flag.NewFlagSet("addresslists",flag.ExitOnError)
	createWalletCmd := flag.NewFlagSet("createwallet",flag.ExitOnError)

	//reset UTXO
	resetUTXOCmd := flag.NewFlagSet("resetUTXO",flag.ExitOnError)

	//node
	startNodeCmd := flag.NewFlagSet("startnode",flag.ExitOnError)
	flagMiner := startNodeCmd.String("miner","","reward address.")


	switch os.Args[1] {
		case "send":
			err := sendBlockCmd.Parse(os.Args[2:])
			if err != nil {
				log.Panic(err)
			}
		case "printchain":
			err := printChainCmd.Parse(os.Args[2:])
			if err != nil {
				log.Panic(err)
			}
		case "createblockchain":
			err := createBlockchainCmd.Parse(os.Args[2:])
			if err != nil {
				log.Panic(err)
			}
		case "getbalance":
			err := getbalanceCmd.Parse(os.Args[2:])
			if err != nil {
				log.Panic(err)
			}
		case "addresslists":
			err := addresslistsCmd.Parse(os.Args[2:])
			if err != nil {
				log.Panic(err)
			}
		case "createwallet":
			err := createWalletCmd.Parse(os.Args[2:])
			if err != nil {
				log.Panic(err)
			}
		case "resetUTXO":
			err := resetUTXOCmd.Parse(os.Args[2:])
			if err != nil {
				log.Panic(err)
			}
		case "startnode":
			err := startNodeCmd.Parse(os.Args[2:])
			if err != nil {
				log.Panic(err)
			}
		default:
			SJB_printUsage()
			os.Exit(1)
	}


	if sendBlockCmd.Parsed() {
		if *flagFrom == "" || *flagTo == "" || *flagAmount == ""{
			SJB_printUsage()
			os.Exit(1)
		}

		from := SJB_JSONToArray(*flagFrom)
		to := SJB_JSONToArray(*flagTo)
		amount := SJB_JSONToArray(*flagAmount)
		cli.SJB_send(from,to,amount,nodeID,*flagMine)
	}

	if printChainCmd.Parsed() {
		cli.SJB_printchain(nodeID)
	}

	if createBlockchainCmd.Parsed() {

		if *flagCreateBlockchainWithAddress == "" {
			fmt.Println("please input address....")
			SJB_printUsage()
			os.Exit(1)
		}

		cli.SJB_createGenesisBlockchain(*flagCreateBlockchainWithAddress,nodeID)
	}

	if getbalanceCmd.Parsed() {

		if *getbalanceWithAdress == "" {
			fmt.Println("please input address....")
			SJB_printUsage()
			os.Exit(1)
		}

		cli.SJB_getBalance(*getbalanceWithAdress,nodeID)
	}

	if addresslistsCmd.Parsed() {
		cli.SJB_addressLists(nodeID)
	}


	if createWalletCmd.Parsed() {
		cli.SJB_createWallet(nodeID)
	}

	if resetUTXOCmd.Parsed(){
		cli.SJB_resetUTXOSet(nodeID)
	}

	if startNodeCmd.Parsed(){
		cli.startNode(nodeID,*flagMiner)
	}
}