package block

import (
	"time"
)

type Block struct {
	Timestamp     int64
	Data          []byte
	Hash          []byte
	PrevBlockHash []byte
	Nounce        int
}

func NewBlock(data string, prevBlockHash []byte) *Block {
	b := Block{
		Data:          []byte(data),
		Timestamp:     time.Now().Unix(),
		PrevBlockHash: prevBlockHash,
	}
	b.SetHash()

	return &b
}

func (b *Block) SetHash() {
	pow := NewProofOfWork(b)
	nounce, hash := pow.Mine()

	b.Hash = hash[:]
	b.Nounce = nounce
}

func NewGenesisBlock() *Block {
	return NewBlock("Genesis Block", nil)
}
