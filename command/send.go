package command

import (
  "flag"
  "fmt"
  "os"

  "github.com/aurora-is-near/evm-bully/nearapi"
)

// Send implements the 'send' command.
func Send(argv0 string, args ...string) error {
  var nodeURL nodeURLFlag
  fs := flag.NewFlagSet(argv0, flag.ContinueOnError)
  fs.Usage = func() {
    fmt.Fprintf(os.Stderr, "Usage: %s <sender> <receiver> <amount>\n", argv0)
    fmt.Fprintf(os.Stderr, "Send tokens to given receiver.\n")
    fs.PrintDefaults()
  }
  nodeURL.registerFlag(fs)
  if err := fs.Parse(args); err != nil {
    return err
  }
  if fs.NArg() != 3 {
    fs.Usage()
    return flag.ErrHelp
  }
  sender := fs.Arg(0)
  receiver := fs.Arg(1)
  amount := fs.Arg(2)
  fmt.Printf("%s %s %s %s\n", nodeURL, sender, receiver, amount)
  c := nearapi.NewConnection(string(nodeURL))
  _, err := nearapi.LoadAccount(c, sender)
  if err != nil {
    return err
  }
  // TODO: send money
  return nil
}
