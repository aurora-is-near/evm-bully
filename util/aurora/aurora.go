// Package aurora contains wrappers around some Aurora CLI commands.
package aurora

import (
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/log"
)

// Install the EVM contract with given accountID owner and chainID.
func Install(accountID string, chainID uint8, contract string) error {
	args := []string{
		"install",
		"--chain", strconv.FormatUint(uint64(chainID), 10),
		"--engine", accountID,
		"--signer", accountID,
		"--owner", accountID,
		contract,
	}
	log.Info("$ aurora " + strings.Join(args, " "))
	cmd := exec.Command("aurora", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// Upgrade the EVM contract with given accountID owner and ChainID.
func Upgrade(accountID string, chainID uint8, contract string) error {
	// `aurora upgrade` is an alias for `aurora install`
	return Install(accountID, chainID, contract)
}
