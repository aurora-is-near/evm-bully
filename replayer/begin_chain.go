package replayer

import (
	"fmt"

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
		b, err := bigIntToRawU256(account.Balance)
		if err != nil {
			return nil, err
		}
		ab.Balance = b
		ga = append(ga, ab)
	}
	return ga, nil
}

func (r *Replayer) beginChainTx(g *core.Genesis) *Tx {
	var args BeginChainArgs
	var err error
	args.ChainID[31] = r.ChainID
	args.GenesisAlloc, err = genesisAlloc(g)
	if err != nil {
		return &Tx{Error: err}
	}
	data, err := borsh.Serialize(args)
	if err != nil {
		return &Tx{Error: err}
	}
	return &Tx{
		Comment:    fmt.Sprintf("begin_chain()"),
		MethodName: "begin_chain",
		Args:       data,
	}
}
