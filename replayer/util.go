package replayer

import (
	"fmt"
	"math/big"
	"path/filepath"

	"github.com/aurora-is-near/evm-bully/util/git"
	"github.com/ethereum/go-ethereum/log"
)

func auroraEngineHead(contract string) (string, error) {
	// get current HEAD in aurora-engine directory
	head, err := git.Head(filepath.Dir(contract))
	if err != nil {
		return "", err
	}
	log.Info(fmt.Sprintf("head=%s", head))
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
