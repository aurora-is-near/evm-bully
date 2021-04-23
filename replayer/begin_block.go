package replayer

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"

	"github.com/aurora-is-near/evm-bully/nearapi"
	"github.com/aurora-is-near/evm-bully/nearapi/utils"
	"github.com/near/borsh-go"
)

// BeginBlockArgs encodes the arguments for 'begin_block'.
type BeginBlockArgs struct {
	Hash       RawU256
	Coinbase   RawAddress
	Timestamp  RawU256
	Number     RawU256
	Difficulty RawU256
	Gaslimit   RawU256
}

func beginBlock(
	a *nearapi.Account,
	evmContract string,
	gas uint64,
	c *blockContext,
) error {
	zeroAmount := big.NewInt(0)

	fmt.Printf("begin_block(%d)\n", c.number)

	var args BeginBlockArgs
	copy(args.Hash[:], c.hash[:])
	copy(args.Coinbase[:], c.coinbase[:])
	binary.LittleEndian.PutUint64(args.Timestamp[:], c.timestamp)
	binary.LittleEndian.PutUint64(args.Number[:], c.number)
	binary.LittleEndian.PutUint64(args.Difficulty[:], c.difficulty)
	binary.LittleEndian.PutUint64(args.Gaslimit[:], c.gaslimit)

	data, err := borsh.Serialize(args)
	if err != nil {
		return err
	}

	txResult, err := a.FunctionCall(evmContract, "begin_block", data, gas, *zeroAmount)
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
