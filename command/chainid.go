package command

import (
  "context"
  "flag"
  "fmt"
  "os"

  "github.com/near/evm-bully/replayer"
)

// ChainID implements the 'chainid' command.
func ChainID(argv0 string, args ...string) error {
  fs := flag.NewFlagSet(argv0, flag.ContinueOnError)
  fs.Usage = func() {
    fmt.Fprintf(os.Stderr, "Usage: %s\n", argv0)
    fmt.Fprintf(os.Stderr, "Print the chain ID to stdout.\n")
    fs.PrintDefaults()
  }
  endpoint := fs.String("endpoint", defaultEndpoint, "Set JSON-RPC endpoint")
  if err := fs.Parse(args); err != nil {
    return err
  }
  if fs.NArg() != 0 {
    fs.Usage()
    return flag.ErrHelp
  }
  cID, err := replayer.ChainID(context.Background(), *endpoint)
  if err != nil {
    return err
  }
  fmt.Printf("%d\n", cID)
  return nil
}
