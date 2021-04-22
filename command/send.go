package command

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/aurora-is-near/evm-bully/nearapi"
	"github.com/aurora-is-near/evm-bully/nearapi/utils"
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
	c := nearapi.NewConnection(string(nodeURL))
	a, err := nearapi.LoadAccount(c, nearapi.GetConfig(), sender)
	if err != nil {
		return err
	}
	amnt, err := utils.ParseNearAmountAsBigInt(amount)
	if err != nil {
		return err
	}
	fmt.Printf("Sending %s NEAR to %s from %s\n", amount, receiver, sender)
	txResult, err := a.SendMoney(receiver, *amnt)
	if err != nil {
		return err
	}
	utils.PrettyPrintResponse(txResult)
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
