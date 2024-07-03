package main

import (
	"fmt"
	"github.com/EddieLuo322/shineblock/part2-basic-prototype/BLC"
)

func main() {
	genesisBlockchain := BLC.CreateBlockchainWithGenesisBlock()
	fmt.Println(genesisBlockchain)
	fmt.Println(genesisBlockchain.Blocks[0])
}
