package blockchain

import (
	"encoding/json"
	"log"

	"github.com/Cijin/gochain/pkg/block"
	"github.com/boltdb/bolt"
)

type BlockchainIterator struct {
	CurrentHash []byte
	Db          *bolt.DB
}

/*
 * Initially currentHash will be the tip of the blockchain
 */
func (bc *Blockchain) Iterator() *BlockchainIterator {
	return &BlockchainIterator{
		CurrentHash: bc.Tip,
		Db:          bc.Db,
	}
}

/*
 * Get the hash of the Prev block
 */
func (bcI *BlockchainIterator) Previous() *block.Block {
	var block block.Block
	err := bcI.Db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BlocksBucket))

		buf := b.Get(bcI.CurrentHash)
		err := json.Unmarshal(buf, &block)
		if err != nil {
			return err
		}

		bcI.CurrentHash = block.Hash
		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	return &block
}
