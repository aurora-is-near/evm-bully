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
	block := fs.Uint64("block", defaultBlockHeight, "Block height")
	dataDir := fs.String("datadir", defaultDataDir, "Data directory containing the database to read")
	defrost := fs.Bool("defrost", false, "Defrost the database first")
	hash := fs.String("hash", defaultBlockhash, "Block hash")
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
	// determine cache directory
	cacheDir, err := determineCacheDir(testnet)
	if err != nil {
		return err
	}
	// calculate statistics
	return replayer.CalcStats(*dataDir, testnet, cacheDir, *block, *hash, *defrost)
}
