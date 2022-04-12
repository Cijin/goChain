package blockchain

import (
	"encoding/hex"
	"log"

	"github.com/Cijin/gochain/pkg/transaction"
)

/*
 * Parse through each block in blockchain
 *	Within the block parse through each transaction
 *   If tx is not coinbase tx
 *   Parse through TXInput's:
 *    - If input has same address as the func param
 *				* Append to array in map with current txId as key
 *				* This will be used to check if output is unspent later
 *
 *
 *   Parse through TXOutput's:
 *	   - If current tx has unspent outputs
 *				 * Check if output is spent, using the map from above
 *				 * If yes continue to next output
 *
 *
 *  If the output can now be unlocked with address, append to
 *  unspent tx's
 *
 *  If all blocks have been traversed, break
 *
 *  Return any unspent transactions
 *
 */
func (bc *Blockchain) FindUnspentTransactions(address string) []transaction.Transaction {
	var unspentTxs []transaction.Transaction
	spentTXs := make(map[string][]int)
	bcI := bc.Iterator()

	for {
		block := bcI.Previous()

		for _, tx := range block.Transactions {
			txId := hex.EncodeToString(tx.Id)

			// look for spent tx's
			if !tx.IsCoinbaseTx() {
				for _, in := range tx.Vin {
					if in.CanUnlockOutputWith(address) {
						txId := hex.EncodeToString(in.TxId)
						spentTXs[txId] = append(spentTXs[txId], in.Vout)
					}
				}
			}

		Outputs:
			for _, out := range tx.Vout {
				for _, spentOut := range spentTXs[txId] {
					// check if current output spent
					if spentOut == out.Value {
						continue Outputs
					}
				}

				if out.CanBeUnlockedWith(address) {
					// found unspent tx
					unspentTxs = append(unspentTxs, *tx)
				}
			}
		}

		if len(block.PrevBlockHash) == 0 {
			break
		}
	}

	return unspentTxs
}

/*
 * Find unspent transaction outputs utilizes the above function to return only
 * outputs, which will make finding balances easier
 */
func (bc *Blockchain) FindUnspentTransactionOutputs(address string) []transaction.TXOutput {
	var unspentTxOutputs []transaction.TXOutput
	unspentTxs := bc.FindUnspentTransactions(address)

	for _, tx := range unspentTxs {
		for _, out := range tx.Vout {
			if out.CanBeUnlockedWith(address) {
				unspentTxOutputs = append(unspentTxOutputs, out)
			}
		}
	}

	return unspentTxOutputs
}

/*
 * amount required for the current transaction
 * @returns spendable outputs, accumulatedAmount
 */
func (bc *Blockchain) FindSpendableTxOutputs(address string, amount int) (map[string][]int, int) {
	var spendableOutputs = make(map[string][]int)
	accumulatedAmount := 0
	unspentTxs := bc.FindUnspentTransactions(address)

	for _, tx := range unspentTxs {
		txId := hex.EncodeToString(tx.Id)

		for outIdx, output := range tx.Vout {

			accumulatedAmount += output.Value
			spendableOutputs[txId] = append(spendableOutputs[txId], outIdx)

			if accumulatedAmount >= amount {
				break
			}
		}
	}

	return spendableOutputs, accumulatedAmount
}

/*
 * Generate Inputs, Output (change if neccessary)
 *  Inputs:
 *		* Get "spendable" outputs
 *		* Generate inputs using TxId
 *		* For now outputs are just indexes
 *		* ScriptSig -> from
 */
func NewUnspentTxs(from, to string, amount int, bc *Blockchain) *transaction.Transaction {
	var inputs []transaction.TXInput
	var outputs []transaction.TXOutput

	spendableOutputs, balance := bc.FindSpendableTxOutputs(from, amount)

	if balance < amount {
		log.Panic("Error: Insufficent balance")
	}

	// inputs
	for txid, output := range spendableOutputs {
		txId, err := hex.DecodeString(txid)

		if err != nil {
			log.Panic(err)
		}

		for _, out := range output {
			inputs = append(inputs, transaction.TXInput{
				TxId:      txId,
				Vout:      out,
				ScriptSig: from,
			})
		}
	}

	// outputs
	outputs = append(outputs, transaction.TXOutput{
		Value: amount, ScriptPubKey: to,
	})

	// generate change
	if amount < balance {
		outputs = append(outputs, transaction.TXOutput{
			Value:        balance - amount,
			ScriptPubKey: from,
		})
	}

	tr := transaction.Transaction{
		Id:   nil,
		Vin:  inputs,
		Vout: outputs,
	}
	tr.SetId()

	return &tr
}
