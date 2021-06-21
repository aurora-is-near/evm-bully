package replayer

import (
	"fmt"
	"sort"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
)

func getGenesisBlock(net string) *core.Genesis {
	switch net {
	case "goerli":
		return core.DefaultGoerliGenesisBlock()
	case "rinkeby":
		return core.DefaultRinkebyGenesisBlock()
	case "ropsten":
		return core.DefaultRopstenGenesisBlock()
	default:
		return core.DefaultGenesisBlock()
	}
}

type AddrSlice []common.Address

func (s AddrSlice) Len() int           { return len(s) }
func (s AddrSlice) Less(i, j int) bool { return s[i].String() < s[j].String() }
func (s AddrSlice) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

func dumpAccounts(g *core.Genesis) error {
	var addresses AddrSlice
	accounts := make(map[common.Address]string)
	for address, account := range g.Alloc {
		addresses = append(addresses, address)
		accounts[address] = account.Balance.String()
	}
	sort.Sort(addresses)
	for _, a := range addresses {
		fmt.Printf("%s: %s\n", a.String(), accounts[a])
	}
	return nil
}

// ProcGenesisBlock processes the genesis block for the given testnet.
func ProcGenesisBlock(testnet string) error {
	g := getGenesisBlock(testnet)
	return dumpAccounts(g)
}
