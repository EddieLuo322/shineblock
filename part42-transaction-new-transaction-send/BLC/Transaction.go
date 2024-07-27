package BLC

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
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
func NewSimpleTransaction(from string, to string, amount int, blockchain *BlockChain, txs []*Transaction) *Transaction {
	fmt.Println("NewSimpleTransaction start ...")
	fmt.Printf("from: %s, to: %s, amount: %d\n", from, to, amount)

	money, spendableUTXODic := blockchain.FindSpendableUTXOs(from, amount, txs)
	fmt.Printf("moeny: %d\n", money)
	fmt.Println("spendableUTXODic: ")
	for key, value := range spendableUTXODic {
		fmt.Printf("%s: %x\n", key, value)
	}

	var txInputs []*TXInput
	var txOutputs []*TXOutput

	// 将 []byte 直接转化成的 hex字符串 再次直接转化为 []byte
	// 注意：[]byte, string, hex string 三者之间的关系
	// []byte 和 string 是对应关系
	// []byte 数据用 hex(16进制)表示，就是 hex string
	// 原来的 tx hash 是 []byte，EncodeString后表示成 hex string传过来，在DecodeString成原来的 []byte
	for txHash, indexArray := range spendableUTXODic {
		txHashBytes, _ := hex.DecodeString(txHash)
		for _, index := range indexArray {
			txInput := &TXInput{txHashBytes, index, from}
			txInputs = append(txInputs, txInput)
		}
	}

	txOutput := &TXOutput{int64(amount), to}
	txOutputs = append(txOutputs, txOutput)
	txOutput = &TXOutput{money - int64(amount), from}
	txOutputs = append(txOutputs, txOutput)

	tx := &Transaction{[]byte{}, txInputs, txOutputs}
	// 设置 tx hash 值
	tx.HashTransaction()
	fmt.Println("new tx: ")
	fmt.Printf("tx.TxHash: %x\n", tx.TxHash)
	fmt.Println("tx.Vins:")
	for _, in := range tx.Vins {
		fmt.Printf("in.TxHash: %x, in.Vout： %d， in.ScriptSig： %s\n", in.TxHash, in.Vout, in.ScriptSig)
	}
	fmt.Println("tx.Vouts:")
	for _, out := range tx.Vouts {
		fmt.Printf("out.Money: %d, out.ScriptPubKey: %s\n", out.Money, out.ScriptPubKey)
	}
	fmt.Println("NewSimpleTransaction end ...")
	return tx
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
