// Package aurora contains wrappers around some Aurora CLI commands.
package aurora

import (
	"os"
	"os/exec"
)

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
