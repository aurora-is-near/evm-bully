package command

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/aurora-is-near/evm-bully/replayer"
	"github.com/aurora-is-near/near-api-go"
)

// CreateAccount implements the 'create-account' command.
func CreateAccount(argv0 string, args ...string) error {
	fs := flag.NewFlagSet(argv0, flag.ContinueOnError)
	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s <accountId>\n", argv0)
		fmt.Fprintf(os.Stderr, "Create a new developer account (subaccount of the masterAccount, ex: app.alice.test)\n")
		fs.PrintDefaults()
	}
	initialBalance := fs.String("initial-balance", "100", "Number of tokens to transfer to newly created account")
	masterAccount := fs.String("master-account", "", "Account used to create requested account")
	cfg := near.GetConfig()
	registerCfgFlag(fs, cfg)
	if err := fs.Parse(args); err != nil {
		return err
	}
	if fs.NArg() != 1 {
		fs.Usage()
		return flag.ErrHelp
	}
	if *masterAccount == "" {
		return errors.New("option -master-account is mandatory")
	}
	accountID := fs.Arg(0)
	// create account
	c := replayer.CreateAccount{
		Config:         cfg,
		InitialBalance: *initialBalance,
		MasterAccount:  *masterAccount,
	}
	return c.Create(accountID)
}
