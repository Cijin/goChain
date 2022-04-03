package main

import (
	"github.com/Cijin/gochain/pkg/blockchain"
	"github.com/Cijin/gochain/pkg/cli"
)

func main() {
	bc := blockchain.NewBlockchain()
	defer bc.Db.Close()

	cli := cli.CLI{
		Bc: bc,
	}
	cli.Run()
}
