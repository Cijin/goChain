package main

import (
	"fmt"

	"github.com/Cijin/gochain/pkg/blockchain"
)

func main() {

	bc := blockchain.NewBlockchain()

	// close db when main exits
	defer bc.Db.Close()

	bc.AddBlock("Send 1 GoCoin to Seagin")
	bc.AddBlock("Send 2 GoCoin to Seagin")

	for _, block := range bc.Blocks {
		fmt.Printf("Data: %s\n", block.Data)
		fmt.Printf("Prev Hash: %x\n", block.PrevBlockHash)
		fmt.Printf("Hash: %x\n", block.Hash)
		fmt.Println("")
	}

}
