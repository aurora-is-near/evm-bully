package command

import (
	"flag"
	"fmt"
	"os"

	"github.com/aurora-is-near/evm-bully/replayer"
)

// Stats implements the 'stats' command.
func Stats(argv0 string, args ...string) error {
	var f testnetFlags
	fs := flag.NewFlagSet(argv0, flag.ContinueOnError)
	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s\n", argv0)
		fmt.Fprintf(os.Stderr, "Calculate testnet statistics.\n")
		fs.PrintDefaults()
	}
	block := fs.Uint64("block", defaultGoerliBlockHeight, "Block height")
	dataDir := fs.String("datadir", defaultDataDir, "Data directory containing the database to read")
	defrost := fs.Bool("defrost", false, "Defrost the database first")
	dump := fs.Bool("dump", false, "Use dump file instead of database")
	hash := fs.String("hash", defaultGoerliBlockHash, "Block hash")
	f.registerFlags(fs)
	if err := fs.Parse(args); err != nil {
		return err
	}
	_, testnet, err := f.determineTestnet()
	if err != nil {
		return err
	}
	adjustBlockDefaults(block, hash, testnet)
	if fs.NArg() != 0 {
		fs.Usage()
		return flag.ErrHelp
	}
	// calculate statistics
	return replayer.CalcStats(*dataDir, testnet, *block, *hash, *defrost, *dump)
}
