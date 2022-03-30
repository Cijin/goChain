package block

import (
	"bytes"
	"crypto/sha256"
	"time"
)

type Block struct {
	Timestamp     int64
	Data          []byte
	Hash          []byte
	PrevBlockHash []byte
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
	concatenatedData := bytes.Join([][]byte{b.Data, b.PrevBlockHash, b.Data}, []byte{})
	hash := sha256.Sum256(concatenatedData)

	b.Hash = hash[:]
}

func NewGenesisBlock() *Block {
	return NewBlock("Genesis Block", nil)
}
