package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
)

type BlockChain struct { // it contains blocks, in future hash will be used to point ot blocks
	blocks []*Block
}

type Block struct {
	Hash     []byte // calc by hashing of prevHash and data
	Data     []byte
	PrevHash []byte
}

func (b *Block) DeriveHash() { // method to create hash from previous hash and data
	info := bytes.Join([][]byte{b.Data, b.PrevHash}, []byte{}) // 2d slice of bytes
	// after comma is empty bytes
	hash := sha256.Sum256(info) // sha256 is simpler version of actual hasshing mechanism of blockchain
	b.Hash = hash[:]
}

func CreateBlock(data string, prevHash []byte) *Block {
	block := &Block{[]byte{}, []byte(data), prevHash} // taking empty data for hash feild
	// taking data string and converting it into a slice of bytes
	block.DeriveHash()
	return block // returns hash
}

func (chain *BlockChain) AddBlock(data string) { // method to add a block to a chain, string is the data to put into the block
	// chain *BlockChain -> pointer to blockchain
	prevBlock := chain.blocks[len(chain.blocks)-1] // length of blocks -1 block
	new := CreateBlock(data, prevBlock.Hash)
	chain.blocks = append(chain.blocks, new) // adding this new block
}

func Genesis() *Block { // func to create 1st block
	return CreateBlock("Genesis", []byte{})
}

func InitBlockChain() *BlockChain { // it will build initial blockchain using genesis block
	return &BlockChain{[]*Block{Genesis()}}
}

func main() {
	chain := InitBlockChain()

	chain.AddBlock("First Block after Genesis")
	chain.AddBlock("Second Block after Genesis")
	chain.AddBlock("Third Block after Genesis")

	for _, block := range chain.blocks {// for loop
		// fmt.Printf("Previous Hash: %x\n", block.PrevHash)
		fmt.Printf("Data in Block: %s\n", block.Data) // %s -> view data in string
		fmt.Printf("Hash: %x\n", block.Hash)
	}
}
