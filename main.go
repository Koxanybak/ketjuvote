package main

import (
	"crypto/sha256"
	"encoding/json"
	"log"
	"github.com/libp2p/go-libp2p"
	"context"
)

type VoteData struct {
	eventID		int			`json:"eventID"`
}

// Block represents a single block in the blockchain
type Block struct {
	ID			int			`json:"id"`
	Data		VoteData	`json:"data"`
	Hash		[32]byte	`json:"-"`
	PrevHash	[32]byte	`json:"-"`
}

// BlockChain (kinda obvious)
var BlockChain []Block

// Calculate sha256 hash for the block
func (block *Block) calculateHash() [32]byte {
	jsonEnc, err := json.Marshal(*block)
	if err != nil {
		panic(err)
	}
	return sha256.Sum256(jsonEnc)
}

// Checks if the block is valid
func (block *Block) isValid(prev *Block) bool {
	if prev.Hash != block.PrevHash {
		return false
	}
	if block.Hash != block.calculateHash() {
		return false
	}
	return true
}

// Checks if the chain is valid
func isChainValid(chain []Block) bool {
	for i := 1; i < len(chain); i++ {
		if !chain[i].isValid(&chain[i-1]) {
			return false
		}
	}
	return true
}

// Generates the first block
func generateGenesisBlock() Block {
	ID := 0
	data := VoteData{ 0 }
	return Block{
		ID,
		data,
		
	}
}

func main() {
	node, err := libp2p.New(context.Background(), )
	if err !=  nil {
		panic(err)
	}
}