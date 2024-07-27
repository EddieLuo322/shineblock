package BLC

import (
	"flag"
	"fmt"
	"log"
	"os"
)

type CLI struct {
}

func (cli *CLI) Run() {
	isValidArgs()
	addBlockCmd := flag.NewFlagSet("addblcok", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)
	createBlockChainCmd := flag.NewFlagSet("createblockchain", flag.ExitOnError)

	flagAddBlockData := addBlockCmd.String("data", "caobizhenshaung", "交易数据")
	flagCreateBlockChainWithData := createBlockChainCmd.String("data", "Genesis blockca data ...", "创世区块交易数据")

	switch os.Args[1] {
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
	case "createblockchain":
		err := createBlockChainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	default:
		printUsage()
		os.Exit(1)
	}

	if addBlockCmd.Parsed() {
		if *flagAddBlockData == "" {
			os.Exit(1)
		}
		cli.addBlock([]*Transaction{})
	}
	if printChainCmd.Parsed() {
		cli.printChain()
	}
	if createBlockChainCmd.Parsed() {
		if *flagCreateBlockChainWithData == "" {
			fmt.Println("交易数据不能为空")
			printUsage()
			os.Exit(1)
		}
		cli.createGenesisBlockChain([]*Transaction{})
	}
}

func (cli *CLI) addBlock(txs []*Transaction) {
	if DBExists() == false {
		fmt.Println("区块链不存在")
		os.Exit(1)
	}
	blockChain := GetBlockChainFromBolt()
	defer blockChain.DB.Close()
	blockChain.AddBlockToBlockchain(txs)
}

func (cli *CLI) printChain() {
	if DBExists() == false {
		fmt.Println("区块链不存在")
		os.Exit(1)
	}
	blockChain := GetBlockChainFromBolt()
	defer blockChain.DB.Close()
	blockChain.PrintBlockChain()
}

func (cli *CLI) createGenesisBlockChain(txs []*Transaction) {
	CreateBlockchainWithGenesisBlock(txs)
}

func printUsage() {
	fmt.Println("Usage: ")
	fmt.Println("\tcreateblockchain -data  -- 交易数据")
	fmt.Println("\taddblock -data DATA -- 交易数据")
	fmt.Println("\tprintchain -- 输出区块信息")
}

func isValidArgs() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}
}
