package command

import (
  "context"
  "flag"
  "fmt"
  "os"

  "github.com/near/evm-bully/replayer"
)

// Replay implements the 'replay' command.
func Replay(net, argv0 string, args ...string) error {
  fs := flag.NewFlagSet(argv0, flag.ContinueOnError)
  fs.Usage = func() {
    fmt.Fprintf(os.Stderr, "Usage: %s\n", argv0)
    fmt.Fprintf(os.Stderr, "Replay transactions to NEAR EVM.\n")
    fs.PrintDefaults()
  }
  block := fs.Uint64("block", defaultBlockHeight, "Block height")
  datadir := fs.String("datadir", defaultDataDir, "Data directory containing the database to read")
  endpoint := fs.String("endpoint", defaultEndpoint, "Set default JSON-RPC endpoint")
  hash := fs.String("hash", defaultBlockhash, "Block hash")
  if err := fs.Parse(args); err != nil {
    return err
  }
  if fs.NArg() != 0 {
    fs.Usage()
    return flag.ErrHelp
  }
  // run replayer
  return replayer.ReadTxs(context.Background(), *endpoint, *datadir, net, *block, *hash)
}
