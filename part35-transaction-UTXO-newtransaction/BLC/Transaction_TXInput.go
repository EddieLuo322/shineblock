package BLC

type TXInput struct {
	// 交易的 Hash
	TxHash []byte
	// 存储 TXOutput 在 Vout 里面的索引
	Vout int
	// 用户名
	ScriptSig string
}
