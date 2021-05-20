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
	var (
		nodeURL      nodeURLFlag
		testnetFlags testnetFlags
	)
	fs := flag.NewFlagSet(argv0, flag.ContinueOnError)
	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s <evmContract>\n", argv0)
		fmt.Fprintf(os.Stderr, "Replay transactions to NEAR EVM.\n")
		fs.PrintDefaults()
	}
	accountID := fs.String("accountId", "", "Unique identifier for the account that will be used to sign this call")
	batch := fs.Bool("batch", false, "Batch transactions")
	batchSize := fs.Int("size", 10, "Batch size when batching transactions")
	block := fs.Uint64("block", defaultGoerliBlockHeight, "Block height to replay to")
	dataDir := fs.String("datadir", defaultDataDir, "Data directory containing the database to read")
	defrost := fs.Bool("defrost", false, "Defrost the database first")
	gas := fs.Uint64("gas", defaultGas, "Max amount of gas a call can use (in gas units)")
	hash := fs.String("hash", defaultGoerliBlockHash, "Block hash to replay to")
	skip := fs.Bool("skip", false, "Skip empty blocks")
	startBlock := fs.Int("startblock", 0, "Start replaying at this block height")
	startTx := fs.Int("starttx", 0, "Start replaying at this transaction (in block given by -startblock)")
	timeout := fs.Duration("timeout", 0, "Timeout for JSON-RPC client")
	cfg := near.GetConfig()
	nodeURL.registerFlag(fs, cfg)
	testnetFlags.registerFlags(fs)
	if err := fs.Parse(args); err != nil {
		return err
	}
	if *accountID == "" {
		return errors.New("option -accountId is mandatory")
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
	// load account
	c := near.NewConnectionWithTimeout(string(nodeURL), *timeout)
	a, err := near.LoadAccount(c, cfg, *accountID)
	if err != nil {
		return err
	}
	// run replayer
	r := replayer.Replayer{
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
	}
	return r.Replay(a, evmContract)
}
