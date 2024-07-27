package BLC

import (
	"fmt"
	"os"
)

func (cli *CLI) send(from []string, to []string, amount []string) {
	if DBExists() == false {
		fmt.Println("区块链不存在")
		os.Exit(1)
	}
	blockchain := GetBlockChainFromBolt()
	defer blockchain.DB.Close()
	blockchain.MineNewBlock(from, to, amount)
}
