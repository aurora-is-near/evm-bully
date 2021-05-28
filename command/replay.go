package command

import (
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
		fmt.Fprintf(os.Stderr, "Usage: %s <evmContract>\n", argv0)
		fmt.Fprintf(os.Stderr, "Replay transactions to NEAR EVM.\n")
		fs.PrintDefaults()
	}
	accountID := fs.String("accountId", "", "Unique identifier for the account that will be used to sign this call")
	batch := fs.Bool("batch", false, "Batch transactions")
	batchSize := fs.Int("size", 10, "Batch size when batching transactions")
	breakBlock := fs.Int("breakblock", 0, "Break replaying at this block height")
	breakTx := fs.Int("breaktx", 0, "Break replaying at this transaction (in block given by -breakblock)")
	block := fs.Uint64("block", defaultGoerliBlockHeight, "Block height to replay to")
	dataDir := fs.String("datadir", defaultDataDir, "Data directory containing the database to read")
	defrost := fs.Bool("defrost", false, "Defrost the database first")
	gas := fs.Uint64("gas", defaultGas, "Max amount of gas a call can use (in gas units)")
	hash := fs.String("hash", defaultGoerliBlockHash, "Block hash to replay to")
	release := fs.Bool("release", false, "Run release version of neard")
	setup := fs.Bool("setup", false, "Setup and run neard before replaying")
	skip := fs.Bool("skip", false, "Skip empty blocks")
	startBlock := fs.Int("startblock", 0, "Start replaying at this block height")
	startTx := fs.Int("starttx", 0, "Start replaying at this transaction (in block given by -startblock)")
	timeout := fs.Duration("timeout", 0, "Timeout for JSON-RPC client")
	cfg := near.GetConfig()
	registerCfgFlags(fs, cfg, true)
	testnetFlags.registerFlags(fs)
	if err := fs.Parse(args); err != nil {
		return err
	}
	if *accountID == "" {
		return errors.New("option -accountId is mandatory")
	}
	if *startBlock != 0 && *breakBlock != 0 {
		return errors.New("options -startblock and -breakblock exclude each other")
	}
	if *startBlock != 0 && *breakTx != 0 {
		return errors.New("options -startblock and -breaktx exclude each other")
	}
	if *startTx != 0 && *breakBlock != 0 {
		return errors.New("options -starttx and -breakblock exclude each other")
	}
	if *startTx != 0 && *breakTx != 0 {
		return errors.New("options -starttx and -breaktx exclude each other")
	}
	if *release && !*setup {
		return errors.New("option -release requires option -setup")
	}
	chainID, testnet, err := testnetFlags.determineTestnet()
	if err != nil {
		return err
	}
	adjustBlockDefaults(block, hash, testnet)
	if fs.NArg() != 1 {
		fs.Usage()
		return flag.ErrHelp
	}
	evmContract := fs.Arg(0)

	// run replayer
	r := replayer.Replayer{
		AccountID:   *accountID,
		Config:      cfg,
		Timeout:     *timeout,
		ChainID:     chainID,
		Gas:         *gas,
		DataDir:     *dataDir,
		Testnet:     testnet,
		BlockHeight: *block,
		BlockHash:   *hash,
		Defrost:     *defrost,
		Skip:        *skip,
		Batch:       *batch,
		BatchSize:   *batchSize,
		StartBlock:  *startBlock,
		StartTx:     *startTx,
		BreakBlock:  *breakBlock,
		BreakTx:     *breakTx,
		Release:     *release,
		Setup:       *setup,
	}
	return r.Replay(evmContract)
}
