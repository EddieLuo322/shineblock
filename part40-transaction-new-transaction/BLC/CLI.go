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

	printChainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)
	createBlockChainCmd := flag.NewFlagSet("createblockchain", flag.ExitOnError)
	sendCmd := flag.NewFlagSet("send", flag.ExitOnError)
	getBalanceCmd := flag.NewFlagSet("getbalance", flag.ExitOnError)

	flagFrom := sendCmd.String("from", "", "转账转出地址")
	flagTo := sendCmd.String("to", "", "转账转入地址")
	flagAmount := sendCmd.String("amount", "", "转账金额")
	flagCreateBlockChainWithAddress := createBlockChainCmd.String("address", "", "创建创世区块的地址")
	getBalanceWithAddress := getBalanceCmd.String("address", "", "要查询的地址")

	switch os.Args[1] {
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
	case "send":
		err := sendCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "getbalance":
		err := getBalanceCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	default:
		printUsage()
		os.Exit(1)
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
	if sendCmd.Parsed() {
		if *flagFrom == "" || *flagTo == "" || *flagAmount == "" {
			printUsage()
			os.Exit(1)
		}
		//fmt.Println(*flagFrom)
		//fmt.Println(*flagTo)
		//fmt.Println(*flagAmount)
		//fmt.Println(JsonToArray(*flagFrom))
		//fmt.Println(JsonToArray(*flagTo))
		//fmt.Println(JsonToArray(*flagAmount))
		//cli.addBlock([]*Transaction{})

		from := JsonToArray(*flagFrom)
		to := JsonToArray(*flagTo)
		amount := JsonToArray(*flagAmount)
		cli.send(from, to, amount)
	}
	if getBalanceCmd.Parsed() {
		if *getBalanceWithAddress == "" {
			printUsage()
			os.Exit(1)
		}
		cli.getBalance(*getBalanceWithAddress)
	}
}

func printUsage() {
	fmt.Println("Usage: ")
	fmt.Println("\tcreateblockchain -address  -- 交易地址")
	fmt.Println("\tsend -from FROM -to TO -amount AMOUNT -- 交易明细")
	fmt.Println("\tprintchain -- 输出区块信息")
	fmt.Println("\tgetbalance -address -- 输出地址余额")
}

func isValidArgs() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}
}
