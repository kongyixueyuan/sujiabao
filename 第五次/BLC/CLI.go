package BLC

import (
	"fmt"
	"os"
	"flag"
	"log"
)

type CLI struct {}


func printUsage()  {


	fmt.Println("---------使用说明--------")
	fmt.Println("\tcreateblockchain -address -- 交易数据.")
	fmt.Println("\tsend -from FROM -to TO -amount AMOUNT -- 交易明细")
	fmt.Println("\tcreatewallet -- 创建钱包")
	fmt.Println("\tprintchain -- 输出区块信息")
	fmt.Println("\tgetbalance -address -- 输出区块信息.")
	fmt.Println("\taddresslists -- 输出所有钱包地址")
	fmt.Println("---------------------")

}

func isValidArgs()  {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}
}



func (cli *CLI) Run()  {

	isValidArgs()

	sendBlockCmd := flag.NewFlagSet("send",flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("printchain",flag.ExitOnError)
	createBlockchainCmd := flag.NewFlagSet("createblockchain",flag.ExitOnError)
	getbalanceCmd := flag.NewFlagSet("getbalance",flag.ExitOnError)

	flagFrom := sendBlockCmd.String("from","","transfer frome......")
	flagTo := sendBlockCmd.String("to","","transfer to......")
	flagAmount := sendBlockCmd.String("amount","","transfer amount......")

	flagCreateBlockchainWithAddress := createBlockchainCmd.String("address","","creat genesis blockchain")
	getbalanceWithAdress := getbalanceCmd.String("address","","account address")


	addresslistsCmd := flag.NewFlagSet("addresslists",flag.ExitOnError)
	createWalletCmd := flag.NewFlagSet("createwallet",flag.ExitOnError)




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
		default:
			printUsage()
			os.Exit(1)
	}

	if sendBlockCmd.Parsed() {
		if *flagFrom == "" || *flagTo == "" || *flagAmount == ""{
			printUsage()
			os.Exit(1)
		}

		from := JSONToArray(*flagFrom)
		to := JSONToArray(*flagTo)
		amount := JSONToArray(*flagAmount)
		cli.send(from,to,amount)
	}

	if printChainCmd.Parsed() {
		cli.printchain()
	}

	if createBlockchainCmd.Parsed() {

		if *flagCreateBlockchainWithAddress == "" {
			fmt.Println("please input address....")
			printUsage()
			os.Exit(1)
		}

		cli.createGenesisBlockchain(*flagCreateBlockchainWithAddress)
	}

	if getbalanceCmd.Parsed() {

		if *getbalanceWithAdress == "" {
			fmt.Println("please input address....")
			printUsage()
			os.Exit(1)
		}

		cli.getBalance(*getbalanceWithAdress)
	}

	if addresslistsCmd.Parsed() {
		cli.addressLists()
	}


	if createWalletCmd.Parsed() {
		cli.createWallet()
	}

}