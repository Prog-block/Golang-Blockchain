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

//  A block is only marked as “valid” if the hash value of the entire block is below the difficulty hash. A block contains crucial transaction information that can’t be changed. So, the Miners change the nonce to get the hash lower than the difficulty threshold.
// Take the data from the block

// create a counter (nonce) which starts at 0

// create a hash of the data plus the counter

// check the hash to see if it meets a set of requirements

// Requirements:
// The First few bytes must contain 0s

const Difficulty = 18 // in realiy difficulty is dynamic controlled by a algo
// miners or transactions increase -> relults in more computation power ->  difficulty increase, as we want the block rate to be same and also the time required to mine a block to be same

type ProofOfWork struct {
	Block  *Block   // pointer to a block // block is block of the blockchain
	Target *big.Int // target which is a big.Int pointer // target represents the Requirements(above)
}

func NewProof(b *Block) *ProofOfWork { // takes pointer to a block and produces a pointer to a pow // new proof for a new block
	target := big.NewInt(1)                  // casting 1 as newBig int(binary) // 001
	target.Lsh(target, uint(256-Difficulty)) // left shifting the target by: 256-Difficulty // 001 * 2*(256-18)  = 0010000000...238 zeros....
	// valid block hash must be less than this target. Therefore more left shifting will generate a bigger target hash and thus block hash will become easy to calculate
	pow := &ProofOfWork{b, target} // putting block and left shifted target to the instance of ProofOfWork

	return pow
}

func (pow *ProofOfWork) InitData(nonce int) []byte { // it justs joins the data
	data := bytes.Join(
		[][]byte{
			pow.Block.PrevHash,
			pow.Block.Data,
			ToHex(int64(nonce)), // hex to convert int into bytes
			ToHex(int64(Difficulty)),
		},
		[]byte{},
	)

	return data
}

func (pow *ProofOfWork) Run() (int, []byte) { // Calculating hash of block, adjusting nonce to compare to target
	var intHash big.Int
	var hash [32]byte // array of size 32 bytes

	nonce := 0

	for nonce < math.MaxInt64 { // infinite loop
		data := pow.InitData(nonce)
		hash = sha256.Sum256(data)

		fmt.Printf("\r%x", hash)
		intHash.SetBytes(hash[:]) // convert hash to bigInt // compare this big int with our target bigInt

		if intHash.Cmp(pow.Target) == -1 { // compare pow target bigInt with intHash bigInt, -1 would mean our hash is actually less than the target we are looking for
			break
		} else {
			nonce++
		}
	}
	fmt.Println()

	return nonce, hash[:]
}

func (pow *ProofOfWork) Validate() bool { // validating if the nonce is valid
	var intHash big.Int // bigInt version of hash

	data := pow.InitData(pow.Block.Nonce)

	hash := sha256.Sum256(data)
	intHash.SetBytes(hash[:]) // convert hash to bigInt

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
