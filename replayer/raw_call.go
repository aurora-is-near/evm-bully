package replayer

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/big"

	"github.com/aurora-is-near/evm-bully/nearapi"
	"github.com/aurora-is-near/evm-bully/nearapi/utils"
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
	for i, tx := range txs {
		// get signed transaction in RLP encoding
		rlp, err := tx.MarshalBinary()
		if err != nil {
			return err
		}

		fmt.Printf("raw_call(%d, tx=%d, tx_size=%d)\n", blockHeight, i, len(rlp))
		txResult, err := a.FunctionCall(evmContract, "raw_call", rlp, gas, *zeroAmount)
		if err != nil {
			return err
		}
		utils.PrettyPrintResponse(txResult)
		status := txResult["status"].(map[string]interface{})
		jsn, err := json.MarshalIndent(status, "", "  ")
		if err != nil {
			return err
		}
		fmt.Println(string(jsn))
		if status["Failure"] != nil {
			return errors.New("replayer: transaction failed")
		}
	}
	return nil
}
