package BLC

type BlockChain struct {
	Blocks []*Block // 存储有序的区块

}

// AddBlockToBlockchain 增加区块到区块链里面
func (blc *BlockChain) AddBlockToBlockchain(data string) {
	// 创建新区块
	newBlock := NewBlock(
		data,
		blc.Blocks[len(blc.Blocks)-1].Height+1,
		blc.Blocks[len(blc.Blocks)-1].Hash,
	)
	// 往链中添加区块
	blc.Blocks = append(blc.Blocks, newBlock)
}

// CreateBlockchainWithGenesisBlock 创建带有创世区块的区块链
func CreateBlockchainWithGenesisBlock() *BlockChain {
	genesisBlock := CreateGenesisBlock("Genesis Data")
	return &BlockChain{
		[]*Block{genesisBlock},
	}
}
