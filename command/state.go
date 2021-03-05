package command

import (
  "encoding/json"
  "flag"
  "fmt"
  "os"

  "github.com/near/evm-bully/nearapi"
  "github.com/near/evm-bully/nearapi/utils"
)

// State implements the 'state' command.
func State(argv0 string, args ...string) error {
  fs := flag.NewFlagSet(argv0, flag.ContinueOnError)
  fs.Usage = func() {
    fmt.Fprintf(os.Stderr, "Usage: %s <accountId>\n", argv0)
    fmt.Fprintf(os.Stderr, "View account state.\n")
    fs.PrintDefaults()
  }
  nodeURL := fs.String("nodeUrl", defaultNodeURL, "NEAR node URL")
  if err := fs.Parse(args); err != nil {
    return err
  }
  if fs.NArg() != 1 {
    fs.Usage()
    return flag.ErrHelp
  }
  accountID := fs.Arg(0)
  c := nearapi.NewConnection(*nodeURL)
  res, err := c.State(accountID)
  if err != nil {
    return err
  }
  m := res.(map[string]interface{})
  amount, ok := m["amount"].(string)
  if ok {
    fa, err := utils.FormatNearAmount(amount)
    if err != nil {
      return err
    }
    m["formattedAmount"] = fa
  }
  jsn, err := json.MarshalIndent(res, "", "  ")
  if err != nil {
    return err
  }
  fmt.Println(string(jsn))
  return nil
}
