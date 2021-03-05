package nearapi

import (
	"bytes"
	"crypto/ed25519"
	"encoding/json"
	"fmt"
	"io/ioutil"
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
}

func FindAccessKey(receiverID string) (ed25519.PrivateKey, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	// TODO: extend this function to allow loading from "local" as well
	fn := filepath.Join(home, ".near-credentials", "default", receiverID+".json")
	return ReadAccessKey(fn, receiverID)
}

func ReadAccessKey(filename, receiverID string) (ed25519.PrivateKey, error) {
	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var a Account
	err = json.Unmarshal(buf, &a)
	if err != nil {
		return nil, err
	}
	// account ID
	if a.AccountID != receiverID {
		return nil, fmt.Errorf("nearapi: parsed account_id '%s' does not match with receiverID '%s'",
			a.AccountID, receiverID)
	}
	// public key
	if !strings.HasPrefix(a.PublicKey, ed25519Prefix) {
		return nil, fmt.Errorf("nearapi: parsed public_key '%s' is not an Ed25519 key",
			a.PublicKey)
	}
	pubKey := base58.Decode(strings.TrimPrefix(a.PublicKey, ed25519Prefix))
	// private key
	if !strings.HasPrefix(a.PrivateKey, ed25519Prefix) {
		return nil, fmt.Errorf("nearapi: parsed private_key '%s' is not an Ed25519 key",
			a.PrivateKey)
	}
	privKey := base58.Decode(strings.TrimPrefix(a.PrivateKey, ed25519Prefix))
	if err != nil {
		return nil, err
	}
	privateKey := ed25519.PrivateKey(privKey)
	// make sure keys match
	if !bytes.Equal(pubKey, privateKey.Public().(ed25519.PublicKey)) {
		return nil, fmt.Errorf("nearapi: public_key does not match private_key: %s", filename)
	}
	return privateKey, nil
}
