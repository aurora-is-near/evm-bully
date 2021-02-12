package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/node"
	"github.com/near/evm-bully/replayer"
)

const (
	defaultBlockHeight = 100000
	defaultBlockhash   = "0x94b9c7be22783a3ee1e1f1f31e35f7de4612905d6fd24d3fe90a26091dca43fe"
)

var (
	defaultDataDir = filepath.Join(node.DefaultDataDir(), "goerli", "geth", "chaindata")
)

func fatal(err error) {
	fmt.Fprintf(os.Stderr, "%s: error: %s\n", os.Args[0], err)
	os.Exit(1)
}

func main() {
	block := flag.Uint64("block", defaultBlockHeight, "Block height")
	datadir := flag.String("datadir", defaultDataDir, "Data directory containing the database to read")
	hash := flag.String("hash", defaultBlockhash, "Block hash")
	verbose := flag.Bool("v", false, "Be verbose")

	flag.Parse()

	if *verbose {
		log.Root().SetHandler(log.StdoutHandler)
	}
	if err := replayer.ReadTxs(*datadir, *block, *hash); err != nil {
		fatal(err)
	}
}
