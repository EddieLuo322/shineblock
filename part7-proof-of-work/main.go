package main

import (
	"fmt"
	"github.com/EddieLuo322/shineblock/part7-proof-of-work/BLC"
)

func main() {
	blockchain := BLC.CreateBlockchainWithGenesisBlock()
	blockchain.AddBlockToBlockchain("Send 1000RMB to Tony")
	//blockchain.AddBlockToBlockchain("Send 2000RMB to Sonny")
	fmt.Println(blockchain)
	fmt.Println(blockchain.Blocks)
}
