// Package replayer implements an Ethereum transaction replayer.
package replayer

import (
	"context"
	"encoding/hex"
	"fmt"
	"path/filepath"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/rpc"
)

// traverse blockchain backwards starting at block b with given blockHeight
// and return list of block hashes starting with the genesis block.
func traverse(
	db ethdb.Database,
	b *types.Block,
	blockHeight uint64,
) ([]common.Hash, error) {
	var blocks []common.Hash
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
	}
	// reverse blocks
	for i, j := 0, len(blocks)-1; i < j; i, j = i+1, j-1 {
		blocks[i], blocks[j] = blocks[j], blocks[i]
	}

	return blocks, nil
}

// generateTransactions starting at genesis block.
func generateTransactions(
	ctx context.Context,
	endpoint string,
	db ethdb.Database,
	blocks []common.Hash,
) error {
	c, err := rpc.DialContext(ctx, endpoint)
	if err != nil {
		return err
	}
	ec := ethclient.NewClient(c)
	defer ec.Close()

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
		c.dump()

		// transactions
		for i, tx := range b.Transactions() {
			fmt.Printf("b=%d tx=%d chainid=%s data=%s\n", blockHeight, i,
				tx.ChainId().String(), hex.EncodeToString(tx.Data()))

			// submit transaction to JSON-RPC endpoint
			if err := ec.SendTransaction(ctx, tx); err != nil {
				return err
			}
		}

	}
	return nil
}

// ReadTxs reads transactions from datadir, starting at block with given
// blockHeight and blockHash.
func ReadTxs(
	ctx context.Context,
	endpoint, datadir, testnet string,
	blockHeight uint64,
	blockHash string,
) error {
	dbDir := filepath.Join(datadir, testnet, "geth", "chaindata")

	log.Info(fmt.Sprintf("opening DB in '%s'", dbDir))
	db, err := rawdb.NewLevelDBDatabaseWithFreezer(dbDir, 0, 0, filepath.Join(dbDir, "ancient"), "")
	if err != nil {
		return err
	}
	defer func() {
		log.Info(fmt.Sprintf("closing DB in '%s'", dbDir))
		db.Close()
	}()

	// TODO: we might have to "defrost" the database first in some cases
	// rawdb.InitDatabaseFromFreezer(db)

	// read starting block
	b := rawdb.ReadBlock(db, common.HexToHash(blockHash), blockHeight)
	if b == nil {
		return fmt.Errorf("cannot read block at height %d with hash %s",
			blockHeight, blockHash)
	}
	log.Info(fmt.Sprintf("read block at height %d with hash %s", blockHeight,
		blockHash))

	// traverse backwards from there
	blocks, err := traverse(db, b, blockHeight)
	if err != nil {
		return err
	}
	blocks = append(blocks, common.HexToHash(blockHash))

	// generate transactions starting at genesis block
	if err := generateTransactions(ctx, endpoint, db, blocks); err != nil {
		return err
	}

	return nil
}
