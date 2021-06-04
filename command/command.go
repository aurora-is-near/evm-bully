// Package command implements the evm-bully commands.
package command

import (
	"github.com/ethereum/go-ethereum/node"
)

const (
	// 2021-05-07
	defaultGoerliBlockHeight = 4747554
	defaultGoerliBlockHash   = "0xca3f0a8bcbfadf60994423da4009b9519ccab6e1e91c637888d313ecf24f0a1a"
	//defaultRinkebyBlockHeight = 8541193
	//defaultRinkebyBlockHash   = "0x9afd56145c5a771967b0d86800338694214e9c83d9d89d52215d916b555d9cd5"
	defaultRinkebyBlockHeight = 100
	defaultRinkebyBlockHash   = "0xb9dc3149942ec58aeb2832691e6dcc3d7c5d682e537239c5689c01a131c6a575"
	//defaultRopstenBlockHeight = 10187164
	//defaultRopstenBlockHash   = "0x000f6b7fc929f6f3a493cad3cee9d65274e51228a05970cf03ad7fb6664007a6"
	defaultRopstenBlockHeight = 100
	defaultRopstenBlockHash   = "0xb40a0dfde1b270d7c58c3cb505c7e773c50198b28cce3e442c4e2f33ff764582"
)

const (
	defaultGas            = 800000000000000
	defaultInitialBalance = "100"
)

var (
	defaultDataDir = node.DefaultDataDir()
)
