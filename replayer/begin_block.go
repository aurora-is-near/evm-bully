package replayer

import (
	"encoding/binary"
	"fmt"

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

func beginBlockTx(gas uint64, c *blockContext) *Tx {
	var args BeginBlockArgs
	copy(args.Hash[:], c.hash[:])
	copy(args.Coinbase[:], c.coinbase[:])
	binary.LittleEndian.PutUint64(args.Timestamp[:], c.timestamp)
	binary.LittleEndian.PutUint64(args.Number[:], c.number)
	binary.LittleEndian.PutUint64(args.Difficulty[:], c.difficulty)
	binary.LittleEndian.PutUint64(args.Gaslimit[:], c.gaslimit)

	data, err := borsh.Serialize(args)
	if err != nil {
		return &Tx{Error: err}
	}

	return &Tx{
		Comment:    fmt.Sprintf("begin_block(%d)", c.number),
		MethodName: "begin_block",
		Args:       data,
	}
}
