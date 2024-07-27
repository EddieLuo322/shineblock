package BLC

type TXInput struct {
	// 交易的 Hash
	TxHash []byte
	// 存储 TXOutput 在 Vout 里面的索引
	Vout int
	// 用户名
	ScriptSig string
}

// UnLockTXInputWithAddress 判断当前的消费是谁的钱
func (txInput *TXInput) UnLockTXInputWithAddress(address string) bool {
	return txInput.ScriptSig == address
}
