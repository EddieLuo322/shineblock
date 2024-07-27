package BLC

type TXOutput struct {
	Money        int64
	ScriptPubKey string // 公钥
}

// UnLockTXOutputWithAddress 解锁
func (txOutput *TXOutput) UnLockTXOutputWithAddress(address string) bool {
	return txOutput.ScriptPubKey == address
}
