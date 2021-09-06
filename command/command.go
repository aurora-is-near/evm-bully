// Package command implements the evm-bully commands.
package command

import (
	"github.com/ethereum/go-ethereum/node"
)

const (
	// 2021-05-07
	defaultGoerliBlockHeight  = 4747554
	defaultGoerliBlockHash    = "0xca3f0a8bcbfadf60994423da4009b9519ccab6e1e91c637888d313ecf24f0a1a"
	defaultRinkebyBlockHeight = 8541193
	defaultRinkebyBlockHash   = "0x9afd56145c5a771967b0d86800338694214e9c83d9d89d52215d916b555d9cd5"
	defaultRopstenBlockHeight = 10187164
	defaultRopstenBlockHash   = "0x000f6b7fc929f6f3a493cad3cee9d65274e51228a05970cf03ad7fb6664007a6"
)

const (
	defaultGas            = 800000000000000
	defaultInitialBalance = "100"
)

var (
	defaultDataDir = node.DefaultDataDir()
)
