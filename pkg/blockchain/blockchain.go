package blockchain

import (
	"encoding/json"
	"log"

	"github.com/Cijin/gochain/pkg/block"
	"github.com/boltdb/bolt"
)

const BlocksBucket = "blocksBucket"
const blockchainDb = "blockchainDb"
const LeafKey = "l"

type Blockchain struct {
	Tip []byte
	Db  *bolt.DB
}

func (bc *Blockchain) AddBlock(data string) {
	/*
	 * get the tip of blockchain
	 * mine new block
	 * update tip ("l") and add new block to blockchain
	 */
	newBlock := block.NewBlock(data, bc.Tip)

	err := bc.Db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BlocksBucket))
		buf, err := json.Marshal(newBlock)
		if err != nil {
			return err
		}

		err = b.Put(newBlock.Hash, buf)
		if err != nil {
			return err
		}

		err = b.Put([]byte(LeafKey), newBlock.Hash)
		if err != nil {
			return err
		}

		// update blockchain
		bc.Tip = newBlock.Hash
		return nil
	})

	if err != nil {
		log.Panic(err)
	}
}

func NewBlockchain() *Blockchain {
	var tip []byte
	var bl block.Block

	db, err := bolt.Open(blockchainDb, 0600, nil)
	if err != nil {
		log.Panic(err)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BlocksBucket))

		// Check if blockchain present in database
		if b == nil {
			// marshal JSON and write to bucket
			block := block.NewGenesisBlock()
			buf, err := json.Marshal(block)
			if err != nil {
				return err
			}

			// create bucket
			b, err := tx.CreateBucket([]byte(BlocksBucket))
			if err != nil {
				return err
			}

			err = b.Put([]byte(LeafKey), buf)
			if err != nil {
				return err
			}

			tip = block.Hash
			return b.Put(block.Hash, buf)
		}

		// blockchain present
		buf := b.Get([]byte(LeafKey))
		err = json.Unmarshal(buf, &bl)
		if err != nil {
			log.Panic(err)
		}

		tip = bl.Hash
		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	return &Blockchain{tip, db}
}
