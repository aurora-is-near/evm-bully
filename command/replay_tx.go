package command

import (
	"flag"
	"fmt"
	"os"

	"github.com/aurora-is-near/evm-bully/replayer"
	"github.com/aurora-is-near/near-api-go"
)

// ReplayTx implements the 'replay-tx' command.
func ReplayTx(argv0 string, args ...string) error {
	fs := flag.NewFlagSet(argv0, flag.ContinueOnError)
	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s <breakpointDir>\n", argv0)
		fmt.Fprintf(os.Stderr, "Replay transaction from breakpoint directory.\n")
		fs.PrintDefaults()
	}
	build := fs.Bool("build", false, "Build nearcore and aurora-engine before replaying tx")
	contract := fs.String("contract", "", "Upgrade EVM contract with file before replaying tx")
	gas := fs.Uint64("gas", defaultGas, "Max amount of gas a call can use (in gas units)")
	release := fs.Bool("release", false, "Run release version of neard")
	cfg := near.GetConfig()
	registerCfgFlags(fs, cfg, true)
	if err := fs.Parse(args); err != nil {
		return err
	}
	if fs.NArg() != 1 {
		fs.Usage()
		return flag.ErrHelp
	}
	breakpointDir := fs.Arg(0)
	return replayer.ReplayTx(breakpointDir, *build, *contract, *release, *gas)
}
