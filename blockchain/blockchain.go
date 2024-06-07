package blockchain

import (
	"fmt"

	"github.com/dgraph-io/badger"
)

// moving everything related to BlockChain struct here
// bitcoin has two types of data:
// blocks stored with metadata & chain state object

// in Bitcoin each block has a seprate file on the disk, this is done to increase performance // we wont need to open multiple files just to read one block

const (
	dbPath = "./tmp/blocks" // path to our database
)

type BlockChain struct {
	LastHash []byte
	Database *badger.DB // pointer to badger database
}

type BlockChainIterator struct { // structure to iterate over the blocks
	CurrentHash []byte
	Database    *badger.DB
}

func InitBlockChain() *BlockChain {
	var lastHash []byte

	opts := badger.DefaultOptions // struct
	// pointing both of them towards same path
	opts.Dir = dbPath      // dir part of data base will store keys and metadata
	opts.ValueDir = dbPath // ValueDir part will store all values // this does not matter because both are stored in same place

	db, err := badger.Open(opts) // open the database, returns db and a error,  db -> pointer to data base
	Handle(err)

	err = db.Update(func(txn *badger.Txn) error { // update allows to read and write transaction on database //closure is pointer to a badger transaction
		//func(txn *badger.Txn) == closer// closer takes in  pointer to a badger transaction and passes back an error
		if _, err := txn.Get([]byte("lh")); err == badger.ErrKeyNotFound { // if err = true-> database does not exist -> blockchain does not exist // lh is key -> last hash; if err == badger.ErrKeyNotFound is true than this will work // []byte("lh") -> bytes representation of lh
			fmt.Println("No existing blockchain found")
			genesis := Genesis() // creating genesis block
			fmt.Println("Genesis proved")
			err = txn.Set(genesis.Hash, genesis.Serialize()) // using genesis hash as key for genesis block, serializing the genesis block into the database usnig txn.Set
			Handle(err)
			err = txn.Set([]byte("lh"), genesis.Hash) // genesis is the only block rn so we are setting its hash as the last hash in the database

			lastHash = genesis.Hash // putting hash in memory

			return err // handling error outside of the dbUpdate function
		} else { // if blockchain already exists
			item, err := txn.Get([]byte("lh")) // calling get on the key []byte("lh")
			Handle(err)
			lastHash, err = item.Value() // getting vlaue from item struct
			return err                   //   handling error outside of the dbUpdate function
		}
	})

	Handle(err) // handling db.Update func error

	blockchain := BlockChain{lastHash, db} // putting blockchain in memory // {lastHash, db}-> last hash and pointer to badger database
	return &blockchain // return pointer to blockchain
}

func (chain *BlockChain) AddBlock(data string) {
	var lastHash []byte

	err := chain.Database.View(func(txn *badger.Txn) error { // read only func to read the lastHash
		item, err := txn.Get([]byte("lh"))
		Handle(err)
		lastHash, err = item.Value()

		return err
	})
	Handle(err)

	newBlock := CreateBlock(data, lastHash) // creating new block with the last hash and the data

	err = chain.Database.Update(func(txn *badger.Txn) error { // read and write transaction on databse-> new block-> assign 
		err := txn.Set(newBlock.Hash, newBlock.Serialize()) // 
		Handle(err)
		err = txn.Set([]byte("lh"), newBlock.Hash)

		chain.LastHash = newBlock.Hash // updating BlockChain last hash to point to new block

		return err
	})
	Handle(err)
}

func (chain *BlockChain) Iterator() *BlockChainIterator {
	iter := &BlockChainIterator{chain.LastHash, chain.Database}

	return iter
}

func (iter *BlockChainIterator) Next() *Block { // calling next until it hits the end of data
	var block *Block

	err := iter.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get(iter.CurrentHash)
		Handle(err)
		encodedBlock, err := item.Value() // encodedBlock -> byte representation of a block
		block = Deserialize(encodedBlock)

		return err
	})
	Handle(err)

	iter.CurrentHash = block.PrevHash // using previous hash to find the next block, thus we will keep going backwards
	return block
}
