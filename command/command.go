// Package command implements the evm-bully commands.
package command

import (
	"github.com/ethereum/go-ethereum/node"
)

const (
	defaultEndpoint    = "http://localhost:8545"
	defaultBlockHeight = 100000
	defaultBlockhash   = "0x94b9c7be22783a3ee1e1f1f31e35f7de4612905d6fd24d3fe90a26091dca43fe"
)

var (
	defaultDataDir = node.DefaultDataDir()
)
