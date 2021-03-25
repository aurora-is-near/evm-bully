package replayer

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/big"

	"github.com/aurora-is-near/evm-bully/nearapi"
	"github.com/aurora-is-near/evm-bully/nearapi/utils"
	"github.com/ethereum/go-ethereum/core"
	"github.com/near/borsh-go"
)

type BeginChainArgs struct {
	ChainID uint64 // TODO: use correct type for aurora-engine compatibility
}

func beginChain(
	chainID uint64,
	a *nearapi.Account,
	evmContract string,
	gas uint64,
	g *core.Genesis,
) error {
	zeroAmount := big.NewInt(0)

	fmt.Println("begin_chain()")

	args := BeginChainArgs{
		ChainID: chainID,
	}

	data, err := borsh.Serialize(args)
	if err != nil {
		return err
	}

	txResult, err := a.FunctionCall(evmContract, "begin_chain", data, gas, *zeroAmount)
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

	return nil
}
