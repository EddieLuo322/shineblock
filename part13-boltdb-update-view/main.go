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

	// 创建表
	err = db.Update(func(tx *bolt.Tx) error {
		// 创建 bucket, 相当于表
		b := tx.Bucket([]byte("MyBucket"))

		// 往表里存储数据
		if b != nil {
			data := b.Get([]byte("l"))
			fmt.Println(string(data))
			data = b.Get([]byte("ll"))
			fmt.Println(string(data))
		}

		// 返回nil，以便数据库进行响应操作
		return nil
	})
	// 更新失败
	if err != nil {
		log.Fatal(err)
	}
}
