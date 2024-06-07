package blockchain

import (
	"bytes"
	"encoding/gob"
	"log"
)
// bitcoin and other currencies use levelDB as there database, it is a key value based database, we will use badger DB (native database for golang) it is also key value based like levelDB // byte keys byte values, putting them into files
// badger db only accepts array/slices of bytes
// therefore we need function to cerealize and decerealize the Block struct to bytes // we will make it Block.go
type Block struct {
	Hash     []byte
	Data     []byte
	PrevHash []byte
	Nonce    int
}

func CreateBlock(data string, prevHash []byte) *Block {
	block := &Block{[]byte{}, []byte(data), prevHash, 0}
	pow := NewProof(block)
	nonce, hash := pow.Run()

	block.Hash = hash[:]
	block.Nonce = nonce

	return block
}

func Genesis() *Block {
	return CreateBlock("Genesis", []byte{})
}

func (b *Block) Serialize() []byte { // method on the BlockStruct, outputs a slice of bytes(bytes representation of our block)
	var res bytes.Buffer
	encoder := gob.NewEncoder(&res) // encoder on our results bytes buffer

	err := encoder.Encode(b) // calling encode on the block itself

	Handle(err)

	return res.Bytes() // returning bytes portion of results
}

func Deserialize(data []byte) *Block {
	var block Block

	decoder := gob.NewDecoder(bytes.NewReader(data))
	// bytes.NewReader(data) creates bytes reader which gets passed to NewDecoder
	err := decoder.Decode(&block)

	Handle(err)

	return &block
}

func Handle(err error) {
	if err != nil {
		log.Panic(err)
	}
}
