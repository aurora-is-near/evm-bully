package replayer

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/log"
)

func calcStatsForTxs(blockHeight int, txs types.Transactions) (uint64, error) {
	var contractTxCounter uint64
	for i, tx := range txs {
		// only look at contract creating transactions
		if tx.To() == nil {
			log.Info(fmt.Sprintf("block=%d, tx=%d", blockHeight, i))
			contractTxCounter++
		}
	}
	return contractTxCounter, nil
}

func calcStatsForBlocks(db ethdb.Database, blocks []common.Hash) error {
	var contractTxCounter uint64
	var totalTxCounter uint64
	for blockHeight, blockHash := range blocks {
		// read block from DB
		b := rawdb.ReadBlock(db, blockHash, uint64(blockHeight))
		if b == nil {
			return fmt.Errorf("cannot read block at height %d with hash %s",
				blockHeight, blockHash.Hex())
		}

		// transactions
		if len(b.Transactions()) > 0 {
			counter, err := calcStatsForTxs(blockHeight, b.Transactions())
			if err != nil {
				return err
			}
			contractTxCounter += counter
			totalTxCounter += uint64(len(b.Transactions()))
		}
	}
	fmt.Printf("contract creating txs: %d\n", contractTxCounter)
	fmt.Printf("total txs: %d\n", totalTxCounter)
	fmt.Printf("percentage of contract creating txs: %.2f%%\n", float64(contractTxCounter)/float64(totalTxCounter)*100.0)
	return nil
}

// CalcStats calculates some statistics for the given testnet and prints them to stdout.
func CalcStats(
	dataDir, testnet string,
	blockHeight uint64,
	blockHash string,
	defrost bool,
) error {
	// determine cache directory
	cacheDir, err := determineCacheDir(testnet)
	if err != nil {
		return err
	}
	// open database
	db, blocks, err := openDB(dataDir, testnet, cacheDir, blockHeight,
		blockHash, defrost)
	if err != nil {
		return err
	}
	defer func() {
		log.Info("closing DB")
		db.Close()
	}()
	// calculate statistics
	return calcStatsForBlocks(db, blocks)
}
