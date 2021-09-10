package command

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/aurora-is-near/evm-bully/replayer"
	"github.com/aurora-is-near/near-api-go"
)

// Replay implements the 'replay' command.
func Replay(argv0 string, args ...string) error {
	var testnetFlags testnetFlags
	fs := flag.NewFlagSet(argv0, flag.ContinueOnError)
	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [<evmContract>]\n", argv0)
		fmt.Fprintf(os.Stderr, "Replay transactions to NEAR EVM installed in account <evmContract>.\n")
		fs.PrintDefaults()
	}
	autobreak := fs.Bool("autobreak", false, "Automatically repeat with a break point after an error")
	accountID := fs.String("accountId", "", "Unique identifier for the account that will be used to sign this call")
	batch := fs.Bool("batch", false, "Batch transactions")
	batchSize := fs.Int("size", 10, "Batch size when batching transactions")
	breakBlock := fs.Int("breakblock", -1, "Break replaying at this block height")
	breakTx := fs.Int("breaktx", 0, "Break replaying at this transaction (in block given by -breakblock)")
	contract := fs.String("contract", "", "EVM contract file to deploy")
	dataDir := fs.String("datadir", defaultDataDir, "Data directory containing the database to read")
	defrost := fs.Bool("defrost", false, "Defrost the database first")
	gas := fs.Uint64("gas", defaultGas, "Max amount of gas a call can use (in gas units)")
	initialBalance := fs.String("initial-balance", defaultInitialBalance, "Number of tokens to transfer to newly created account")
	release := fs.Bool("release", false, "Run release version of neard (instead of debug version)")
	setup := fs.Bool("setup", false, "Setup and run neard before replaying (auto-deploys contract)")
	neardPath := fs.String("neard", "", "Path to neard binary (won't build neard if -setup is provided)")
	neardHead := fs.String("neardhead", "", "Git hash of neard (required if -neard is provided)")
	skip := fs.Bool("skip", false, "Skip empty blocks during replay")
	startBlock := fs.Int("startblock", 0, "Start replaying at this block height")
	startTx := fs.Int("starttx", 0, "Start replaying at this transaction (in block given by -startblock)")
	timeout := fs.Duration("timeout", 0, "Timeout for JSON-RPC client")
	cfg := near.GetConfig()
	registerCfgFlags(fs, cfg, true)
	testnetFlags.registerFlags(fs)
	if err := fs.Parse(args); err != nil {
		return err
	}
	if !*setup && *initialBalance != defaultInitialBalance {
		return errors.New("option -initial-balance requires -setup")
	}
	if !*setup && *accountID == "" {
		return errors.New("option -accountId is mandatory")
	}
	if *autobreak && *breakBlock != -1 {
		return errors.New("options -autobreak and -breakblock exclude each other")
	}
	if *autobreak && *breakTx != 0 {
		return errors.New("options -autobreak and -breaktx exclude each other")
	}
	if *autobreak && *startBlock != 0 {
		return errors.New("options -autobreak and -startblock exclude each other")
	}
	if *autobreak && *startTx != 0 {
		return errors.New("options -autobreak and -starttx exclude each other")
	}
	if *startBlock != 0 && *breakBlock != -1 {
		return errors.New("options -startblock and -breakblock exclude each other")
	}
	if *startBlock != 0 && *breakTx != 0 {
		return errors.New("options -startblock and -breaktx exclude each other")
	}
	if *startTx != 0 && *breakBlock != -1 {
		return errors.New("options -starttx and -breakblock exclude each other")
	}
	if *startTx != 0 && *breakTx != 0 {
		return errors.New("options -starttx and -breaktx exclude each other")
	}
	if *release && !*setup {
		return errors.New("option -release requires option -setup")
	}
	if *setup && *contract == "" {
		return errors.New("option -setup requires option -contract")
	}
	if *contract != "" && !*setup {
		return errors.New("option -contract requires option -setup")
	}
	if *neardPath != "" && *neardHead == "" {
		return errors.New("option -neard requires option -neardhead")
	}
	chainID, testnet, err := testnetFlags.determineTestnet()
	if err != nil {
		return err
	}
	if !*setup {
		if fs.NArg() != 1 {
			fs.Usage()
			return flag.ErrHelp
		}
	} else {
		if fs.NArg() > 1 {
			fs.Usage()
			return flag.ErrHelp
		}
	}

	// determine evmContract
	var evmContract string
	if fs.NArg() == 1 {
		evmContract = fs.Arg(0)
	} else {
		var b [16]byte
		if _, err := rand.Read(b[:]); err != nil {
			return err
		}
		evmContract = hex.EncodeToString(b[:]) + ".test.near"
		fmt.Fprintf(os.Stderr, "evmContract name generated: %s\n", evmContract)
	}

	// set accountID, if necessary
	if *accountID == "" {
		*accountID = evmContract
	}

	// run replayer
	r := replayer.Replayer{
		Config:         cfg,
		Timeout:        *timeout,
		ChainID:        chainID,
		Gas:            *gas,
		DataDir:        *dataDir,
		Testnet:        testnet,
		Defrost:        *defrost,
		Skip:           *skip,
		Batch:          *batch,
		BatchSize:      *batchSize,
		StartBlock:     *startBlock,
		StartTx:        *startTx,
		Autobreak:      *autobreak,
		BreakBlock:     *breakBlock,
		BreakTx:        *breakTx,
		Release:        *release,
		Setup:          *setup,
		NeardPath:      *neardPath,
		NeardHead:      *neardHead,
		InitialBalance: *initialBalance,
		Contract:       *contract,
		Breakpoint: replayer.Breakpoint{
			AccountID: *accountID,
		},
	}
	if err := r.Replay(evmContract); err != nil {
		return err
	}
	return nil
}
