package BLC

import "fmt"

func (cli *CLI) getBalance(address string) {
	blockchain := GetBlockChainFromBolt()
	defer blockchain.DB.Close()
	amount := blockchain.GetBalance(address)
	fmt.Printf("%s balance： %d", address, amount)
}
