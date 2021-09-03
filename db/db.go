// Package db implements functions to create and read Ethereum database dumps.
package db

import (
	"encoding/gob"
	"fmt"
	"io"
	"math/big"
	"os"
	"path/filepath"

	"github.com/aurora-is-near/evm-bully/util"
	"github.com/aurora-is-near/evm-bully/util/hashcache"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/log"
	"github.com/frankbraun/codechain/util/file"
)

// Block defines an Ethereum block.
type Block struct {
	Header       *types.Header
	Coinbase     common.Address
	Time         uint64
	Hash         common.Hash
	Transactions []*Transaction
}

// Transaction defines an Ethereum transaction.
type Transaction struct {
	RLP      []byte
	Nonce    uint64
	GasPrice *big.Int
	GasLimit uint64
	To       *common.Address
	Value    *big.Int
	Data     []byte
}

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

// Open database.
func Open(
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

func readTx(tx *types.Transaction) (*Transaction, error) {
	var (
		encTx Transaction
		err   error
	)
	encTx.RLP, err = tx.MarshalBinary()
	if err != nil {
		return nil, err
	}
	encTx.Nonce = tx.Nonce()
	encTx.GasPrice = tx.GasPrice()
	encTx.GasLimit = tx.Gas()
	encTx.To = tx.To()
	encTx.Value = tx.Value()
	encTx.Data = tx.Data()
	return &encTx, nil
}

// Dump dumps the Ethereum database for the given testnet stored in dataDir
// up to blockHeight with given blockHash into the evm-bully cache directory:
//  ~/.config/evm-bully/tetstnet/dbdump.
func Dump(
	dataDir, testnet string,
	blockHeight uint64,
	blockHash string,
	defrost bool,
) error {
	// determine cache directory
	cacheDir, err := util.DetermineCacheDir(testnet)
	if err != nil {
		return err
	}

	// check dump file
	dumpFile := filepath.Join(cacheDir, "dump.db")
	exists, err := file.Exists(dumpFile)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("db: file '%s' exists already", dumpFile)
	}

	// open database
	db, blocks, err := Open(dataDir, testnet, cacheDir, blockHeight,
		blockHash, defrost)
	if err != nil {
		return err
	}
	defer func() {
		log.Info("closing DB")
		db.Close()
	}()

	// open dump file
	fp, err := os.Create(dumpFile)
	if err != nil {
		return err
	}
	defer func() {
		fp.Close()
		log.Info(fmt.Sprintf("'%s' written", dumpFile))
	}()
	/*
		gw := gzip.NewWriter(fp)
		defer func() {
			gw.Close()
			log.Info("gzip writer closed")
		}()
	*/
	enc := gob.NewEncoder(fp)

	// read DB
	for blockHeight, blockHash := range blocks {
		// read block from DB
		b := rawdb.ReadBlock(db, blockHash, uint64(blockHeight))
		if b == nil {
			return fmt.Errorf("cannot read block at height %d with hash %s",
				blockHeight, blockHash.Hex())
		}

		// transactions
		var encBlock Block
		encBlock.Header = b.Header()
		encBlock.Coinbase = b.Coinbase()
		encBlock.Time = b.Time()
		encBlock.Hash = b.Hash()
		if len(b.Transactions()) > 0 {
			for _, tx := range b.Transactions() {
				encTx, err := readTx(tx)
				if err != nil {
					return err
				}
				encBlock.Transactions = append(encBlock.Transactions, encTx)
			}
		}
		// save block
		if err := enc.Encode(encBlock); err != nil {
			return err
		}
		log.Info(fmt.Sprintf("block %d/%d written", blockHeight, len(blocks)))
	}
	return nil
}

// Reader implments a DB dump reader.
type Reader struct {
	fp  *os.File
	dec *gob.Decoder
}

// NewReader returns a new DB dump reader for the given testnet.
func NewReader(testnet string) (*Reader, error) {
	// determine cache directory
	cacheDir, err := util.DetermineCacheDir(testnet)
	if err != nil {
		return nil, err
	}

	// check dump file
	dumpFile := filepath.Join(cacheDir, "dump.db")
	exists, err := file.Exists(dumpFile)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, fmt.Errorf("db: file '%s' doesn't exist", dumpFile)
	}

	// open reader
	var r Reader
	r.fp, err = os.Open(dumpFile)
	if err != nil {
		return nil, err
	}
	/*
		gr, err := gzip.NewReader(r.fp)
		if err != nil {
			return nil, err
		}
	*/
	r.dec = gob.NewDecoder(r.fp)
	return &r, nil
}

// Next returns the next Block for the given reader or nil.
func (r *Reader) Next() (*Block, error) {
	var b Block
	if err := r.dec.Decode(&b); err != nil {
		if err == io.EOF {
			return nil, nil
		}
		return nil, err
	}
	return &b, nil
}

// Close closes the reader.
func (r *Reader) Close() error {
	return r.fp.Close()
}
