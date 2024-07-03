package BLC

import "math/big"

// 生成Hash的难度，256位 Hash 前面至少要有16个 0
const targetBit = 16

type ProofOfWork struct {
	Block  *Block   // 当前要验证的区块
	target *big.Int // 验证难度
}

func (proofOfWork ProofOfWork) Run() ([]byte, int64) {
	return nil, 0
}

func NewProofOfWork(block *Block) *ProofOfWork {
	target := big.NewInt(1)
	target = target.Lsh(target, 256-targetBit)

	return &ProofOfWork{
		block,
		target,
	}
}
