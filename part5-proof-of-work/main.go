package main

import (
	"fmt"
	"github.com/EddieLuo322/shineblock/part5-proof-of-work/BLC"
)

func main() {
	blockchain := BLC.CreateBlockchainWithGenesisBlock()
	fmt.Println(blockchain)
	fmt.Println(blockchain.Blocks)
	fmt.Println(blockchain.Blocks[0])

	blockchain.AddBlockToBlockchain("Send 1000RMB to Tony",
		blockchain.Blocks[len(blockchain.Blocks)-1].Height+1,
		blockchain.Blocks[len(blockchain.Blocks)-1].Hash)
	blockchain.AddBlockToBlockchain("Send 2000RMB to Sonny",
		blockchain.Blocks[len(blockchain.Blocks)-1].Height+1,
		blockchain.Blocks[len(blockchain.Blocks)-1].Hash)
	fmt.Println(blockchain)
	fmt.Println(blockchain.Blocks)
}
