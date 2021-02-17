package replayer

import (
	"fmt"
	"sort"

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

func dumpAccounts(g *core.Genesis) {
	var addresses []string
	accounts := make(map[string]string)
	for address, account := range g.Alloc {
		a := address.String()
		addresses = append(addresses, a)
		accounts[a] = account.Balance.String()
	}
	sort.Strings(addresses)
	for _, a := range addresses {
		fmt.Printf("%s: %s\n", a, accounts[a])
	}
}

// ProcGenesisBlock processes the genesis block for the given testnet.
func ProcGenesisBlock(testnet string) error {
	g := getGenesisBlock(testnet)
	dumpAccounts(g)
	return nil
}
