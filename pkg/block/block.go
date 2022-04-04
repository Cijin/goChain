package block

import (
	"bytes"
	"crypto/sha256"
	"time"

	transaction "github.com/Cijin/gochain/pkg/transaction"
)

type Block struct {
	Timestamp     int64
	Transactions  []*transaction.Transaction
	Hash          []byte
	PrevBlockHash []byte
	Nounce        int
}

func NewBlock(tx []*transaction.Transaction, prevBlockHash []byte) *Block {
	b := Block{
		Transactions:  tx,
		Timestamp:     time.Now().Unix(),
		PrevBlockHash: prevBlockHash,
	}
	b.SetHash()

	return &b
}

func (b *Block) SetHash() {
	pow := NewProofOfWork(b)
	nounce, hash := pow.Run()

	b.Hash = hash[:]
	b.Nounce = nounce
}

func (b *Block) HashTransactions() []byte {
	var txIds [][]byte
	var hash [32]byte

	for _, tx := range b.Transactions {
		txIds = append(txIds, tx.Id)
	}

	hash = sha256.Sum256(bytes.Join(txIds, []byte{}))
	return hash[:]
}

func NewGenesisBlock(coinbase *transaction.Transaction) *Block {
	return NewBlock([]*transaction.Transaction{coinbase}, []byte{})
}
