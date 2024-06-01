// proof.go and block.go are both part of the "package blockchain"
package blockchain // can use functionality of block.go

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"log"
	"math"
	"math/big"
)

// Take the data from the block

// create a counter (nonce) which starts at 0

// create a hash of the data plus the counter

// check the hash to see if it meets a set of requirements

// Requirements:
// The First few bytes must contain 0s

const Difficulty = 18 // in realiy difficulty is dynamic controlled by a algo
// miners or transactions increase -> relults in more computation power ->  difficulty increase, as we want the block rate to be same and also the time required to mine a block to be same

type ProofOfWork struct {
	Block  *Block   // pointer to a block
	Target *big.Int // target which is a big.Int pointer // target represents the Requirements(above)
}

func NewProof(b *Block) *ProofOfWork { // takes pointer to a block and produces a pointer to a pow
	target := big.NewInt(1)                  // casting 1 as newBig int
	target.Lsh(target, uint(256-Difficulty)) // left shifting the target by: 256-Difficulty

	pow := &ProofOfWork{b, target} // putting block and left shifted target to the instance of ProofOfWork

	return pow
}

func (pow *ProofOfWork) InitData(nonce int) []byte { // output slash of bytes
	data := bytes.Join(
		[][]byte{
			pow.Block.PrevHash,
			pow.Block.Data,
			ToHex(int64(nonce)),
			ToHex(int64(Difficulty)),
		},
		[]byte{},
	)

	return data
}

func (pow *ProofOfWork) Run() (int, []byte) {
	var intHash big.Int
	var hash [32]byte

	nonce := 0

	for nonce < math.MaxInt64 {
		data := pow.InitData(nonce)
		hash = sha256.Sum256(data) //

		fmt.Printf("\r%x", hash)
		intHash.SetBytes(hash[:]) // convert hast to bigInt

		if intHash.Cmp(pow.Target) == -1 { // compare pow target with the new bigInt version of hash, -1 would mean our hash is actually less than the target we are looking for
			break
		} else {
			nonce++
		}

	}
	fmt.Println()

	return nonce, hash[:]
}

func (pow *ProofOfWork) Validate() bool {
	var intHash big.Int

	data := pow.InitData(pow.Block.Nonce)

	hash := sha256.Sum256(data)
	intHash.SetBytes(hash[:])

	return intHash.Cmp(pow.Target) == -1
}

func ToHex(num int64) []byte {
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, num) // decode num into bytes
	//binary.BigEndian signifies how we want the bytes to be organised
	if err != nil {
		log.Panic(err)

	}

	return buff.Bytes() // returning bytes portion of buffer
}
