package command

import (
	"flag"
	"fmt"
	"os"

	"github.com/aurora-is-near/evm-bully/replayer"
)

// Genesis implements the 'genesis' command.
func Genesis(argv0 string, args ...string) error {
	var f testnetFlags
	fs := flag.NewFlagSet(argv0, flag.ContinueOnError)
	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s\n", argv0)
		fmt.Fprintf(os.Stderr, "Process genesis block.\n")
		fs.PrintDefaults()
	}
	f.registerFlags(fs)
	if err := fs.Parse(args); err != nil {
		return err
	}
	_, testnet, err := f.determineTestnet()
	if err != nil {
		return err
	}
	if fs.NArg() != 0 {
		fs.Usage()
		return flag.ErrHelp
	}
	return replayer.ProcGenesisBlock(testnet)
}
