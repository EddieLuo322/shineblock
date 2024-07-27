package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

func printUsage() {
	fmt.Println("Usage: ")
	fmt.Println("\taddblock -data DATA -- 交易数据")
	fmt.Println("\tprintchain -- 输出区块信息")
}

func isValidArgs() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}
}

func main() {
	isValidArgs()
	addBlockCmd := flag.NewFlagSet("addblcok", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)
	flagAddBlockData := addBlockCmd.String("data", "caobizhenshaung", "交易数据")
	switch os.Args[1] {
	case "addblock":
		err := addBlockCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "printchain":
		err := addBlockCmd.Parse(os.Args[2:])
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
		fmt.Println(*flagAddBlockData)
	}
	if printChainCmd.Parsed() {
		fmt.Println("输出所有区块的数据")
	}
}

// go build -o bc main.go

// bc
// ./bc addBlock -data "send $100 to YaDi"

// ./bc printchain
