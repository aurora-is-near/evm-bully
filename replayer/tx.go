package replayer

// Tx defines a replayer transaction.
type Tx struct {
	Comment    string // comment for the transaction
	MethodName string // the Aurora Engine method name to call
	Args       []byte // the argument to call the method with
	Error      error  // error during transaction construction
}
