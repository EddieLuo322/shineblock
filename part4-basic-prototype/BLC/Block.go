package BLC

import (
	"bytes"
	"crypto/sha256"
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
}

func (block *Block) SetHash() {
	// 将 Height 转化为 []byte 类型
	heightBytes := Int64ToBytes(block.Height)
	// 将 TimeStamp 转化为 []byte 类型
	timeBytes := Int64ToBytes(block.Timestamp)
	// 拼接
	blockBytes := bytes.Join([][]byte{heightBytes, block.PrevBlockHash, block.Data, timeBytes}, []byte{})
	// 生成 Hash
	hash := sha256.Sum256(blockBytes)

	block.Hash = hash[:]
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
	}
	// 设置 Hash
	block.SetHash()
	return block
}

// CreateGenesisBlock 单独写一个方法，生成创世区块
func CreateGenesisBlock(data string) *Block {
	return NewBlock(data, 1, []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0})
}
