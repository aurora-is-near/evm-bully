package nearapi

import (
	"crypto"
	"crypto/ed25519"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/big"

	"github.com/aurora-is-near/evm-bully/nearapi/utils"
	"github.com/near/borsh-go"
)

type Transfer struct {
	Deposit *big.Int // TODO
}

type Transaction struct {
	SignerID   string
	PublicKey  *utils.PublicKey
	Nonce      int64
	ReceiverID string
	Actions    []interface{}
	BlockHash  []byte
}

type Signature struct {
	KeyType int
	Data    []byte
}

type SignedTransaction struct {
	Transaction *Transaction
	Signature   *Signature
}

func createTransaction(
	signerID string,
	publicKey *utils.PublicKey,
	receiverID string,
	nonce int64,
	actions []interface{},
	blockHash []byte,
) *Transaction {
	return &Transaction{
		SignerID:   signerID,
		PublicKey:  publicKey,
		ReceiverID: receiverID,
		Nonce:      nonce,
		Actions:    actions,
		BlockHash:  blockHash,
	}
}

func signTransactionObject(
	tx *Transaction,
	privKey ed25519.PrivateKey,
	accountID string,
) (txHash []byte, signedTx *SignedTransaction, err error) {
	buf, err := borsh.Serialize(tx)
	if err != nil {
		return nil, nil, err
	}

	fmt.Printf("tx=%s\n", hex.EncodeToString(buf))

	hash := sha256.Sum256(buf)

	sig, err := privKey.Sign(rand.Reader, hash[:], crypto.Hash(0))
	if err != nil {
		return nil, nil, err
	}

	stx := &SignedTransaction{
		Transaction: tx,
		Signature: &Signature{
			KeyType: utils.ED25519,
			Data:    sig,
		},
	}

	return hash[:], stx, nil
}

func signTransaction(
	receiverID string,
	nonce int64,
	actions []interface{},
	blockHash []byte,
	publicKey ed25519.PublicKey,
	privKey ed25519.PrivateKey,
	accountID string,
) (txHash []byte, signedTx *SignedTransaction, err error) {
	// create transaction
	tx := createTransaction(accountID, utils.PublicKeyFromEd25519(publicKey),
		receiverID, nonce, actions, blockHash)

	// sign transaction object
	txHash, signedTx, err = signTransactionObject(tx, privKey, accountID)
	if err != nil {
		return nil, nil, err
	}
	return txHash, signedTx, nil
}
