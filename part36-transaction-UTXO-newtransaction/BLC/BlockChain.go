package BLC

import (
	"fmt"
	"github.com/boltdb/bolt"
	"log"
	"math/big"
	"os"
	"strconv"
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
func (blc *BlockChain) AddBlockToBlockchain(txs []*Transaction) {
	// 往链中添加区块
	err := blc.DB.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blockBucketName))
		if bucket != nil {
			// 从数据库中获取最新的区块
			latestBlockBytes := bucket.Get(blc.Tip)
			// 反序列化
			latestBlock := DeserializeBlock(latestBlockBytes)
			// 构造新区块
			newBlock := NewBlock(txs, latestBlock.Height+1, latestBlock.Hash)
			// 把新区块放入数据库
			err := bucket.Put(newBlock.Hash, newBlock.Serialize())
			if err != nil {
				return fmt.Errorf("put newBlock bucket: %s", err)
			}
			// 更新数据库的最新hash latestHash
			err = bucket.Put([]byte("latestHash"), newBlock.Hash)
			if err != nil {
				return fmt.Errorf("put latestHash into bucket: %s", err)
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

		fmt.Printf("height: %d, nonce: %d, timestamp: %s, prehash: %x, hash: %x\n",
			block.Height,
			block.Nonce,
			time.Unix(block.Timestamp, 0).Format("2006-01-02 03:04:05 PM"),
			block.PrevBlockHash,
			block.Hash,
		)
		fmt.Println("Txs:")
		for _, tx := range block.Txs {
			fmt.Printf("%x\n", tx.TxHash)
			fmt.Println("Vins:")
			for _, in := range tx.Vins {
				fmt.Printf("%x\n", in.TxHash)
				fmt.Printf("%d\n", in.Vout)
				fmt.Printf("%s\n", in.ScriptSig)
			}
			fmt.Println("Vouts:")
			for _, out := range tx.Vouts {
				fmt.Printf("%d\n", out.Money)
				fmt.Printf("%s\n", out.ScriptPubKey)
			}
		}
		fmt.Println("--------------------------------------------------")

		var hashInt big.Int
		hashInt.SetBytes(block.PrevBlockHash)
		if big.NewInt(0).Cmp(&hashInt) == 0 {
			break
		}
	}
}

// MineNewBlock 挖掘新的区块
func (blc *BlockChain) MineNewBlock(from []string, to []string, amount []string) *BlockChain {
	fmt.Println(from)
	fmt.Println(to)
	fmt.Println(amount)

	value, _ := strconv.Atoi(amount[0])
	tx := NewSimpleTransaction(from[0], to[0], value)

	var txs []*Transaction
	txs = append(txs, tx)

	var block *Block
	err := blc.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blockBucketName))
		if b != nil {
			hash := b.Get([]byte("latestHash"))
			blockBytes := b.Get(hash)
			block = DeserializeBlock(blockBytes)
		}
		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	block = NewBlock(txs, block.Height+1, block.Hash)
	// 将新区块存储到数据库
	err = blc.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blockBucketName))
		if b != nil {
			err := b.Put(block.Hash, block.Serialize())
			if err != nil {
				return fmt.Errorf("put new block into bucket: %s", err)
			}
			err = b.Put([]byte("latestHash"), block.Hash)
			if err != nil {
				return fmt.Errorf("put latestHash into bucket: %s", err)
			}
			blc.Tip = block.Hash
		}
		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	return blc
}

// 如果一个地址对应的 TxOutput 未花费，那么这个 Transaction 就应该添加到数组中返回
func unSpentTransactionsWithAddress(address string) []*Transaction {
	return nil
}

// CreateBlockchainWithGenesisBlock 创建带有创世区块的区块链
func CreateBlockchainWithGenesisBlock(address string) *BlockChain {
	// 判断数据库是否存在
	if DBExists() {
		fmt.Println("数据库和创世区块已经存在")
		os.Exit(1)
	}

	fmt.Println("正在创建数据库，区块链和创世区块")

	// 创建或者打开数据库
	db, err := bolt.Open(dbName, 0600, nil)
	if err != nil {
		log.Panic(err)
	}

	var genesisHash []byte

	err = db.Update(func(tx *bolt.Tx) error {
		// 创建 bucket
		b, err := tx.CreateBucket([]byte(blockBucketName))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}

		// 数据库存储或者更新数据
		if b != nil {
			// 创建创世区块
			txCoinbase := NewCoinbaseTransaction(address)
			genesisBlock := CreateGenesisBlock([]*Transaction{txCoinbase})
			// 将创世区块存储到数据库
			err = b.Put(genesisBlock.Hash, genesisBlock.Serialize())
			if err != nil {
				return fmt.Errorf("put genesisBlock into bucket: %s", err)
			}

			// 存储创世区块的 Hash 到数据库
			err = b.Put([]byte("latestHash"), genesisBlock.Hash)
			if err != nil {
				return fmt.Errorf("put latestHash into bucket: %s", err)
			}
			genesisHash = genesisBlock.Hash
		}

		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	return &BlockChain{genesisHash, db}
}

// DBExists 判断数据库是否存在
func DBExists() bool {
	if _, err := os.Stat(dbName); os.IsNotExist(err) {
		return false
	}
	return true
}

func GetBlockChainFromBolt() *BlockChain {
	db, err := bolt.Open(dbName, 0600, nil)
	if err != nil {
		log.Panic(err)
	}

	var hash []byte

	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blockBucketName))

		if b != nil {
			hash = b.Get([]byte("latestHash"))
		}
		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	return &BlockChain{
		Tip: hash,
		DB:  db,
	}
}
