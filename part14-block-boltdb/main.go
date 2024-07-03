package main

import (
	"fmt"
	"github.com/EddieLuo322/shineblock/part14-block-boltdb/BLC"
	"github.com/boltdb/bolt"
	"log"
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

	// 创建或者打开数据库
	db, err := bolt.Open("my.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("blocks"))

		if b == nil {
			b, err = tx.CreateBucket([]byte("blocks"))
			if err != nil {
				return fmt.Errorf("创建 bucket 失败: %s", err)
			}
		}

		// 往 bucket 里存储数据
		blockBytes := block.Serialize()
		err = b.Put([]byte("l"), blockBytes)
		if err != nil {
			return fmt.Errorf("数据存储失败: %s", err)
		}

		// 返回nil，以便数据库进行响应操作
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	// 查看
	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("blocks"))
		if b != nil {
			blockData := b.Get([]byte("l"))
			block1 := BLC.DeserializeBlock(blockData)
			fmt.Printf("%v\n", block1)
		}
		return nil
	})
	if err != nil {
		log.Panic(err)
	}
}
