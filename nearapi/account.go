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
	"github.com/near/borsh-go"
)

const ed25519Prefix = "ed25519:"

// Account defines access credentials for a NEAR account.
type Account struct {
	AccountID                 string `json:"account_id"`
	PublicKey                 string `json:"public_key"`
	PrivateKey                string `json:"private_key"`
	conn                      *Connection
	pubKey                    ed25519.PublicKey
	privKey                   ed25519.PrivateKey
	accessKeyByPublicKeyCache map[string]map[string]interface{}
}

// LoadAccount loads the credential for the receiverID account, to be used via
// connection c, and returns it.
func LoadAccount(c *Connection, receiverID string) (*Account, error) {
	var a Account
	a.conn = c
	if err := a.locateAccessKey(receiverID); err != nil {
		return nil, err
	}
	a.accessKeyByPublicKeyCache = make(map[string]map[string]interface{})
	return &a, nil
}

func (a *Account) locateAccessKey(receiverID string) error {
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
	a.pubKey = ed25519.PublicKey(pubKey)
	// private key
	if !strings.HasPrefix(a.PrivateKey, ed25519Prefix) {
		return fmt.Errorf("nearapi: parsed private_key '%s' is not an Ed25519 key",
			a.PrivateKey)
	}
	privateKey := base58.Decode(strings.TrimPrefix(a.PrivateKey, ed25519Prefix))
	a.privKey = ed25519.PrivateKey(privateKey)
	// make sure keys match
	if !bytes.Equal(pubKey, a.privKey.Public().(ed25519.PublicKey)) {
		return fmt.Errorf("nearapi: public_key does not match private_key: %s", filename)
	}
	return nil
}

// SendMoney sends amount NEAR from account to receiverID.
func (a *Account) SendMoney(
	receiverID string,
	amount big.Int,
) (map[string]interface{}, error) {
	return a.signAndSendTransaction(receiverID, []Action{Action{
		Enum:     3,
		Transfer: Transfer{amount},
	}})
}

func (a *Account) signAndSendTransaction(
	receiverID string,
	actions []Action,
) (map[string]interface{}, error) {
	// TODO: exponential backoff
	_, signedTx, err := a.signTransaction(receiverID, actions)
	if err != nil {
		return nil, err
	}

	buf, err := borsh.Serialize(*signedTx)
	if err != nil {
		return nil, err
	}

	return a.conn.SendTransaction(buf)
}

func (a *Account) signTransaction(
	receiverID string,
	actions []Action,
) (txHash []byte, signedTx *SignedTransaction, err error) {
	_, ak, err := a.findAccessKey()
	if err != nil {
		return nil, nil, err
	}

	// get current block hash
	block, err := a.conn.Block()
	if err != nil {
		return nil, nil, err
	}
	blockHash := block["header"].(map[string]interface{})["hash"].(string)

	// create next nonce
	nonce, err := ak["nonce"].(json.Number).Int64()
	if err != nil {
		return nil, nil, err
	}
	nonce++

	// sign transaction
	return signTransaction(receiverID, uint64(nonce), actions, base58.Decode(blockHash),
		a.pubKey, a.privKey, a.AccountID)

}

func (a *Account) findAccessKey() (publicKey ed25519.PublicKey, accessKey map[string]interface{}, err error) {
	// TODO: Find matching access key based on transaction
	// TODO: use accountId and networkId?
	pk := a.pubKey
	if ak := a.accessKeyByPublicKeyCache[string(publicKey)]; ak != nil {
		return pk, ak, nil
	}
	ak, err := a.conn.ViewAccessKey(a.AccountID, a.PublicKey)
	if err != nil {
		fmt.Println("ERROR")
		return nil, nil, err
	}
	a.accessKeyByPublicKeyCache[string(publicKey)] = ak
	return pk, ak, nil
}
