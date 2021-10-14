package command

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/aurora-is-near/near-api-go"
	"github.com/aurora-is-near/near-api-go/utils"
)

// State implements the 'state' command.
func State(argv0 string, args ...string) error {
	fs := flag.NewFlagSet(argv0, flag.ContinueOnError)
	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s <accountId>\n", argv0)
		fmt.Fprintf(os.Stderr, "View account state.\n")
		fs.PrintDefaults()
	}
	cfg := near.GetConfig()
	registerCfgFlags(fs, cfg, false)
	if err := fs.Parse(args); err != nil {
		return err
	}
	if fs.NArg() != 1 {
		fs.Usage()
		return flag.ErrHelp
	}
	accountID := fs.Arg(0)
	c := near.NewConnection(cfg.NodeURL)
	res, err := c.GetAccountState(accountID)
	if err != nil {
		return err
	}
	amount, ok := res["amount"].(string)
	if ok {
		fa, err := utils.FormatNearAmount(amount)
		if err != nil {
			return err
		}
		res["formattedAmount"] = fa
	}
	jsn, err := json.MarshalIndent(res, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(jsn))
	return nil
}
