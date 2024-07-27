package BLC

import (
	"fmt"
	"github.com/boltdb/bolt"
	"log"
	"math/big"
	"os"
	"time"
)

// 数据库名字
const dbName = "blockchain.db"

// bucket 名字
const blockBucketName = "blocks"

type BlockChain struct {
	Tip []byte   // 最新的区块的 Hash
	DB  *bolt.DB // 数据库
}

func (blc *BlockChain) Iterator() *BlockChainIterator {
	return &BlockChainIterator{blc.Tip, blc.DB}
}

// AddBlockToBlockchain 增加区块到区块链里面
func (blc *BlockChain) AddBlockToBlockchain(data string) {
	// 往链中添加区块
	err := blc.DB.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blockBucketName))
		if bucket != nil {
			// 从数据库中获取最新的区块
			latestBlockBytes := bucket.Get(blc.Tip)
			// 反序列化
			latestBlock := DeserializeBlock(latestBlockBytes)
			// 构造新区块
			newBlock := NewBlock(data, latestBlock.Height+1, latestBlock.Hash)
			// 把新区块放入数据库
			err := bucket.Put(newBlock.Hash, newBlock.Serialize())
			if err != nil {
				return fmt.Errorf("put bucket: %s", err)
			}
			blc.Tip = newBlock.Hash
		}
		return nil
	})
	if err != nil {
		log.Panic(err)
	}
}

// PrintBlockChain 遍历输出所有区块的信息
func (blc *BlockChain) PrintBlockChain() {
	blcIterator := blc.Iterator()
	for {
		block := blcIterator.Next()

		fmt.Printf("\nheight: %d, nonce: %d, timestamp: %s, prehash: %x, hash: %x, data: %s\n",
			block.Height,
			block.Nonce,
			time.Unix(block.Timestamp, 0).Format("2006-01-02 03:04:05 PM"),
			block.PrevBlockHash,
			block.Hash,
			block.Data,
		)

		var hashInt big.Int
		hashInt.SetBytes(block.PrevBlockHash)
		if big.NewInt(0).Cmp(&hashInt) == 0 {
			break
		}
	}
}

// CreateBlockchainWithGenesisBlock 创建带有创世区块的区块链
func CreateBlockchainWithGenesisBlock() *BlockChain {
	if dbExists() {
		fmt.Println("创世区块已经存在")

		db, err := bolt.Open(dbName, 0600, nil)
		if err != nil {
			log.Panic(err)
		}

		var blockChain *BlockChain

		err = db.View(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte(blockBucketName))
			hash := b.Get([]byte("latestHash"))
			blockChain = &BlockChain{hash, db}
			return nil
		})
		if err != nil {
			log.Panic(err)
		}

		return blockChain
	}

	db, err := bolt.Open(dbName, 0600, nil)
	if err != nil {
		log.Panic(err)
	}

	var blockHash []byte

	err = db.Update(func(tx *bolt.Tx) error {
		// 获取 bucket
		b := tx.Bucket([]byte(blockBucketName))
		if b == nil {
			// 创建 bucket
			b, err = tx.CreateBucket([]byte(blockBucketName))
			if err != nil {
				return fmt.Errorf("create bucket: %s", err)
			}
		}

		if b != nil {
			// 创建创世区块，并且存储到数据库中
			genesisBlock := CreateGenesisBlock("Genesis Data")
			err = b.Put(genesisBlock.Hash, genesisBlock.Serialize())
			if err != nil {
				return fmt.Errorf("put genesisBlock into bucket: %s", err)
			}

			// 存储最新的区块的 Hash 到数据库
			err = b.Put([]byte("latestHash"), genesisBlock.Hash)
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

// 判断数据库是否存在
func dbExists() bool {
	if _, err := os.Stat(dbName); os.IsNotExist(err) {
		return false
	}
	return true
}
