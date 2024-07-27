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
	flagCreateBlockChainWithAddress := createBlockChainCmd.String("address", "", "创建创世区块的地址")

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
		if *flagCreateBlockChainWithAddress == "" {
			fmt.Println("地址不能为空")
			printUsage()
			os.Exit(1)
		}
		cli.createGenesisBlockChain(*flagCreateBlockChainWithAddress)
	}
}

func (cli *CLI) createGenesisBlockChain(address string) {
	CreateBlockchainWithGenesisBlock(address)
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

func printUsage() {
	fmt.Println("Usage: ")
	fmt.Println("\tcreateblockchain -address  -- 交易数据")
	fmt.Println("\taddblock -data DATA -- 交易数据")
	fmt.Println("\tprintchain -- 输出区块信息")
}

func isValidArgs() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}
}
