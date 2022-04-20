package cli

import (
	"fmt"

	"github.com/Cijin/gochain/pkg/blockchain"
)

func (cli *CLI) createBlockchain(address string) {
	bc := blockchain.CreateBlockchain(address)
	bc.Db.Close()
	fmt.Println("Blockchain created!")
}
