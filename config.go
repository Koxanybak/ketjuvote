package main

import (
	"flag"
	"github.com/joho/godotenv"
	"os"
	"log"
)

type config struct {
	Seed		int64
	ListenPort	string
	TargetAddr 	string
}

func getConfig() *config {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Unable to load .env file: %v", err)
	}
	listenPort := flag.String("p", os.Getenv("LISTEN_PORT"), "port to listen to")
	targetAddr := flag.String("t", "", "target node's multiaddress")
	seed := flag.Int64("s", 0, "seed for identity generation")
	flag.Parse()

	return &config{ *seed, *listenPort, *targetAddr }
}