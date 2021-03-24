package command

import (
  "encoding/base64"
  "encoding/json"
  "errors"
  "flag"
  "fmt"
  "os"

  "github.com/aurora-is-near/evm-bully/nearapi"
  "github.com/aurora-is-near/evm-bully/nearapi/utils"
)

func parseArgs(args string, base64enc bool) ([]byte, error) {
  if base64enc {
    b, err := base64.URLEncoding.DecodeString(args)
    if err != nil {
      return nil, err
    }
    return b, nil
  } else {
    var obj map[string]interface{}
    if err := json.Unmarshal([]byte(args), &obj); err != nil {
      return nil, err
    }
    jsn, err := json.Marshal(&obj)
    if err != nil {
      return nil, err
    }
    return jsn, nil
  }
}

// Call implements the 'call' command.
func Call(argv0 string, args ...string) error {
  var nodeURL nodeURLFlag
  fs := flag.NewFlagSet(argv0, flag.ContinueOnError)
  fs.Usage = func() {
    fmt.Fprintf(os.Stderr, "Usage: %s <contractName> <methodName>\n", argv0)
    fmt.Fprintf(os.Stderr, "Schedule smart contract call which can modify state.\n")
    fs.PrintDefaults()
  }
  nodeURL.registerFlag(fs)
  accountID := fs.String("accountId", "", "Unique identifier for the account that will be used to sign this call")
  contractArgs := fs.String("args", "{}", "Arguments to the contract call, in JSON format by default (e.g. '{\"param_a\": \"value\"}')")
  base64enc := fs.Bool("base64", false, "Treat arguments as base64-encoded BLOB")
  gas := fs.Uint64("gas", 100000000000000, "Max amount of gas this call can use (in gas units)")
  amount := fs.String("amount", "0", "Number of tokens to attach (in NEAR)")
  if err := fs.Parse(args); err != nil {
    return err
  }
  if *accountID == "" {
    return errors.New("option -accountId is mandatory")
  }
  if fs.NArg() != 2 {
    fs.Usage()
    return flag.ErrHelp
  }
  contract := fs.Arg(0)
  method := fs.Arg(1)
  c := nearapi.NewConnection(string(nodeURL))
  a, err := nearapi.LoadAccount(c, *accountID)
  if err != nil {
    return err
  }
  amnt, err := utils.ParseNearAmountAsBigInt(*amount)
  if err != nil {
    return err
  }
  parsedArgs, err := parseArgs(*contractArgs, *base64enc)
  if err != nil {
    return err
  }
  txResult, err := a.FunctionCall(contract, method, parsedArgs, *gas, *amnt)
  if err != nil {
    return err
  }
  res, err := nearapi.GetTransactionLastResult(txResult)
  if err != nil {
    return err
  }
  if res != nil {
    jsn, err := json.MarshalIndent(res, "", "  ")
    if err != nil {
      return err
    }
    fmt.Println(string(jsn))
  }
  return nil
}
