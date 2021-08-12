package replayer

import "github.com/aurora-is-near/evm-bully/db"

// Tx defines a replayer transaction.
type Tx struct {
	BlockNum   int             // block number (-1 if undefined)
	TxNum      int             // transaction number
	Comment    string          // comment for the transaction
	MethodName string          // the Aurora Engine method name to call
	Args       []byte          // the argument to call the method with
	EthTx      *db.Transaction // pointer to original Ethereum transaction (for 'submit')
	Error      error           // error during transaction construction
}
