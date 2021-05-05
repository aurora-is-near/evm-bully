package replayer

import (
	"github.com/ethereum/go-ethereum/core/types"
)

// Tx defines a replayer transaction.
type Tx struct {
	Comment    string             // comment for the transaction
	MethodName string             // the Aurora Engine method name to call
	Args       []byte             // the argument to call the method with
	EthTx      *types.Transaction // pointer to original Ethereum transaction (for 'submit')
	Error      error              // error during transaction construction
}
