package cli

import (
	"fmt"
	"log"

	"github.com/Cijin/gochain/pkg/blockchain"
	"github.com/Cijin/gochain/pkg/transaction"
	"github.com/Cijin/gochain/pkg/utils"
)

func (cli *CLI) getBalance(address string) {
	if !transaction.ValidateAddress(address) {
		log.Panic("Error: Address is not valid")
	}

	bc := blockchain.NewBlockchain()
	defer bc.Db.Close()

	// get public key hash from address
	pubKeyHash := utils.Base58Decode([]byte(address))
	pubKeyHashLen := len(pubKeyHash) - transaction.ChecksumLen
	// element 0 is the version
	pubKeyHash = pubKeyHash[1:pubKeyHashLen]

	var balance int
	unspentTxOutputs := bc.FindUnspentTransactionOutputs(pubKeyHash)

	for _, out := range unspentTxOutputs {
		balance += out.Value
	}

	fmt.Printf("Balance of '%s': %d\n", address, balance)
}
