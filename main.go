package main

import (
	"bufio"
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"log"
	mrand "math/rand"
	"os"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/protocol"
	"github.com/multiformats/go-multiaddr"
)

// VoteData represents a single vote
type VoteData struct {
	EventID int `json:"eventID"`
}

// Block represents a single block in the blockchain
type Block struct {
	ID       uint32   `json:"id"`
	Data     VoteData `json:"data"`
	Hash     [32]byte `json:"-"`
	PrevHash [32]byte `json:"-"`
}

// BlockChain (kinda obvious)
var BlockChain []Block

// Calculate sha256 hash for the block
func calculateHash(ID uint32, data *VoteData) [32]byte {
	byteData, err := json.Marshal(*data)
	if err != nil {
		panic(err)
	}
	buf := make([]byte, 32)
	binary.LittleEndian.PutUint32(buf, ID)
	byteData = append(byteData, buf...)
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
	var ID uint32 = 0
	data := VoteData{0}
	hsh := calculateHash(ID, &data)
	return Block{
		ID,
		data,
		hsh,
		hsh,
	}
}

// Creates a p2p node
func createNode(listenPort string, seed int64) (host.Host, error) {
	var ir io.Reader
	if seed == 0 {
		ir = rand.Reader
	} else {
		ir = mrand.New(mrand.NewSource(seed))
	}
	priv, _, err := crypto.GenerateRSAKeyPair(2048, ir)
	if err != nil {
		panic(err)
	}

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

func readData(rw *bufio.ReadWriter, comms chan int) {
	defer close(comms)
	for {
		msg, err := rw.ReadString('\n')
		if err != nil {
			fmt.Printf("\nUnable to read from the data stream: %v\n", err)
			// Terminate write go routine
			comms <- 1
			break
		} else {
			fmt.Printf("\nSait viestin: %s\nKirjoita viesti > ", msg)
		}
	}
}

func writeData(rw *bufio.ReadWriter, comms chan int) {
	stdReader := bufio.NewReader(os.Stdin)
	for {
		select {
		default:
			fmt.Print("Kirjoita viesti > ")
			dataToSend, err := stdReader.ReadString('\n')
			if err != nil {
				fmt.Println("\nError reading from stdin")
				panic(err)
			}
			_, err = rw.WriteString(dataToSend)
			if err != nil {
				fmt.Printf("\nUnable to write to the buffer: %v\n", err)
			}
			err = rw.Flush()
			if err != nil {
				fmt.Printf("\nUnable to flush the write buffer: %v\n", err)
			}

		case <-comms:
			break
		}
	}
}

// Stream handler
func handler(s network.Stream) {
	fmt.Println("Got a stream I guess")

	// Channel that lets read and write functions communicate when connections stops
	comms := make(chan int, 1)

	rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))
	go readData(rw, comms)
	go writeData(rw, comms)
}

// Main func
func main() {
	genesisBlock := generateGenesisBlock()
	BlockChain = append(BlockChain, genesisBlock)

	conf := getConfig()

	node, err := createNode(conf.ListenPort, conf.Seed)
	if err != nil {
		log.Fatal(err)
	}
	addrs, _ := peer.AddrInfoToP2pAddrs(&peer.AddrInfo{ID: node.ID(), Addrs: node.Addrs()})
	fmt.Printf("Listening on multiaddr: %v\n", addrs)
	node.SetStreamHandler("/p2p/1.0.0", handler)

	if conf.TargetAddr != "" {
		addr, _ := multiaddr.NewMultiaddr(conf.TargetAddr)
		targetInfo, _ := peer.AddrInfoFromP2pAddr(addr)
		node.Connect(context.Background(), *targetInfo)

		stream, err := node.NewStream(context.Background(), targetInfo.ID, protocol.ID("/p2p/1.0.0"))
		if err != nil {
			fmt.Printf("Unable to connect")
		} else {
			handler(stream)
		}
	}

	select {}
}
