package BLC

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"log"
	"time"
)

type Block struct {
	// 区块高度
	Height int64
	// 当前区块 HASH 值
	Hash []byte
	// 上一个区块的 HASH 值
	PrevBlockHash []byte
	// 交易数据
	Txs []*Transaction
	// 时间戳
	Timestamp int64
	// Nonce 用于工作量整蒙
	Nonce int64
}

// Serialize 将区块序列化成字节数组
func (block *Block) Serialize() []byte {
	var result bytes.Buffer

	encoder := gob.NewEncoder(&result)
	err := encoder.Encode(block)
	if err != nil {
		log.Panic(err)
	}

	return result.Bytes()
}

// HashTransactions 将 transactions 转化为 []byte
func (block *Block) HashTransactions() []byte {
	var txHashes [][]byte
	var txHash [32]byte
	for _, tx := range block.Txs {
		//bytes.Join([][]byte{tx.TxHash}, []byte{})
		txHashes = append(txHashes, tx.TxHash)
	}
	txHash = sha256.Sum256(bytes.Join(txHashes, []byte{}))
	return txHash[:]
}

// DeserializeBlock 反序列化区块
func DeserializeBlock(blockBytes []byte) *Block {
	var block Block

	decoder := gob.NewDecoder(bytes.NewReader(blockBytes))
	err := decoder.Decode(&block)
	if err != nil {
		log.Panic(err)
	}

	return &block
}

// NewBlock 创建新的区块
func NewBlock(txs []*Transaction, height int64, prevBlockHash []byte) *Block {
	// 创建区块
	block := &Block{
		Height:        height,
		PrevBlockHash: prevBlockHash,
		Txs:           txs,
		Timestamp:     time.Now().Unix(),
		Hash:          nil,
		Nonce:         0,
	}

	// 调用工作量证明的方法，并且返回有效的 Hash 和 Nonce 值
	pow := NewProofOfWork(block)
	// 挖矿验证
	hash, nonce := pow.Run()
	block.Hash = hash[:]
	block.Nonce = nonce

	return block
}

// CreateGenesisBlock 单独写一个方法，生成创世区块
func CreateGenesisBlock(txs []*Transaction) *Block {
	return NewBlock(txs, 1, []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0})
}
