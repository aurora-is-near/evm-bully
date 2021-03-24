package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/aurora-is-near/evm-bully/command"
	"github.com/ethereum/go-ethereum/log"
)

func init() {
	log.Root().SetHandler(log.DiscardHandler())
}

func fatal(err error) {
	fmt.Fprintf(os.Stderr, "%s: error: %s\n", os.Args[0], err)
	os.Exit(1)
}

func usage() {
	cmd := os.Args[0] + " [-v]"
	fmt.Fprintf(os.Stderr, "Usage: %s genesis\n", cmd)
	fmt.Fprintf(os.Stderr, "       %s replay\n", cmd)
	fmt.Fprintf(os.Stderr, "       %s block\n", cmd)
	fmt.Fprintf(os.Stderr, "       %s state <accountId>\n", cmd)
	fmt.Fprintf(os.Stderr, "       %s call <contractName> <methodName>\n", cmd)
	fmt.Fprintf(os.Stderr, "       %s send <sender> <receiver> <amount>\n", cmd)
	fmt.Fprintf(os.Stderr, "Stress test and benchmark the NEAR EVM.\n")
	fmt.Fprintf(os.Stderr, "Global flags:\n")
	flag.PrintDefaults()
}

func main() {
	// define flags
	verbose := flag.Bool("v", false, "Be verbose")

	// parse flags
	flag.Usage = usage
	flag.Parse()

	// enable logging, if necessary
	if *verbose {
		log.Root().SetHandler(log.StdoutHandler)
	}

	// makes sure a command was given
	if flag.NArg() == 0 {
		usage()
		os.Exit(2)
	}

	// call command
	var err error
	argv0 := os.Args[0] + " " + flag.Args()[0]
	args := flag.Args()[1:]
	switch flag.Arg(0) {
	case "genesis":
		err = command.Genesis(argv0, args...)
	case "replay":
		err = command.Replay(argv0, args...)
	case "block":
		err = command.Block(argv0, args...)
	case "state":
		err = command.State(argv0, args...)
	case "call":
		err = command.Call(argv0, args...)
	case "send":
		err = command.Send(argv0, args...)
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
