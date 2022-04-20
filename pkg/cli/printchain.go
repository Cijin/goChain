package cli

import (
	"fmt"

	"github.com/Cijin/gochain/pkg/blockchain"
)

func (cli *CLI) printChain() {
	bc := blockchain.NewBlockchain()
	bci := bc.Iterator()

	for {
		prevBlock := bci.Previous()

		fmt.Printf("Previous hash: %x\n", prevBlock.PrevBlockHash)
		fmt.Printf("Hash: %x\n", prevBlock.Hash)
		fmt.Println()

		if len(prevBlock.PrevBlockHash) == 0 {
			break
		}
	}
}
