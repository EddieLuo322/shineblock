package BLC

type TXOutput struct {
	Money        int64
	ScriptPubKey string // 公钥
}

// UnLockScriptPubKeyWithAddress 解锁
func (txOutput *TXOutput) UnLockScriptPubKeyWithAddress(address string) bool {
	return txOutput.ScriptPubKey == address
}
