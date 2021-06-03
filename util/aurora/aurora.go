// Package aurora contains wrappers around some Aurora CLI commands.
package aurora

import (
	"os"
	"os/exec"
	"strconv"
)

// Install the EVM contract with given accountID owner and chainID.
func Install(accountID string, chainID uint8, contract string) error {
	cmd := exec.Command(
		"aurora", "install",
		"--chain", strconv.FormatUint(uint64(chainID), 10),
		"--engine", accountID,
		"--signer", accountID,
		"--owner", accountID,
		contract,
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// Upgrade the EVM contract with given accountID owner and ChainID.
func Upgrade(accountID string, chainID uint8, contract string) error {
	// `aurora upgrade` is an alias for `aurora install`
	return Install(accountID, chainID, contract)
}
