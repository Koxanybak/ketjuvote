package main

import (
	"crypto/sha256"
	"encoding/json"
	"log"
	"github.com/libp2p/go-libp2p"
	"context"
)

// Block represents a single block in the blockchain
type Block struct {
	ID			int			`json:"id"`
	Hash		[32]byte	`json:"-"`
	PrevHash	[32]byte	`json:"-"`
	Data		string		`json:"data"`
}

// BlockChain (kinda obvious)
var BlockChain []Block

// Calculate sha256 hash for the block
func (block *Block) calculateHash() [32]byte {
	jsonEnc, err := json.Marshal(*block)
	if err != nil {
		log.Fatal(err)
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

func main() {
	libp2p.New(context.Background(), )
}