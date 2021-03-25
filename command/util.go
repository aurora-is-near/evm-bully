package command

import (
	"errors"
	"flag"
)

type nodeURLFlag string

func (n *nodeURLFlag) registerFlag(fs *flag.FlagSet) {
	fs.StringVar((*string)(n), "nodeUrl", defaultNodeURL, "NEAR node URL")
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

func (f *testnetFlags) determineTestnet() (string, error) {
	if !f.goerli && !f.rinkeby && !f.ropsten {
		return "", errors.New("one of the options -goerli, -rinkeby, or -ropsten is mandatory")
	}
	if f.goerli && f.rinkeby {
		return "", errors.New("the options -goerli and -rinkeby exclude each other")
	}
	if f.goerli && f.ropsten {
		return "", errors.New("the options -goerli and -ropsten exclude each other")
	}
	if f.rinkeby && f.ropsten {
		return "", errors.New("the options -rinkeby and -ropsten exclude each other")
	}
	if f.rinkeby {
		return "rinkeby", nil
	} else if f.ropsten {
		return "ropsten", nil
	}
	// use Görli as the default
	return "goerli", nil
}
