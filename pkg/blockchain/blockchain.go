package blockchain

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/Cijin/gochain/pkg/block"
	"github.com/Cijin/gochain/pkg/transaction"
	"github.com/boltdb/bolt"
)

const BlocksBucket = "blocksBucket"
const blockchainDb = "blockchainDb"
const LeafKey = "l"
const genesisCoinbaseData = "It's fucking 2100, I need to sleep :D"

type Blockchain struct {
	Tip []byte
	Db  *bolt.DB
}

func (bc *Blockchain) MineBlock(tx []*transaction.Transaction) {
	/*
	 * get the tip of blockchain
	 * mine new block
	 * update tip ("l") and add new block to blockchain
	 */
	newBlock := block.NewBlock(tx, bc.Tip)

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

func isBlockchainDbPresent() bool {
	_, err := os.Stat(blockchainDb)

	return !os.IsNotExist(err)
}

/*
* @TODO: check if address required as a param
 */
func NewBlockchain() *Blockchain {
	var bl block.Block
	var tip []byte

	if !isBlockchainDbPresent() {
		fmt.Println("Blockchain does not exist, create one first")
		os.Exit(1)
	}

	db, err := bolt.Open(blockchainDb, 0600, nil)
	if err != nil {
		log.Panic(err)
	}

	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BlocksBucket))

		buf := b.Get([]byte(LeafKey))
		err = json.Unmarshal(buf, &bl)
		if err != nil {
			return err
		}

		tip = bl.Hash
		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	bc := Blockchain{tip, db}
	return &bc
}

func CreateBlockchain(address string) *Blockchain {
	var tip []byte

	if isBlockchainDbPresent() {
		fmt.Println("Blockchain already exists")
		os.Exit(1)
	}

	db, err := bolt.Open(blockchainDb, 0600, nil)
	if err != nil {
		log.Panic(err)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BlocksBucket))

		// Check if blockchain present in database
		if b == nil {
			cbtx := transaction.NewCoinbaseTX(address, genesisCoinbaseData)
			block := block.NewGenesisBlock(cbtx)

			// marshal JSON and write to bucket
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

		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	return &Blockchain{tip, db}
}
