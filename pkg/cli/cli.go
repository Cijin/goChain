package cli

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/Cijin/gochain/pkg/blockchain"
	"github.com/Cijin/gochain/pkg/transaction"
)

type CLI struct{}

func (cli *CLI) printUsage() {
	fmt.Println("Usage:")
	fmt.Println("  printchain - print all the blocks of the blockchain")
	fmt.Println("  create -address <WALLET_ADDRESS> - create a new blockchain")
	fmt.Println("  getbalance -address <WALLET_ADDRESS> - get your current balance")
	fmt.Println("  send -from FROM -to TO -amount AMOUNT - Send AMOUNT of coins from FROM address to TO")
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

		fmt.Printf("Previous hash: %x\n", prevBlock.PrevBlockHash)
		fmt.Printf("Hash: %x\n", prevBlock.Hash)
		fmt.Println()

		if len(prevBlock.PrevBlockHash) == 0 {
			break
		}
	}
}

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

func (cli *CLI) send(from, to string, amount int) {
	bc := blockchain.NewBlockchain()
	defer bc.Db.Close()

	tx := blockchain.NewUnspentTxs(from, to, amount, bc)
	bc.MineBlock([]*transaction.Transaction{tx})
	fmt.Println("Success!")
}

// Run parses command line arguments and processes commands
func (cli *CLI) Run() {
	cli.validateArgs()

	createChainCmd := flag.NewFlagSet("create", flag.ExitOnError)
	sendCmd := flag.NewFlagSet("send", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)
	getbalanceCmd := flag.NewFlagSet("getbalance", flag.ExitOnError)

	createBlockchainAddress := createChainCmd.String("address", "", "The address to send genesis block reward to")
	getBalanceAddress := getbalanceCmd.String("address", "", "The address to get balance for")
	sendFrom := sendCmd.String("from", "", "Source wallet address")
	sendTo := sendCmd.String("to", "", "Destination wallet address")
	sendAmount := sendCmd.Int("amount", 0, "Amount to send")

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

	case "getbalance":
		err := getbalanceCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}

	case "send":
		err := sendCmd.Parse(os.Args[2:])
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

	if getbalanceCmd.Parsed() {
		if *getBalanceAddress == "" {
			getbalanceCmd.Usage()
			os.Exit(1)
		}
		cli.getBalance(*getBalanceAddress)
	}

	if printChainCmd.Parsed() {
		cli.printChain()
	}

	if sendCmd.Parsed() {
		if *sendFrom == "" || *sendTo == "" || *sendAmount <= 0 {
			sendCmd.Usage()
			os.Exit(1)
		}

		cli.send(*sendFrom, *sendTo, *sendAmount)
	}
}
