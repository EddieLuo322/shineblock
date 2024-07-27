package main

import (
	"github.com/EddieLuo322/shineblock/part24-persistence-cli/BLC"
)

func main() {
	blockchain := BLC.CreateBlockchainWithGenesisBlock()
	defer blockchain.DB.Close()

	cli := BLC.CLI{BlockChain: blockchain}
	cli.Run()
}
