package command

import (
	"errors"
	"flag"
	"fmt"

	"github.com/aurora-is-near/evm-bully/nearapi"
)

type nodeURLFlag string

func (n *nodeURLFlag) registerFlag(fs *flag.FlagSet, cfg *nearapi.Config) {
	fs.StringVar((*string)(n), "nodeUrl", cfg.NodeURL, "NEAR node URL")
}

type testnetFlags struct {
	goerli  bool
	rinkeby bool
	ropsten bool
}

func (f *testnetFlags) registerFlags(fs *flag.FlagSet) {
	fs.BoolVar(&f.goerli, "goerli", false, "Use the Görli testnet")
	fs.BoolVar(&f.rinkeby, "rinkeby", false, "Use the Rinkeby testnet")
	fs.BoolVar(&f.ropsten, "ropsten", false, "Use the Ropsten testnet")
}

func (f *testnetFlags) determineTestnet() (chainID uint8, testnet string, err error) {
	if !f.goerli && !f.rinkeby && !f.ropsten {
		return 0, "", errors.New("one of the options -goerli, -rinkeby, or -ropsten is mandatory")
	}
	if f.goerli && f.rinkeby {
		return 0, "", errors.New("the options -goerli and -rinkeby exclude each other")
	}
	if f.goerli && f.ropsten {
		return 0, "", errors.New("the options -goerli and -ropsten exclude each other")
	}
	if f.rinkeby && f.ropsten {
		return 0, "", errors.New("the options -rinkeby and -ropsten exclude each other")
	}
	if f.rinkeby {
		return 4, "rinkeby", nil
	} else if f.ropsten {
		return 3, "ropsten", nil
	}
	// use Görli as the default
	return 5, "goerli", nil
}

func adjustBlockDefaults(block *uint64, hash *string, testnet string) {
	switch testnet {
	case "rinkeby":
		if *block == defaultGoerliBlockHeight && *hash == defaultGoerliBlockHash {
			fmt.Printf("changing -block value to %d\n", defaultRinkebyBlockHeight)
			*block = defaultRinkebyBlockHeight
			fmt.Printf("changing -hash value to %s\n", defaultRinkebyBlockHash)
			*hash = defaultRinkebyBlockHash
		}
	case "ropsten":
		if *block == defaultGoerliBlockHeight && *hash == defaultGoerliBlockHash {
			fmt.Printf("changing -block value to %d\n", defaultRopstenBlockHeight)
			*block = defaultRopstenBlockHeight
			fmt.Printf("changing -hash value to %s\n", defaultRopstenBlockHash)
			*hash = defaultRopstenBlockHash
		}

	}
}
