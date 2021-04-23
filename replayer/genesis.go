package replayer

import (
	"encoding/json"
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

func dumpAccounts(g *core.Genesis) error {
	var addresses []string
	accounts := make(map[string]string)
	for address, account := range g.Alloc {
		a := address.String()
		addresses = append(addresses, a)
		jsn, err := json.MarshalIndent(account, "", "  ")
		if err != nil {
			return err
		}
		accounts[a] = string(jsn)
	}
	sort.Strings(addresses)
	for _, a := range addresses {
		fmt.Printf("%s: %s\n", a, accounts[a])
	}
	return nil
}

// ProcGenesisBlock processes the genesis block for the given testnet.
func ProcGenesisBlock(testnet string) error {
	g := getGenesisBlock(testnet)
	return dumpAccounts(g)
}
