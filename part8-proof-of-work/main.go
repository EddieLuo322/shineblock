package main

import (
	"fmt"
	"github.com/EddieLuo322/shineblock/part8-proof-of-work/BLC"
)

func main() {
	//blockchain := BLC.CreateBlockchainWithGenesisBlock()
	//blockchain.AddBlockToBlockchain("Send 1000RMB to Tony")
	////blockchain.AddBlockToBlockchain("Send 2000RMB to Sonny")
	//fmt.Println(blockchain)
	//fmt.Println(blockchain.Blocks)

	block := BLC.NewBlock("test", 1, []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0})
	fmt.Printf("\n%d\n", block.Nonce)
	fmt.Printf("%x\n", block.Hash)
}
