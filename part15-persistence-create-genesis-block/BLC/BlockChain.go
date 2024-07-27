package BLC

import (
	"fmt"
	"github.com/boltdb/bolt"
	"log"
)

// 数据库名字
const dbName = "blockchain.db"

// bucket 名字
const blockBucketName = "blocks"

type BlockChain struct {
	Tip []byte   // 最新的区块的 Hash
	DB  *bolt.DB // 数据库
}

// AddBlockToBlockchain 增加区块到区块链里面
//func (blc *BlockChain) AddBlockToBlockchain(data string) {
//	// 创建新区块
//	newBlock := NewBlock(
//		data,
//		blc.Blocks[len(blc.Blocks)-1].Height+1,
//		blc.Blocks[len(blc.Blocks)-1].Hash,
//	)
//	// 往链中添加区块
//	blc.Blocks = append(blc.Blocks, newBlock)
//}

// CreateBlockchainWithGenesisBlock 创建带有创世区块的区块链
func CreateBlockchainWithGenesisBlock() *BlockChain {
	// 创建或者打开数据库
	db, err := bolt.Open(dbName, 0600, nil)
	if err != nil {
		log.Panic(err)
	}

	var blockHash []byte

	err = db.Update(func(tx *bolt.Tx) error {
		// 创建 bucket
		b, err := tx.CreateBucket([]byte(blockBucketName))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}

		if b != nil {
			genesisBlock := CreateGenesisBlock("Genesis Data")
			err = b.Put(genesisBlock.Hash, genesisBlock.Serialize())
			if err != nil {
				return fmt.Errorf("put genesisBlock into bucket: %s", err)
			}

			// 存储最新的区块的 Hash
			err = b.Put([]byte("latestHash"), genesisBlock.Serialize())
			if err != nil {
				return fmt.Errorf("put latestHash into bucket: %s", err)
			}

			blockHash = genesisBlock.Hash
		}

		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	return &BlockChain{
		blockHash,
		db,
	}
}
