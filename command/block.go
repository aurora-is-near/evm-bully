package command

import (
  "encoding/json"
  "flag"
  "fmt"
  "os"

  "github.com/aurora-is-near/evm-bully/nearapi"
)

// Block implements the 'block' command.
func Block(argv0 string, args ...string) error {
  fs := flag.NewFlagSet(argv0, flag.ContinueOnError)
  fs.Usage = func() {
    fmt.Fprintf(os.Stderr, "Usage: %s\n", argv0)
    fmt.Fprintf(os.Stderr, "Queries network for latest block details.\n")
    fs.PrintDefaults()
  }
  nodeURL := fs.String("nodeUrl", defaultNodeURL, "NEAR node URL")
  if err := fs.Parse(args); err != nil {
    return err
  }
  if fs.NArg() != 0 {
    fs.Usage()
    return flag.ErrHelp
  }
  c := nearapi.NewConnection(*nodeURL)
  res, err := c.Block()
  if err != nil {
    return err
  }
  jsn, err := json.MarshalIndent(res, "", "  ")
  if err != nil {
    return err
  }
  fmt.Println(string(jsn))
  return nil
}
