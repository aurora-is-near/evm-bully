// Package aurora contains wrappers around some Aurora CLI commands.
package aurora

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/log"
)

var auroraCliPath string

func init() {
	auroraCliPath = "aurora"
}

// SetAuroraCliPath sets aurora-cli path (or alias)
func SetAuroraCliPath(path string) {
	auroraCliPath = path
}

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
	log.Info(fmt.Sprintf("$ %v %v", auroraCliPath, strings.Join(args, " ")))
	cmd := exec.Command(auroraCliPath, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// Upgrade the EVM contract with given accountID owner and ChainID.
func Upgrade(accountID string, chainID uint8, contract string) error {
	// `aurora upgrade` is an alias for `aurora install`
	return Install(accountID, chainID, contract)
}
