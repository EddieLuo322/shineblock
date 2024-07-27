package BLC

import (
	"github.com/boltdb/bolt"
	"log"
)

type BlockChainIterator struct {
	CurrentHash []byte
	DB          *bolt.DB
}

func (blcIterator *BlockChainIterator) Next() *Block {
	var block *Block

	err := blcIterator.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blockBucketName))
		if b != nil {
			currentBlockBytes := b.Get(blcIterator.CurrentHash)
			block = DeserializeBlock(currentBlockBytes)
			blcIterator.CurrentHash = block.PrevBlockHash
		}
		return nil
	})
	if err != nil {
		log.Panic(err)
	}
	return block
}
