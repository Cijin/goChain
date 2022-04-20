package cli

import (
	"fmt"

	"github.com/Cijin/gochain/pkg/blockchain"
)

func (cli *CLI) getBalance(address string) {
	bc := blockchain.NewBlockchain()
	defer bc.Db.Close()

	var balance int
	unspentTxOutputs := bc.FindUnspentTransactionOutputs(address)

	for _, out := range unspentTxOutputs {
		balance += out.Value
	}

	fmt.Printf("Balance of '%s': %d\n", address, balance)
}
