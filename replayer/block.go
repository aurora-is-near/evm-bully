package replayer

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type context struct {
	coinbase   common.Address // block.coinbase
	timestamp  uint64         // block.timestamp
	number     uint64         // block.number
	difficulty uint64         // block.difficulty
	gaslimit   uint64         // block.gaslimit
	hash       common.Hash    // hash = block.blockHash(blockNumber)
}

func getBlockContext(b *types.Block) (*context, error) {
	var c context
	var err error
	h := b.Header()
	c.coinbase = b.Coinbase()
	c.timestamp = b.Time()
	c.number, err = bigIntToUint64(h.Number)
	if err != nil {
		return nil, err
	}
	c.difficulty, err = bigIntToUint64(h.Difficulty)
	if err != nil {
		return nil, err
	}
	c.gaslimit = h.GasLimit
	c.hash = b.Hash()
	return &c, nil
}

func (c *context) dump() {
	fmt.Println("block context:")
	fmt.Printf("block.coinbase=%s\n", c.coinbase.String())
	fmt.Printf("block.timestamp=%d\n", c.timestamp)
	fmt.Printf("block.number=%d\n", c.number)
	fmt.Printf("block.difficulty=%d\n", c.difficulty)
	fmt.Printf("block.gaslimit=%d\n", c.gaslimit)
	fmt.Printf("block.hash=%s\n", c.hash.Hex())
}
