package command

import (
  "flag"
  "fmt"
  "os"
)

// BlockNumber implements the 'blocknumber' command.
func BlockNumber(net, argv0 string, args ...string) error {
  fs := flag.NewFlagSet(argv0, flag.ContinueOnError)
  fs.Usage = func() {
    fmt.Fprintf(os.Stderr, "Usage: %s\n", argv0)
    fmt.Fprintf(os.Stderr, "TODO\n")
    fs.PrintDefaults()
  }
  if err := fs.Parse(args); err != nil {
    return err
  }
  if fs.NArg() != 0 {
    fs.Usage()
    return flag.ErrHelp
  }
  // TODO
  return nil
}
