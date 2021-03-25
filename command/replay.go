package command

import (
  "context"
  "flag"
  "fmt"
  "os"
  "path/filepath"

  "github.com/aurora-is-near/evm-bully/replayer"
  "github.com/frankbraun/codechain/util/homedir"
)

// Replay implements the 'replay' command.
func Replay(argv0 string, args ...string) error {
  var f testnetFlags
  fs := flag.NewFlagSet(argv0, flag.ContinueOnError)
  fs.Usage = func() {
    fmt.Fprintf(os.Stderr, "Usage: %s\n", argv0)
    fmt.Fprintf(os.Stderr, "Replay transactions to NEAR EVM.\n")
    fs.PrintDefaults()
  }
  block := fs.Uint64("block", defaultBlockHeight, "Block height")
  dataDir := fs.String("datadir", defaultDataDir, "Data directory containing the database to read")
  defrost := fs.Bool("defrost", false, "Defrost the database first")
  endpoint := fs.String("endpoint", defaultEndpoint, "Set JSON-RPC endpoint")
  hash := fs.String("hash", defaultBlockhash, "Block hash")
  f.registerFlags(fs)
  if err := fs.Parse(args); err != nil {
    return err
  }
  testnet, err := f.determineTestnet()
  if err != nil {
    return err
  }
  if fs.NArg() != 0 {
    fs.Usage()
    return flag.ErrHelp
  }
  // determine cache directory
  homeDir := homedir.Get("evm-bully")
  cacheDir := filepath.Join(homeDir, testnet)
  // make sure cache directory exists
  if err := os.MkdirAll(cacheDir, 0755); err != nil {
    return err
  }
  // run replayer
  return replayer.Replay(context.Background(), *endpoint, *dataDir, testnet,
    cacheDir, *block, *hash, *defrost)
}
