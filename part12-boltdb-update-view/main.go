package main

import (
	"fmt"
	"github.com/boltdb/bolt"
	"log"
)

func main() {
	// 创建或者打开数据库
	db, err := bolt.Open("my.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Update(func(tx *bolt.Tx) error {
		// 获取 bucket
		b := tx.Bucket([]byte("MyBucket"))

		// 往表里存储数据
		if b != nil {
			err = b.Put([]byte("ll"), []byte("Send 10 BTC to BingBing"))
			if err != nil {
				return fmt.Errorf("数据存储失败: %s", err)
			}
		}

		// 返回nil，以便数据库进行响应操作
		return nil
	})
	// 更新失败
	if err != nil {
		log.Fatal(err)
	}
}
