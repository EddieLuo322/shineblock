package main

import (
	"fmt"
	"github.com/EddieLuo322/shineblock/part16-persistence-add-block/BLC"
)

func main() {
	blockchain := BLC.CreateBlockchainWithGenesisBlock()
	defer blockchain.DB.Close()
	blockchain.AddBlockToBlockchain("Send 1000RMB to Tony")
	blockchain.AddBlockToBlockchain("Send 2000RMB to Sonny")
	blockchain.AddBlockToBlockchain("Send 3000RMB to Eddie")
	fmt.Println(blockchain, blockchain.Tip, blockchain.DB)
}
