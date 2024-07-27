package BLC

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"log"
)

// Transaction UTXO
type Transaction struct {
	TxHash []byte
	Vins   []*TXInput
	Vouts  []*TXOutput
}

// IsCoinbaseTransaction 判断当前交易是否是 Coinbase 交易
func (tx *Transaction) IsCoinbaseTransaction() bool {
	return len(tx.TxHash) == 0 && tx.Vins[0].Vout == -1
}

// Transaction 创建分两种情况

// NewCoinbaseTransaction 1. 创世区块创建时的 Transaction
func NewCoinbaseTransaction(address string) *Transaction {
	txInput := &TXInput{[]byte{}, -1, "Genesis Data"}
	txOutput := &TXOutput{10, address}
	txCoinbase := &Transaction{[]byte{}, []*TXInput{txInput}, []*TXOutput{txOutput}}
	txCoinbase.HashTransaction()
	return txCoinbase
}

// NewSimpleTransaction 2. 转账时产生的 Transaction
func NewSimpleTransaction(from string, to string, amount int) *Transaction {
	var txInputs []*TXInput
	var txOutputs []*TXOutput

	// 将 []byte 直接转化成的 hex字符串 再次直接转化为 []byte
	// 注意：[]byte, string, hex string 三者之间的关系
	// []byte 和 string 是对应关系
	// []byte 数据用 hex(16进制)表示，就是 hex string
	// 原来的 tx hash 是 []byte，EncodeString后表示成 hex string传过来，在DecodeString成原来的 []byte
	decodeString, _ := hex.DecodeString("5c7d7b92659749335f82907da6469a2e7f89c12df8dea0ac03b7b779ce4c9abf")
	txInput := &TXInput{decodeString, 0, from}
	txInputs = append(txInputs, txInput)

	txOutput := &TXOutput{int64(amount), to}
	txOutputs = append(txOutputs, txOutput)
	txOutput = &TXOutput{10 - int64(amount), from}
	txOutputs = append(txOutputs, txOutput)

	txCoinbase := &Transaction{[]byte{}, txInputs, txOutputs}
	txCoinbase.HashTransaction()
	return txCoinbase
}

func (tx *Transaction) HashTransaction() {
	var result bytes.Buffer

	encoder := gob.NewEncoder(&result)
	err := encoder.Encode(tx)
	if err != nil {
		log.Panic(err)
	}

	hash := sha256.Sum256(result.Bytes())
	tx.TxHash = hash[:]
}
