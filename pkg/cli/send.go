package cli

import (
	"fmt"

	"github.com/Cijin/gochain/pkg/blockchain"
	"github.com/Cijin/gochain/pkg/transaction"
)

func (cli *CLI) send(from, to string, amount int) {
	bc := blockchain.NewBlockchain()
	defer bc.Db.Close()

	tx := blockchain.NewUnspentTxs(from, to, amount, bc)
	bc.MineBlock([]*transaction.Transaction{tx})
	fmt.Println("Success!")
}
