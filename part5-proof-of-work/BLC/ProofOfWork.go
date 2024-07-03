package BLC

type ProofOfWork struct {
	Block *Block
}

func (proofOfWork ProofOfWork) Run() ([]byte, int64) {
	return nil, 0
}

func NewProofOfWork(block *Block) *ProofOfWork {
	return &ProofOfWork{
		block,
	}
}
