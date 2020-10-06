package main

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/binary"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/network"
)

// VoteData represents a single vote
type VoteData struct {
	EventID		int			`json:"eventID"`
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
func calculateHash(ID int, data *VoteData) [32]byte {
	byteData, err := json.Marshal(*data)
	if err != nil {
		panic(err)
	}
	buf := new(bytes.Buffer)
	err = binary.Write(buf, binary.LittleEndian, ID)
	if err != nil {
		panic(err)
	}
	byteData = append(byteData, buf.Bytes()...)
	return sha256.Sum256(byteData)
}

// Checks if the block is valid
func (block *Block) isValid(prev *Block) bool {
	if prev.Hash != block.PrevHash {
		return false
	}
	if block.Hash != calculateHash(block.ID, &block.Data) {
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
	hsh := calculateHash(ID, &data)
	return Block{
		ID,
		data,
		hsh,
		hsh,
	}
}

// Creates a p2p node
func createNode(listenPort string) (host.Host, error) {
	priv, _, err := crypto.GenerateRSAKeyPair(2048, *new(io.Reader))
	if err != nil {
		panic(err)
	}
	fmt.Println(priv)

	opts := []libp2p.Option{
		libp2p.ListenAddrStrings(fmt.Sprintf("/ip4/127.0.0.1/tcp/%s", listenPort)),
		libp2p.Identity(priv),
	}
	node, err := libp2p.New(context.Background(), opts...)
	if err != nil {
		return nil, err
	}

	return node, nil
}

// Stream handler
func handler(s network.Stream) {
	log.Println("Haloo")
}

func main() {
	genesisBlock := generateGenesisBlock()
	BlockChain = append(BlockChain, genesisBlock)

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Unable to load .env file: %v", err)
	}
	listenPort := flag.String("p", os.Getenv("LISTEN_PORT"), "port to listen to")
	//targetAddr := flag.String("t", "", "target node to connect to")

	node, err := createNode(*listenPort)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Listening on multiaddr: %v\n", node.Addrs())
	node.SetStreamHandler("/p2p/1.0.0", handler)

	select {}
}