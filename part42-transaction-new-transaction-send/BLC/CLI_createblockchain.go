package BLC

func (cli *CLI) createGenesisBlockChain(address string) {
	blockchain := CreateBlockchainWithGenesisBlock(address)
	defer blockchain.DB.Close()
}
