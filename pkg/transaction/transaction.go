package transaction

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
)

const subsidy = 50

type TXOutput struct {
	Value        int
	ScriptPubKey string
}

type TXInput struct {
	TxId      []byte
	Vout      int
	ScriptSig string
}

type Transaction struct {
	Id   []byte
	Vin  []TXInput
	Vout []TXOutput
}

func (tx *Transaction) SetId() {
	var encoded bytes.Buffer
	var hash [32]byte

	enc := gob.NewEncoder(&encoded)
	enc.Encode(tx)

	hash = sha256.Sum256(encoded.Bytes())
	tx.Id = hash[:]
}

func NewCoinbaseTX(to, data string) *Transaction {
	if data == "" {
		data = fmt.Sprintf("Sending Reward to '%s'", to)
	}

	txin := TXInput{[]byte{}, -1, data}
	txout := TXOutput{subsidy, to}
	tx := Transaction{nil, []TXInput{txin}, []TXOutput{txout}}

	tx.SetId()

	return &tx
}
