package BLC

import (
	"fmt"
	"os"
)

func (cli *CLI) printChain() {
	if DBExists() == false {
		fmt.Println("区块链不存在")
		os.Exit(1)
	}
	blockChain := GetBlockChainFromBolt()
	defer blockChain.DB.Close()
	blockChain.PrintBlockChain()
}
