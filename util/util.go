package util

import (
	"os"
	"path/filepath"

	"github.com/frankbraun/codechain/util/homedir"
)

// DetermineCacheDir determines the evm-bully cache directory:
//  ~/.config/evm-bully/tetstnet
func DetermineCacheDir(testnet string) (string, error) {
	homeDir := homedir.Get("evm-bully")
	cacheDir := filepath.Join(homeDir, testnet)
	// make sure cache directory exists
	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		return "", err
	}
	return cacheDir, nil
}
