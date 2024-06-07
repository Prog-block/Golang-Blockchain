package main

import (
	"blockchain/blockchain"
	"flag"
	"fmt"
	"os"
	"runtime"
	"strconv"
)

// persistance layer, key value storage database
// we will use badger DB, it is native golang data base
// Badger db only accepts arrays or series of bytes
type CommandLine struct {
	blockchain *blockchain.BlockChain
}

func (cli *CommandLine) printUsage() {
	fmt.Println("Usage:")
	fmt.Println(" add -block BLOCK_DATA - add a block to the chain")
	fmt.Println(" print - Prints the blocks in the chain")
}

func (cli *CommandLine) validateArgs() { // validate arguments that are passed through the command line
	if len(os.Args) < 2 { // length of operating system arguments
		cli.printUsage() // printing usage because user has not enterd a command
		runtime.Goexit() // using this because unlike os.exit it exits the application by shutting down the go-routine
		// downsides with badger database-> it needs to collect garbage collect values and keys before it shuts down
		// if appplication shuts down without properly closing the database it can corrupt the data
		// thus using Goexit() to initiate a shutdowwn
	}
}

func (cli *CommandLine) addBlock(data string) {
	cli.blockchain.AddBlock(data)
	fmt.Println("Added Block!")
}

func (cli *CommandLine) printChain() {
	iter := cli.blockchain.Iterator() // converts blockchain struct into an iterator struct

	for {
		block := iter.Next() //

		fmt.Printf("Prev. hash: %x\n", block.PrevHash)
		fmt.Printf("Data: %s\n", block.Data)
		fmt.Printf("Hash: %x\n", block.Hash)
		pow := blockchain.NewProof(block)
		fmt.Printf("PoW: %s\n", strconv.FormatBool(pow.Validate()))
		fmt.Println()

		if len(block.PrevHash) == 0 { // genesis does not have a previous hash thus length of its prevHash will be zero
			break
		}
	}
}

func (cli *CommandLine) run() {
	cli.validateArgs()

	addBlockCmd := flag.NewFlagSet("add", flag.ExitOnError)       // create new flag set if user types in add
	printChainCmd := flag.NewFlagSet("print", flag.ExitOnError)   // create new flag set if user types in print
	addBlockData := addBlockCmd.String("block", "", "Block data") // subset to addBlockCmd // if user types add and then block

	switch os.Args[1] { // calling it on the 1st argument after the original call to the program
	case "add": // add block to blockchain
		err := addBlockCmd.Parse(os.Args[2:]) // parsing arguments which come after 1st argument
		blockchain.Handle(err)

	case "print":
		err := printChainCmd.Parse(os.Args[2:])
		blockchain.Handle(err)

	default: // type anything else or nothing
		cli.printUsage()
		runtime.Goexit() // shutting down the program
	}

	if addBlockCmd.Parsed() { // if addBlockCmd is parsed than do this //it will be parsed or not-> addBlockCmd.Parse returns a boolean
		if *addBlockData == "" { // if pointer to addBlockData is an empty string
			addBlockCmd.Usage() //print out addBlockCmd.Usage()
			runtime.Goexit()    //
		}
		cli.addBlock(*addBlockData) // if not empty make a neww block in the blockchain
	}

	if printChainCmd.Parsed() {
		cli.printChain()
	}
}

func main() {
	defer os.Exit(0) // safety to close the database
	chain := blockchain.InitBlockChain()
	defer chain.Database.Close() // closing the database before the main function ends
	// defer only executes if goExit functions properly
	cli := CommandLine{chain}
	cli.run()
}

//  go run main.go print
//  go run main.go add -block "first block"
//  go run main.go print
//  go run main.go add -block "creating 2nd block"
