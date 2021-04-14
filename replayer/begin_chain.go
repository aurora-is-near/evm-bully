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

type RawU256 [32]uint8

type BeginChainArgs struct {
	ChainID RawU256
}

func beginChain(
	chainID uint8,
	a *nearapi.Account,
	evmContract string,
	gas uint64,
	g *core.Genesis,
) error {
	zeroAmount := big.NewInt(0)

	fmt.Println("begin_chain()")

	var args BeginChainArgs
	args.ChainID[0] = chainID

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
