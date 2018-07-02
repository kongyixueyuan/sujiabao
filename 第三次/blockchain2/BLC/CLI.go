package BLC

import (
	"fmt"
	"os"
	"flag"
	"log"
)

type CLI struct {}


func printUsage()  {

	fmt.Println("Usage: ")
	
	fmt.Println("\tcreateblockchain -data -- 交易数据.")
	fmt.Println("\taddblock -data DATA -- 交易数据.")
	fmt.Println("\tprintchain -- 输出区块信息.")

}

func isValidArgs()  {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}
}

func (cli *CLI) addBlock(data string)  {

	if DBExists() == false {
		fmt.Println("........请先新建链.......")
		os.Exit(1)
	}

	blockchain := BlockchainObject()
	defer blockchain.DB.Close()

	blockchain.AddBlockToBlockchain(data)
}

func (cli *CLI) printchain()  {

	if DBExists() == false {
		fmt.Println("未找到区块链")
		os.Exit(1)
	}


	blockchain := BlockchainObject()

	defer blockchain.DB.Close()

	blockchain.PrintAllchain()

}

func (cli *CLI) createGenesisBlockchain(data string)  {

	CreateBlockchainWithGenesisBlock(data)
}


func (cli *CLI) Run()  {

	isValidArgs()

	addBlockCmd := flag.NewFlagSet("addblock",flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("printchain",flag.ExitOnError)
	createBlockchainCmd := flag.NewFlagSet("createblockchain",flag.ExitOnError)

	flagAddBlockData := addBlockCmd.String("data","nothing here","交易数据......")

	flagCreateBlockchainWithData := createBlockchainCmd.String("data","birthday","创世区块交易数据......")


	switch os.Args[1] {
		case "createblockchain":
			err := createBlockchainCmd.Parse(os.Args[2:])
			if err != nil {
				log.Panic(err)
			}
		case "addblock":
			err := addBlockCmd.Parse(os.Args[2:])
			if err != nil {
				log.Panic(err)
			}
		case "printchain":
			err := printChainCmd.Parse(os.Args[2:])
			if err != nil {
				log.Panic(err)
			}
		default:
			printUsage()
			os.Exit(1)
	}

	if addBlockCmd.Parsed() {
		if *flagAddBlockData == "" {
			printUsage()
			os.Exit(1)
		}
		cli.addBlock(*flagAddBlockData)
	}

	if printChainCmd.Parsed() {
		cli.printchain()
	}

	if createBlockchainCmd.Parsed() {

		if *flagCreateBlockchainWithData == "" {
			fmt.Println("交易数据不能为空......")
			printUsage()
			os.Exit(1)
		}

		cli.createGenesisBlockchain(*flagCreateBlockchainWithData)
	}

}