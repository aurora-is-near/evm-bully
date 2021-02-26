package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/node"
	"github.com/near/evm-bully/replayer"
)

const (
	defaultEndpoint    = "http://localhost:8545"
	defaultBlockHeight = 100000
	defaultBlockhash   = "0x94b9c7be22783a3ee1e1f1f31e35f7de4612905d6fd24d3fe90a26091dca43fe"
)

var (
	defaultDataDir = node.DefaultDataDir()
)

func determineTestnet(goerli, rinkeby, ropsten bool) (string, error) {
	if !goerli && !rinkeby && !ropsten {
		return "", errors.New("one of the options -goerli, -rinkeby, or -ropsten is mandatory")
	}
	if goerli && rinkeby {
		return "", errors.New("the options -goerli and -rinkeby exclude each other")
	}
	if goerli && ropsten {
		return "", errors.New("the options -goerli and -ropsten exclude each other")
	}
	if rinkeby && ropsten {
		return "", errors.New("the options -rinkeby and -ropsten exclude each other")
	}
	if rinkeby {
		return "rinkeby", nil
	} else if ropsten {
		return "ropsten", nil
	}
	// use Görli as the default
	return "goerli", nil
}

func fatal(err error) {
	fmt.Fprintf(os.Stderr, "%s: error: %s\n", os.Args[0], err)
	os.Exit(1)
}

func main() {
	// define flags
	block := flag.Uint64("block", defaultBlockHeight, "Block height")
	datadir := flag.String("datadir", defaultDataDir, "Data directory containing the database to read")
	endpoint := flag.String("endpoint", defaultEndpoint, "Set default JSON-RPC endpoint")
	genesis := flag.Bool("genesis", false, "Process genesis block")
	goerli := flag.Bool("goerli", false, "Use the Görli testnet")
	hash := flag.String("hash", defaultBlockhash, "Block hash")
	rinkeby := flag.Bool("rinkeby", false, "Use the Rinkeby testnet")
	ropsten := flag.Bool("ropsten", false, "Use the Ropsten testnet")
	verbose := flag.Bool("v", false, "Be verbose")

	// parse flags
	flag.Parse()

	// enable logging, if necessary
	if *verbose {
		log.Root().SetHandler(log.StdoutHandler)
	}

	// determine testnet name from flags
	testnet, err := determineTestnet(*goerli, *rinkeby, *ropsten)
	if err != nil {
		fatal(err)
	}

	if *genesis {
		if err := replayer.ProcGenesisBlock(testnet); err != nil {
			fatal(err)
		}
	} else {
		// run replayer
		err := replayer.ReadTxs(context.Background(), *endpoint, *datadir, testnet, *block, *hash)
		if err != nil {
			fatal(err)
		}
	}
}
