package BLC

import (
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
	Data []byte
	// 时间戳
	Timestamp int64
	// Nonce 用于工作量整蒙
	Nonce int64
}

// NewBlock 创建新的区块
func NewBlock(data string, height int64, prevBlockHash []byte) *Block {
	// 创建区块
	block := &Block{
		Height:        height,
		PrevBlockHash: prevBlockHash,
		Data:          []byte(data),
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
func CreateGenesisBlock(data string) *Block {
	return NewBlock(data, 1, []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0})
}
