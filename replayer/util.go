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

// bigIntToRawU256 converts a big integer b to a RawU256, if possible.
func bigIntToRawU256(b *big.Int) (RawU256, error) {
	var res RawU256
	bytes := b.Bytes()
	if len(bytes) > 32 {
		return res,
			fmt.Errorf("replayer: big.Int cannot be represented as RawU256: %s",
				b.String())
	}
	// convert big-endian to little-endian
	for i, j := 0, len(bytes)-1; i < j; i, j = i+1, j-1 {
		bytes[i], bytes[j] = bytes[j], bytes[i]
	}
	copy(res[:], bytes[:])
	return res, nil
}
