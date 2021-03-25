// Package replayer implements an Ethereum transaction replayer.
package replayer

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/aurora-is-near/evm-bully/util/hashcache"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
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
		//c.dump()
		if err := beginBlock(c); err != nil {
			return err
		}

		// transactions
		/*
			for i, tx := range b.Transactions() {
				fmt.Printf("b=%d tx=%d chainid=%s data=%s\n", blockHeight, i,
					tx.ChainId().String(), hex.EncodeToString(tx.Data()))

				// submit transaction to JSON-RPC endpoint ("eth_sendRawTransaction")
				if err := ec.SendTransaction(ctx, tx); err != nil {
					return err
				}
			}
		*/
		if err := rawCall(blockHeight, b.Transactions()); err != nil {

		}

	}
	return nil
}

// Replay transactions from dataDir up block with given blockHeight and
// blockHash.
func Replay(
	ctx context.Context,
	endpoint, dataDir, testnet, cacheDir string,
	blockHeight uint64,
	blockHash string,
	defrost bool,
) error {
	dbDir := filepath.Join(dataDir, testnet, "geth", "chaindata")

	log.Info(fmt.Sprintf("opening DB in '%s'", dbDir))
	db, err := rawdb.NewLevelDBDatabaseWithFreezer(dbDir, 0, 0, filepath.Join(dbDir, "ancient"), "")
	if err != nil {
		return err
	}
	defer func() {
		log.Info(fmt.Sprintf("closing DB in '%s'", dbDir))
		db.Close()
	}()

	// "defrost" the database first
	if defrost {
		rawdb.InitDatabaseFromFreezer(db)
	}

	// load block hash cache
	blocks, err := hashcache.Load(cacheDir)
	if err != nil {
		return err
	}

	if blocks == nil || uint64(len(blocks)) < blockHeight+1 || blocks[blockHeight].Hex() != blockHash {
		log.Info("cache doesn't exist, is too small, or hash mismatch")

		// read starting block
		b := rawdb.ReadBlock(db, common.HexToHash(blockHash), blockHeight)
		if b == nil {
			return fmt.Errorf("cannot read block at height %d with hash %s",
				blockHeight, blockHash)
		}
		log.Info(fmt.Sprintf("read block at height %d with hash %s", blockHeight,
			blockHash))

		// traverse backwards from there
		blocks, err = traverse(db, b, blockHeight)
		if err != nil {
			return err
		}
		blocks = append(blocks, common.HexToHash(blockHash))

		// save block hash cache
		if err := hashcache.Save(cacheDir, blocks); err != nil {
			return err
		}
	} else {
		log.Info("block hashes read from cache")
		// truncate blocks to blockHeight
		blocks = blocks[:blockHeight+1]
	}

	// process genesis block
	g := getGenesisBlock(testnet)
	if err := beginChain(g); err != nil {
		return err
	}

	// generate transactions starting at genesis block
	if err := generateTransactions(ctx, endpoint, db, blocks); err != nil {
		return err
	}

	return nil
}

func beginChain(g *core.Genesis) error {
	fmt.Println("begin_chain()")
	return nil
}

func beginBlock(c *blockContext) error {
	fmt.Printf("begin_block(%d)\n", c.number)
	return nil
}

func rawCall(blockHeight int, txs types.Transactions) error {
	// TODO: batching
	for i, _ := range txs {
		fmt.Printf("raw_call(%d, %d)\n", blockHeight, i)
	}
	return nil
}
