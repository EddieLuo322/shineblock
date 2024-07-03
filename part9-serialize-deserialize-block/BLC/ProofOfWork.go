package BLC

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math/big"
)

// 生成Hash的难度，256位 Hash 前面至少要有16个 0
const targetBit = 16

type ProofOfWork struct {
	Block  *Block   // 当前要验证的区块
	target *big.Int // 验证难度
}

// IsValid 判断 block 是否完成了验证
func (proofOfWork ProofOfWork) IsValid() bool {
	var hashInt big.Int

	hashInt.SetBytes(proofOfWork.Block.Hash)
	if proofOfWork.target.Cmp(&hashInt) == 1 {
		return true
	}
	return false
}

func (proofOfWork ProofOfWork) Run() ([]byte, int64) {
	nonce := int64(0)

	var hash [32]byte
	var hashInt big.Int // 存储新生成的 hash

	for {
		// 将 Block 的属性拼接成字节数组
		dataBytes := proofOfWork.prepareData(nonce)
		// 生成 hash
		hash = sha256.Sum256(dataBytes)
		fmt.Printf("\r%x", hash)
		// 将 hash 存储到 hashInt
		hashInt.SetBytes(hash[:])
		// 判断 hashInt 是否小于 Block 里面的 target
		if proofOfWork.target.Cmp(&hashInt) == 1 {
			break
		}
		nonce += 1
	}

	return hash[:], nonce
}

func (proofOfWork ProofOfWork) prepareData(nonce int64) []byte {
	data := bytes.Join(
		[][]byte{
			proofOfWork.Block.PrevBlockHash,
			proofOfWork.Block.Data,
			Int64ToBytes(proofOfWork.Block.Timestamp),
			Int64ToBytes(proofOfWork.Block.Height),
			Int64ToBytes(int64(targetBit)),
			Int64ToBytes(nonce),
		},
		[]byte{},
	)
	return data
}

func NewProofOfWork(block *Block) *ProofOfWork {
	target := big.NewInt(1)
	target = target.Lsh(target, 256-targetBit)

	return &ProofOfWork{
		block,
		target,
	}
}
