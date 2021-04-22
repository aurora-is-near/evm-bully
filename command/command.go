// Package command implements the evm-bully commands.
package command

import (
	"github.com/ethereum/go-ethereum/node"
)

const (
	defaultBlockHeight = 4445666
	defaultBlockhash   = "0xb43dc6e5961f89b9f3c98b15d389d88dad8b0067fae76afcb5c0738aa21bd8be"
	defaultNodeURL     = "https://rpc.testnet.near.org"
)

var (
	defaultDataDir = node.DefaultDataDir()
)
