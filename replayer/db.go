package replayer

import (
	"fmt"
	"path/filepath"

	"github.com/aurora-is-near/evm-bully/util/hashcache"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/log"
)

func openDB(
	dataDir, testnet, cacheDir string,
	blockHeight uint64,
	blockHash string,
	defrost bool,
) (ethdb.Database, []common.Hash, error) {
	dbDir := filepath.Join(dataDir, testnet, "geth", "chaindata")

	log.Info(fmt.Sprintf("opening DB in '%s'", dbDir))
	// open DB readonly
	db, err := rawdb.NewLevelDBDatabaseWithFreezer(dbDir, 0, 0,
		filepath.Join(dbDir, "ancient"), "", true)
	if err != nil {
		return nil, nil, err
	}

	// "defrost" the database first
	if defrost {
		rawdb.InitDatabaseFromFreezer(db)
	}

	// load block hash cache
	blocks, err := hashcache.Load(cacheDir)
	if err != nil {
		return nil, nil, err
	}

	if blocks == nil || uint64(len(blocks)) < blockHeight+1 || blocks[blockHeight].Hex() != blockHash {
		log.Info("cache doesn't exist, is too small, or hash mismatch")

		// read starting block
		b := rawdb.ReadBlock(db, common.HexToHash(blockHash), blockHeight)
		if b == nil {
			return nil, nil, fmt.Errorf("cannot read block at height %d with hash %s",
				blockHeight, blockHash)
		}
		log.Info(fmt.Sprintf("read block at height %d with hash %s", blockHeight,
			blockHash))

		// traverse backwards from there
		blocks, err = traverse(db, b, blockHeight)
		if err != nil {
			return nil, nil, err
		}
		blocks = append(blocks, common.HexToHash(blockHash))

		// save block hash cache
		if err := hashcache.Save(cacheDir, blocks); err != nil {
			return nil, nil, err
		}
	} else {
		log.Info("block hashes read from cache")
		// truncate blocks to blockHeight
		blocks = blocks[:blockHeight+1]
	}

	return db, blocks, nil
}
