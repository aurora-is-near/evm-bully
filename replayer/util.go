package replayer

import (
	"fmt"
	"math/big"
	"os"
	"path/filepath"

	"github.com/aurora-is-near/evm-bully/util/git"
	"github.com/ethereum/go-ethereum/log"
	"github.com/frankbraun/codechain/util/homedir"
)

func auroraEngineHead(contract string) (string, error) {
	// get cwd
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	fmt.Println(cwd)
	// switch to aurora-engine directory
	if err := os.Chdir(filepath.Dir(contract)); err != nil {
		return "", err
	}
	// get current HEAD
	head, err := git.Head()
	if err != nil {
		return "", err
	}
	log.Info(fmt.Sprintf("head=%s", head))
	// switch back to original directory
	if err := os.Chdir(cwd); err != nil {
		return "", err
	}
	return head, nil
}

// bigIntToUint64 converts a big integer b to a uint64, if possible.
func bigIntToUint64(b *big.Int) (uint64, error) {
	if !b.IsUint64() {
		return 0,
			fmt.Errorf("replayer: big.Int cannot be represented as uint64: %s",
				b.String())
	}
	return b.Uint64(), nil
}

// bigIntToRawU256 converts a big integer b to a RawU256, if possible.
func bigIntToRawU256(b *big.Int) (RawU256, error) {
	var res RawU256
	bytes := b.Bytes()
	if len(bytes) > 32 {
		return res,
			fmt.Errorf("replayer: big.Int cannot be represented as RawU256: %s",
				b.String())
	}
	// the encoding is already big-endian
	copy(res[:], bytes[:])
	return res, nil
}

func determineCacheDir(testnet string) (string, error) {
	homeDir := homedir.Get("evm-bully")
	cacheDir := filepath.Join(homeDir, testnet)
	// make sure cache directory exists
	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		return "", err
	}
	return cacheDir, nil
}
