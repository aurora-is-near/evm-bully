package replayer

import (
	"fmt"
	"math/big"
)

// bigIntToUint64 converts a big integer b to a uint64, if possible.
func bigIntToUint64(b *big.Int) (uint64, error) {
	if !b.IsUint64() {
		return 0,
			fmt.Errorf("replayer: big.Int cannot be represented as uint64: %s",
				b.String())
	}
	return b.Uint64(), nil
}
