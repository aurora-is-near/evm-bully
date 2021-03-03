package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/ethereum/go-ethereum/log"
	"github.com/near/evm-bully/command"
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

func usage() {
	cmd := os.Args[0] + " [flags]"
	fmt.Fprintf(os.Stderr, "Usage: %s genesis\n", cmd)
	fmt.Fprintf(os.Stderr, "       %s replay\n", cmd)
	fmt.Fprintf(os.Stderr, "       %s blocknumber\n", cmd)
	fmt.Fprintf(os.Stderr, "Stress test and benchmark the NEAR EVM.\n")
	fmt.Fprintf(os.Stderr, "Flags:\n")
	flag.PrintDefaults()
}

func main() {
	// define flags
	goerli := flag.Bool("goerli", false, "Use the Görli testnet")
	rinkeby := flag.Bool("rinkeby", false, "Use the Rinkeby testnet")
	ropsten := flag.Bool("ropsten", false, "Use the Ropsten testnet")
	verbose := flag.Bool("v", false, "Be verbose")

	// parse flags
	flag.Usage = usage
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

	// makes sure a command was given
	if flag.NArg() == 0 {
		usage()
		os.Exit(2)
	}

	// call command
	argv0 := os.Args[0] + " " + flag.Args()[0]
	args := flag.Args()[1:]
	switch flag.Arg(0) {
	case "genesis":
		err = command.Genesis(testnet, argv0, args...)
	case "replay":
		err = command.Replay(testnet, argv0, args...)
	case "blocknumber":
		err = command.BlockNumber(argv0, args...)
	default:
		usage()
	}
	if err != nil {
		if err != flag.ErrHelp {
			fatal(err)
		}
		os.Exit(2)
	}
}
