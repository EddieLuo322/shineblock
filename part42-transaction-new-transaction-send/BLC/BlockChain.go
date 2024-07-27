package BLC

import (
	"encoding/hex"
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
			fmt.Printf("tx.TxHash: %x\n", tx.TxHash)
			fmt.Println("tx.Vins:")
			for _, in := range tx.Vins {
				fmt.Printf("in.TxHash: %x, in.Vout： %d， in.ScriptSig： %s\n", in.TxHash, in.Vout, in.ScriptSig)
			}
			fmt.Println("tx.Vouts:")
			for _, out := range tx.Vouts {
				fmt.Printf("out.Money: %d, out.ScriptPubKey: %s\n", out.Money, out.ScriptPubKey)
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
	fmt.Println("MineNewBlock start ...")

	var txs []*Transaction
	for i := 0; i < len(from); i++ {
		value, _ := strconv.Atoi(amount[i])
		tx := NewSimpleTransaction(from[i], to[i], value, blc, txs)
		txs = append(txs, tx)
	}

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
	fmt.Println("newblock: ")
	fmt.Printf("height: %d, nonce: %d, timestamp: %s, prehash: %x, hash: %x\n",
		block.Height,
		block.Nonce,
		time.Unix(block.Timestamp, 0).Format("2006-01-02 03:04:05 PM"),
		block.PrevBlockHash,
		block.Hash,
	)
	fmt.Println("Txs:")
	for _, tx := range block.Txs {
		fmt.Printf("tx.TxHash: %x\n", tx.TxHash)
		fmt.Println("tx.Vins:")
		for _, in := range tx.Vins {
			fmt.Printf("in.TxHash: %x, in.Vout： %d， in.ScriptSig： %s\n", in.TxHash, in.Vout, in.ScriptSig)
		}
		fmt.Println("tx.Vouts:")
		for _, out := range tx.Vouts {
			fmt.Printf("out.Money: %d, out.ScriptPubKey: %s\n", out.Money, out.ScriptPubKey)
		}
	}
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
	fmt.Println("MineNewBlock end ...")
	return blc
}

// FindSpendableUTXOs 转账时查找可用的UTXO
func (blc *BlockChain) FindSpendableUTXOs(from string, amount int, txs []*Transaction) (int64, map[string][]int) {
	fmt.Println("FindSpendableUTXOs start ...")
	utxos := blc.GetUTXOs(from, txs)
	spendableUTXOs := make(map[string][]int)

	var value int64
	for _, utxo := range utxos {
		value += utxo.Output.Money

		hash := hex.EncodeToString(utxo.TxHash)
		spendableUTXOs[hash] = append(spendableUTXOs[hash], utxo.Index)

		if value >= int64(amount) {
			break
		}
	}

	if value < int64(amount) {
		fmt.Printf("%s's fund is not enough", from)
		os.Exit(1)
	}
	fmt.Println("spendableUTXOs: ")
	for key, val := range spendableUTXOs {
		fmt.Println(key, val)
	}
	fmt.Println("FindSpendableUTXOs end ...")
	return value, spendableUTXOs
}

// GetUTXOs 如果一个地址对应的 TxOutput 未花费，那么这个 Transaction 就应该添加到数组中返回
func (blc *BlockChain) GetUTXOs(address string, txs []*Transaction) []*UTXO {
	fmt.Println("GetUTXOs start ...")

	// 用来存储遍历出来的 未花费的 TXOutput
	var UTXOs []*UTXO

	spentTXOutputs := make(map[string][]int)

	for _, tx := range txs {
		if tx.IsCoinbaseTransaction() != true {
			for _, in := range tx.Vins {
				// 是否能否解锁
				if in.UnLockTXInputWithAddress(address) {
					key := hex.EncodeToString(in.TxHash)
					spentTXOutputs[key] = append(spentTXOutputs[key], in.Vout)
				}
			}
		}
	}

	for _, tx := range txs {
	work1:
		for index, out := range tx.Vouts {
			if out.UnLockTXOutputWithAddress(address) {
				if len(spentTXOutputs) != 0 {
					for txHash, indexArray := range spentTXOutputs {
						if txHash == hex.EncodeToString(tx.TxHash) {
							for _, i := range indexArray {
								if index == i {
									continue work1
								}
							}
						}
					}
					utxo := &UTXO{TxHash: tx.TxHash, Index: index, Output: out}
					UTXOs = append(UTXOs, utxo)
				} else {
					utxo := &UTXO{TxHash: tx.TxHash, Index: index, Output: out}
					UTXOs = append(UTXOs, utxo)
				}
			}
		}
	}

	blockIterator := blc.Iterator()

	for {
		block := blockIterator.Next()

		for i := len(block.Txs) - 1; i >= 0; i-- {
			tx := block.Txs[i]

			if tx.IsCoinbaseTransaction() != true {
				for _, in := range tx.Vins {
					// 是否能否解锁
					if in.UnLockTXInputWithAddress(address) {
						key := hex.EncodeToString(in.TxHash)
						spentTXOutputs[key] = append(spentTXOutputs[key], in.Vout)
					}
				}
			}

		work:
			for index, out := range tx.Vouts {
				if out.UnLockTXOutputWithAddress(address) {
					if spentTXOutputs != nil {
						if len(spentTXOutputs) != 0 {
							// 对于一个 output
							// 如果本次 tx 的 hash 在spentTXOutputs 字典中
							// 并且本次 output 在 tx outputs 中的 index， 在 spentTXOutputs 字典里 txHash对应的 indexArray 中
							// 那么本 output 则已经被花费
							// 针对双重循环的命中，有 flag 和 continue 两种处理方法
							for txHash, indexArray := range spentTXOutputs {
								for _, i := range indexArray {
									if txHash == hex.EncodeToString(tx.TxHash) && index == i {
										continue work
									}
								}
							}
							utxo := &UTXO{TxHash: tx.TxHash, Index: index, Output: out}
							UTXOs = append(UTXOs, utxo)
						} else {
							utxo := &UTXO{TxHash: tx.TxHash, Index: index, Output: out}
							UTXOs = append(UTXOs, utxo)
						}
					}
				}
			}
		}

		// 遍历区块的所有 Txs
		// 迭代器从后往前遍历
		// 因为 input 一定是在 output 后面的区块中，所以我们可以在一次遍历中进行 input, output 的判断
		// 如果一笔 output 被花费了，对应的 input 一定记录在他之后的区块中，所以当我们从后往前遍历，不断记录 input，就可以通过记录判断一笔
		// output 是否被花费
		//for _, tx := range block.Txs {
		//	// 不能是 创世区块的第一个 tx
		//	// 遍历每一个 tx 的 input
		//	// 如果 input 的地址是本地址，那么说明对应的 output 被花费了
		//	if tx.IsCoinbaseTransaction() != true {
		//		for _, in := range tx.Vins {
		//			// 是否能否解锁
		//			if in.UnLockTXInputWithAddress(address) {
		//				key := hex.EncodeToString(in.TxHash)
		//				spentTXOutputs[key] = append(spentTXOutputs[key], in.Vout)
		//			}
		//		}
		//	}
		//}

		// 遍历每一个 tx
		// 如果 output 的对象是本地址
		// 查看记录字典 spentTXOutputs，查看这笔 output 是否被花费，也就是被 input 作为源
		// 如果没有被花费，记录在 unUTXOs 中，表示这一笔 output 没有被花费
		//for _, tx := range block.Txs {
		//work:
		//	for index, out := range tx.Vouts {
		//		if out.UnLockTXOutputWithAddress(address) {
		//			if spentTXOutputs != nil {
		//				if len(spentTXOutputs) != 0 {
		//					// 对于一个 output
		//					// 如果本次 tx 的 hash 在spentTXOutputs 字典中
		//					// 并且本次 output 在 tx outputs 中的 index， 在 spentTXOutputs 字典里 txHash对应的 indexArray 中
		//					// 那么本 output 则已经被花费
		//					// 针对双重循环的命中，有 flag 和 continue 两种处理方法
		//					for txHash, indexArray := range spentTXOutputs {
		//						for _, i := range indexArray {
		//							if txHash == hex.EncodeToString(tx.TxHash) && index == i {
		//								continue work
		//							}
		//						}
		//					}
		//					utxo := &UTXO{TxHash: tx.TxHash, Index: index, Output: out}
		//					UTXOs = append(UTXOs, utxo)
		//				} else {
		//					utxo := &UTXO{TxHash: tx.TxHash, Index: index, Output: out}
		//					UTXOs = append(UTXOs, utxo)
		//				}
		//			}
		//		}
		//	}
		//}

		var hashInt big.Int
		hashInt.SetBytes(block.PrevBlockHash[:])
		if hashInt.Cmp(big.NewInt(0)) == 0 {
			break
		}
	}
	fmt.Println("UTXOs:")
	for _, utxo := range UTXOs {
		fmt.Println(hex.EncodeToString(utxo.TxHash), utxo.Index)
		fmt.Println("output:", utxo.Output.Money, utxo.Output.ScriptPubKey)
	}
	fmt.Println("GetUTXOs end ...")
	return UTXOs
}

func (blc *BlockChain) GetBalance(address string) int64 {
	fmt.Printf("get balance start, address=%s ...\n", address)
	var amount int64
	utxos := blc.GetUTXOs(address, []*Transaction{})
	for _, utxo := range utxos {
		amount += utxo.Output.Money
	}
	fmt.Println("get balance end...")
	return amount
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
