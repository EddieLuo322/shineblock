package BLC

import (
	"flag"
	"fmt"
	"log"
	"os"
)

type CLI struct {
	BlockChain *BlockChain
}

func (cli *CLI) Run() {
	isValidArgs()
	addBlockCmd := flag.NewFlagSet("addblcok", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)
	createBlockChainCmd := flag.NewFlagSet("createblockchain", flag.ExitOnError)

	flagAddBlockData := addBlockCmd.String("data", "caobizhenshaung", "交易数据")
	flagCreateBlockChainWithData := createBlockChainCmd.String("data", "Genesis data ...", "创世区块交易数据")
	//flagPrintChain := printChainCmd.Bool("open", true, "打印区块链")
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
		cli.addBlock(*flagAddBlockData)
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
		cli.createGenesisBlockChain(*flagCreateBlockChainWithData)
	}
}

func (cli *CLI) addBlock(data string) {
	cli.BlockChain.AddBlockToBlockchain(data)
}

func (cli *CLI) printChain() {
	cli.BlockChain.PrintBlockChain()
}

func (cli *CLI) createGenesisBlockChain(data string) {
	fmt.Println(data)
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
