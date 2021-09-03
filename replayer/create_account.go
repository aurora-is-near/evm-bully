package replayer

import (
	"fmt"
	"os"
	"strings"

	"github.com/aurora-is-near/near-api-go"
	"github.com/aurora-is-near/near-api-go/keystore"
	"github.com/aurora-is-near/near-api-go/utils"
)

// TLAMinLength defines the minimum length for top-level accounts.
const TLAMinLength = 32

// CreateAccount allow to create accounts.
type CreateAccount struct {
	Config         *near.Config
	InitialBalance string
	MasterAccount  string
}

// Create account with accountID.
func (ca *CreateAccount) Create(accountID string) error {
	splitAccount := strings.Split(accountID, ".")
	if len(splitAccount) == 1 {
		// TLA (bob-with-at-least-maximum-characters)
		if len(splitAccount[0]) < TLAMinLength {
			return fmt.Errorf("top-level accounts must be at least %d characters.\n"+
				"Note: this is for advanced usage only. Typical account names are of the form:\n"+
				"app.alice.test, where the masterAccount shares the top-level account (.test).",
				TLAMinLength)
		}
	}
	// Subaccounts (short.alice.near, even.more.bob.test, and eventually peter.potato)
	// Check that master account TLA matches
	if !strings.HasSuffix(accountID, ca.MasterAccount) {
		return fmt.Errorf("new account doesn't share the same top-level account. Expecting account name to end in %s",
			ca.MasterAccount)
	}
	c := near.NewConnection(ca.Config.NodeURL)

	// generate key
	kp, err := keystore.GenerateEd25519KeyPair(accountID)
	if err != nil {
		return err
	}

	// check to see if account already exists
	_, err = c.State(accountID)
	if err == nil {
		return fmt.Errorf("sorry, account '%s' already exists", accountID)
	} else if !strings.Contains(err.Error(), "does not exist while viewing") {
		return err
	}

	amnt, err := utils.ParseNearAmountAsBigInt(ca.InitialBalance)
	if err != nil {
		return err
	}
	a, err := near.LoadAccount(c, ca.Config, ca.MasterAccount)
	if err != nil {
		fmt.Println("error")
		return err
	}

	// save key
	filename, err := kp.Write(ca.Config.NetworkID)
	if err != nil {
		return err
	}
	fmt.Printf("saving key to '%s'\n", filename)

	// create account
	txResult, err := a.CreateAccount(accountID, utils.PublicKeyFromEd25519(kp.Ed25519PubKey), *amnt)
	if err != nil {
		if err == utils.ErrRetriesExceeded {
			fmt.Fprintf(os.Stderr, "Received a timeout when creating account, please run:\n")
			fmt.Fprintf(os.Stderr, "`evm state %s\n", accountID)
			fmt.Fprintf(os.Stderr, "to confirm creation. Keyfile for this account has been saved.\n")
		} else {
			// remove keyfile
			os.Remove(filename)
		}
		return err
	}
	utils.PrettyPrintResponse(txResult)
	fmt.Printf("account %s for network \"%s\" was created\n", accountID, ca.Config.NetworkID)
	return nil
}
