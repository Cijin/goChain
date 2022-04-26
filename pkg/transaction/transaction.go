package transaction

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
)

const subsidy = 50

type Transaction struct {
	Id   []byte
	Vin  []TXInput
	Vout []TXOutput
}

func (tx Transaction) Serialize() []byte {
	var encoded bytes.Buffer

	enc := gob.NewEncoder(&encoded)
	err := enc.Encode(tx)
	if err != nil {
		log.Panic(err)
	}

	return encoded.Bytes()
}

func (tx *Transaction) Hash() []byte {
	var hash [32]byte

	txCopy := *tx
	txCopy.Id = []byte{}

	hash = sha256.Sum256(txCopy.Serialize())

	return hash[:]
}

func NewCoinbaseTX(to, data string) *Transaction {
	if data == "" {
		data = fmt.Sprintf("Sending Reward to '%s'", to)
	}

	txin := TXInput{[]byte{}, -1, nil, []byte(data)}
	txout := NewTXOutput(subsidy, to)
	tx := Transaction{nil, []TXInput{txin}, []TXOutput{*txout}}
	tx.Hash()

	return &tx
}

func (tx *Transaction) IsCoinbaseTx() bool {
	return len(tx.Vin) == 1 && len(tx.Vin[0].TxId) == 0 && tx.Vin[0].Vout == -1
}

func (tx *Transaction) Sign(privateKey ecdsa.PrivateKey, prevTxs map[string]Transaction) {
	// Coinbase tx's are not signed cause there are no inputs assosiated with them
	if tx.IsCoinbaseTx() {
		return
	}

	for _, vin := range tx.Vin {
		if prevTxs[hex.EncodeToString(vin.TxId)].Id == nil {
			log.Panic("Error: Previous transaction is not valid")
		}
	}

	trimmedTx := tx.TrimmedTransaction()

	for inIdx, vin := range trimmedTx.Vin {
		prevTx := prevTxs[hex.EncodeToString(vin.TxId)]
		// set public key to pubKeyHash of ref. output
		trimmedTx.Vin[inIdx].PubKey = prevTx.Vout[inIdx].PubKeyHash
		trimmedTx.Id = trimmedTx.Hash()
		// reset public key so that it does not interfere with the next iteration
		trimmedTx.Vin[inIdx].PubKey = nil

		signature, err := ecdsa.SignASN1(rand.Reader, &privateKey, trimmedTx.Id)
		if err != nil {
			log.Panic(err)
		}

		tx.Vin[inIdx].Signature = signature
	}
}

func (tx *Transaction) TrimmedTransaction() Transaction {
	var inputs []TXInput
	var outputs []TXOutput

	for _, vin := range tx.Vin {
		inputs = append(inputs, TXInput{vin.TxId, vin.Vout, nil, nil})
	}

	for _, vout := range tx.Vout {
		outputs = append(outputs, vout)
	}

	trimmedTx := Transaction{tx.Id, inputs, outputs}

	return trimmedTx
}

func (tx *Transaction) Verify(prevTxs map[string]Transaction) bool {
	curve := elliptic.P256()
	trimmedTx := tx.TrimmedTransaction()

	for inIdx, vin := range trimmedTx.Vin {
		prevTx := prevTxs[hex.EncodeToString(vin.TxId)]
		// set public key to pubKeyHash of ref. output
		trimmedTx.Vin[inIdx].PubKey = prevTx.Vout[inIdx].PubKeyHash
		trimmedTx.Id = trimmedTx.Hash()
		// reset public key so that it does not interfere with the next iteration
		trimmedTx.Vin[inIdx].PubKey = nil

		x := big.Int{}
		y := big.Int{}
		keyLen := len(vin.PubKey)

		x.SetBytes(vin.PubKey[:(keyLen / 2)])
		y.SetBytes(vin.PubKey[(keyLen / 2):])
		rawPubKey := ecdsa.PublicKey{curve, &x, &y}

		if !ecdsa.VerifyASN1(&rawPubKey, trimmedTx.Id, vin.Signature) {
			return false
		}
	}

	return true
}
