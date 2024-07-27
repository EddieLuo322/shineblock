package main

import (
	"github.com/EddieLuo322/shineblock/part18-persistence-iterator-time-format/BLC"
)

func main() {
	blockchain := BLC.CreateBlockchainWithGenesisBlock()
	defer blockchain.DB.Close()
	blockchain.AddBlockToBlockchain("Send 1000RMB to Tony")
	blockchain.AddBlockToBlockchain("Send 2000RMB to Sonny")
	blockchain.AddBlockToBlockchain("Send 3000RMB to Eddie")
	blockchain.PrintBlockChain()
}
