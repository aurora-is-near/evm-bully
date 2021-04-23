package replayer

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"

	"github.com/aurora-is-near/evm-bully/nearapi"
	"github.com/aurora-is-near/evm-bully/nearapi/utils"
	"github.com/ethereum/go-ethereum/core"
	"github.com/near/borsh-go"
)

// AccountBalance encodes (genesis) account balances used by the 'begin_chain' function.
type AccountBalance struct {
	Account RawAddress
	Balance RawU256
}

// BeginChainArgs encodes the arguments for 'begin_chain'.
type BeginChainArgs struct {
	ChainID      RawU256
	GenesisAlloc []AccountBalance
}

func genesisAlloc(g *core.Genesis) ([]AccountBalance, error) {
	ga := make([]AccountBalance, 0, len(g.Alloc))
	for address, account := range g.Alloc {
		var ab AccountBalance
		copy(ab.Account[:], address[:])
		b, err := bigIntToUint64(account.Balance)
		if err != nil {
			return nil, err
		}
		binary.LittleEndian.PutUint64(ab.Balance[:], b)
		ga = append(ga, ab)
	}
	return ga, nil
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
	var err error
	args.ChainID[31] = chainID
	args.GenesisAlloc, err = genesisAlloc(g)
	if err != nil {
		return err
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
