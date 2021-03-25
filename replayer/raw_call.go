package replayer

import (
	"fmt"
	"math/big"

	"github.com/aurora-is-near/evm-bully/nearapi"
	"github.com/ethereum/go-ethereum/core/types"
)

func rawCall(
	a *nearapi.Account,
	evmContract string,
	gas uint64,
	blockHeight int,
	txs types.Transactions,
) error {
	zeroAmount := big.NewInt(0)

	// TODO: batching
	for i, _ := range txs {
		fmt.Printf("raw_call(%d, %d)\n", blockHeight, i)
	}

	// TODO:
	var args []byte

	_, err := a.FunctionCall(evmContract, "raw_call", args, gas, *zeroAmount)
	if err != nil {
		return err
	}

	return nil
}
