// Package aurora contains wrappers around some Aurora CLI commands.
package aurora

import (
	"os"
	"os/exec"
)

// Install the EVM contract with given accountID owner.
func Install(accountID, contract string) error {
	cmd := exec.Command(
		"aurora", "install",
		"--chain", "1313161556",
		"--engine", accountID,
		"--signer", accountID,
		"--owner", accountID,
		contract,
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// Upgrade the EVM contract with given accountID owner.
func Upgrade(accountID, contract string) error {
	// `aurora upgrade` is an alias for `aurora install`
	return Install(accountID, contract)
}
