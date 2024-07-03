package BLC

type BlockChain struct {
	Blocks []*Block // 存储有序的区块

}

// CreateBlockchainWithGenesisBlock 创建带有创世区块的区块链
func CreateBlockchainWithGenesisBlock() *BlockChain {
	genesisBlock := CreateGenesisBlock("Genesis Data")
	return &BlockChain{
		[]*Block{genesisBlock},
	}
}
