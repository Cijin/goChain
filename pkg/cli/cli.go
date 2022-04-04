package cli

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/Cijin/gochain/pkg/blockchain"
)

type CLI struct{}

func (cli *CLI) printUsage() {
	fmt.Println("Usage:")
	fmt.Println("  printchain - print all the blocks of the blockchain")
	fmt.Println("  create -address <WALLET_ADDRESS> - create a new blockchain")
}

func (cli *CLI) validateArgs() {
	if len(os.Args) < 2 {
		cli.printUsage()
		os.Exit(1)
	}
}

func (cli *CLI) createBlockchain(address string) {
	bc := blockchain.CreateBlockchain(address)
	bc.Db.Close()
	fmt.Println("Blockchain created!")
}

func (cli *CLI) printChain() {
	bc := blockchain.NewBlockchain()
	bci := bc.Iterator()

	for {
		prevBlock := bci.Previous()

		fmt.Printf("Prev. hash: %x\n", prevBlock.PrevBlockHash)
		fmt.Printf("Hash: %x\n", prevBlock.Hash)
		fmt.Println()

		if len(prevBlock.PrevBlockHash) == 0 {
			break
		}
	}
}

// Run parses command line arguments and processes commands
func (cli *CLI) Run() {
	cli.validateArgs()

	createChainCmd := flag.NewFlagSet("create", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)

	createBlockchainAddress := createChainCmd.String("address", "", "The address to send genesis block reward to")

	switch os.Args[1] {
	case "create":
		err := createChainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "printchain":
		err := printChainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	default:
		cli.printUsage()
		os.Exit(1)
	}

	if createChainCmd.Parsed() {
		if *createBlockchainAddress == "" {
			createChainCmd.Usage()
			os.Exit(1)
		}
		cli.createBlockchain(*createBlockchainAddress)
	}

	if printChainCmd.Parsed() {
		cli.printChain()
	}
}
