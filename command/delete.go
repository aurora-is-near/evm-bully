package command

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/aurora-is-near/evm-bully/nearapi"
	"github.com/aurora-is-near/evm-bully/nearapi/utils"
)

// Delete implements the 'delete' command.
func Delete(argv0 string, args ...string) error {
	var nodeURL nodeURLFlag
	fs := flag.NewFlagSet(argv0, flag.ContinueOnError)
	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s <accountId> <beneficiaryId>\n", argv0)
		fmt.Fprintf(os.Stderr, "Delete an account and transfer funds to beneficiary account.\n")
		fs.PrintDefaults()
	}
	cfg := nearapi.GetConfig()
	nodeURL.registerFlag(fs, cfg)
	if err := fs.Parse(args); err != nil {
		return err
	}
	if fs.NArg() != 2 {
		fs.Usage()
		return flag.ErrHelp
	}
	accountID := fs.Arg(0)
	beneficiaryID := fs.Arg(1)
	c := nearapi.NewConnection(string(nodeURL))
	a, err := nearapi.LoadAccount(c, cfg, accountID)
	if err != nil {
		return err
	}
	fmt.Printf("Deleting account. Account ID: %s, node: %s, beneficiary: %s\n",
		accountID, string(nodeURL), beneficiaryID)
	txResult, err := a.DeleteAccount(beneficiaryID)
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
