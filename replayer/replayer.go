// Package replayer implements an Ethereum transaction replayer.
package replayer

import (
	"fmt"

	"github.com/aurora-is-near/evm-bully/nearapi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/log"
)

// traverse blockchain backwards starting at block b with given blockHeight
// and return list of block hashes starting with the genesis block.
func traverse(
	db ethdb.Database,
	b *types.Block,
	blockHeight uint64,
) ([]common.Hash, error) {
	var (
		blocks  []common.Hash
		txCount uint64
	)
	for blockHeight > 0 {
		blockHash := b.ParentHash()
		blockHeight--
		b = rawdb.ReadBlock(db, blockHash, blockHeight)
		if b == nil {
			return nil, fmt.Errorf("cannot read block at height %d with hash %s",
				blockHeight, blockHash.Hex())
		}
		log.Info(fmt.Sprintf("read block at height %d with hash %s",
			blockHeight, blockHash.Hex()))
		blocks = append(blocks, blockHash)
		txCount += uint64(len(b.Transactions()))
	}
	// reverse blocks
	for i, j := 0, len(blocks)-1; i < j; i, j = i+1, j-1 {
		blocks[i], blocks[j] = blocks[j], blocks[i]
	}
	log.Info(fmt.Sprintf("total number of transactions: %d", txCount))
	return blocks, nil
}

// generateTransactions starting at genesis block.
func (r *Replayer) generateTransactions(
	a *nearapi.Account,
	evmContract string,
	db ethdb.Database,
	blocks []common.Hash,
) error {
	// process genesis block
	genesisBlock := getGenesisBlock(r.Testnet)
	err := r.beginChain(a, evmContract, genesisBlock)
	if err != nil {
		return err
	}

	for blockHeight, blockHash := range blocks {
		// read block from DB
		b := rawdb.ReadBlock(db, blockHash, uint64(blockHeight))
		if b == nil {
			return fmt.Errorf("cannot read block at height %d with hash %s",
				blockHeight, blockHash.Hex())
		}

		// block context
		c, err := getBlockContext(b)
		if err != nil {
			return err
		}
		//c.dump()
		if len(b.Transactions()) > 0 {
			if err := beginBlock(a, evmContract, r.Gas, c); err != nil {
				return err
			}
		} else {
			// TODO
			fmt.Printf("begin_block() skipped for empty block %d\n", blockHeight)
		}

		// transactions
		err = submit(a, evmContract, r.Gas, blockHeight, b.Transactions())
		if err != nil {
			return err
		}

	}
	return nil
}

// A Replayer replays transactions.
type Replayer struct {
	ChainID     uint8
	Gas         uint64
	DataDir     string
	Testnet     string
	BlockHeight uint64
	BlockHash   string
	Defrost     bool
	Skip        bool // skip empty blocks
	Batch       bool // batch transactions
	BatchSize   int  // batch size when batching transactions
}

// Replay transactions with evmContract owned by account a.
func (r *Replayer) Replay(a *nearapi.Account, evmContract string) error {
	// determine cache directory
	cacheDir, err := determineCacheDir(r.Testnet)
	if err != nil {
		return err
	}

	// open database
	db, blocks, err := openDB(r.DataDir, r.Testnet, cacheDir, r.BlockHeight,
		r.BlockHash, r.Defrost)
	if err != nil {
		return err
	}
	defer func() {
		log.Info("closing DB")
		db.Close()
	}()

	// generate transactions starting at genesis block
	err = r.generateTransactions(a, evmContract, db, blocks)
	if err != nil {
		return err
	}

	return nil
}
