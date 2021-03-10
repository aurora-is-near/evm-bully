package nearapi

import (
	"bytes"
	"crypto/ed25519"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
	"path/filepath"
	"strings"

	"github.com/btcsuite/btcutil/base58"
)

const ed25519Prefix = "ed25519:"

// Account defines access credentials for a NEAR account.
type Account struct {
	AccountID  string `json:"account_id"`
	PublicKey  string `json:"public_key"`
	PrivateKey string `json:"private_key"`
	conn       *Connection
	privKey    ed25519.PrivateKey
}

// LoadAccount loads the credential for the receiverID account, to be used via
// connection c, and returns it.
func LoadAccount(c *Connection, receiverID string) (*Account, error) {
	var a Account
	a.conn = c
	if err := a.findAccessKey(receiverID); err != nil {
		return nil, err
	}
	return &a, nil
}

func (a *Account) findAccessKey(receiverID string) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	// TODO: extend this function to allow loading from "local" as well
	fn := filepath.Join(home, ".near-credentials", "default", receiverID+".json")
	return a.readAccessKey(fn, receiverID)
}

func (a *Account) readAccessKey(filename, receiverID string) error {
	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	err = json.Unmarshal(buf, &a)
	if err != nil {
		return err
	}
	// account ID
	if a.AccountID != receiverID {
		return fmt.Errorf("nearapi: parsed account_id '%s' does not match with receiverID '%s'",
			a.AccountID, receiverID)
	}
	// public key
	if !strings.HasPrefix(a.PublicKey, ed25519Prefix) {
		return fmt.Errorf("nearapi: parsed public_key '%s' is not an Ed25519 key",
			a.PublicKey)
	}
	pubKey := base58.Decode(strings.TrimPrefix(a.PublicKey, ed25519Prefix))
	// private key
	if !strings.HasPrefix(a.PrivateKey, ed25519Prefix) {
		return fmt.Errorf("nearapi: parsed private_key '%s' is not an Ed25519 key",
			a.PrivateKey)
	}
	privateKey := base58.Decode(strings.TrimPrefix(a.PrivateKey, ed25519Prefix))
	if err != nil {
		return err
	}
	a.privKey = ed25519.PrivateKey(privateKey)
	// make sure keys match
	if !bytes.Equal(pubKey, a.privKey.Public().(ed25519.PublicKey)) {
		return fmt.Errorf("nearapi: public_key does not match private_key: %s", filename)
	}
	return nil
}

// SendMoney sends amount NEAR from account to receiverID.
func (a *Account) SendMoney(receiverID string, amount *big.Int) error {
	// TODO
	return nil
}
