// Package replayer implements an Ethereum transaction replayer.
package replayer

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/big"

	"github.com/aurora-is-near/evm-bully/nearapi"
	"github.com/aurora-is-near/evm-bully/nearapi/utils"
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

// startGenerator starts a goroutine that feeds transactions into the returned tx channel.
func (r *Replayer) startTxGenerator(
	a *nearapi.Account,
	evmContract string,
	db ethdb.Database,
	blocks []common.Hash,
) chan *Tx {
	c := make(chan *Tx, 10*r.BatchSize)

	go func() {
		// process genesis block
		genesisBlock := getGenesisBlock(r.Testnet)
		c <- r.beginChainTx(a, evmContract, genesisBlock)

		for blockHeight, blockHash := range blocks {
			// read block from DB
			b := rawdb.ReadBlock(db, blockHash, uint64(blockHeight))
			if b == nil {
				c <- &Tx{Error: fmt.Errorf("cannot read block at height %d with hash %s",
					blockHeight, blockHash.Hex())}
				return
			}

			// block context
			ctx, err := getBlockContext(b)
			if err != nil {
				c <- &Tx{Error: err}
				return
			}
			if !r.Skip || len(b.Transactions()) > 0 {
				c <- beginBlockTx(a, evmContract, r.Gas, ctx)
			} else {
				c <- &Tx{Comment: fmt.Sprintf("begin_block() skipped for empty block %d", blockHeight)}
			}

			// actual transactions
			for i, tx := range b.Transactions() {
				// get signed transaction in RLP encoding
				rlp, err := tx.MarshalBinary()
				if err != nil {
					c <- &Tx{Error: err}
					return
				}
				c <- &Tx{
					Comment:    fmt.Sprintf("submit(%d, tx=%d, tx_size=%d)", blockHeight, i, len(rlp)),
					MethodName: "submit",
					Args:       rlp,
				}
			}
		}
		close(c)
	}()

	return c
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

	// process transactions
	batch := make([]nearapi.Action, 0, r.BatchSize)
	zeroAmount := big.NewInt(0)
	c := r.startTxGenerator(a, evmContract, db, blocks)
	for tx := range c {
		if tx.Error != nil {
			return tx.Error
		}
		if tx.MethodName != "" {
			var (
				txResult map[string]interface{}
				err      error
			)
			if !r.Batch {
				// no tx batching
				if tx.Comment != "" {
					fmt.Println(tx.Comment)
				}
				txResult, err = a.FunctionCall(evmContract, tx.MethodName, tx.Args, r.Gas, *zeroAmount)
				if err != nil {
					return err
				}
			} else {
				// batch mode
				if tx.Comment != "" {
					fmt.Println("batching: " + tx.Comment)
				}
				batch = append(batch, nearapi.Action{
					Enum: 2,
					FunctionCall: nearapi.FunctionCall{
						MethodName: tx.MethodName,
						Args:       tx.Args,
						Gas:        r.Gas / uint64(r.BatchSize),
						Deposit:    *zeroAmount,
					},
				})
				if len(batch) == r.BatchSize {
					fmt.Println("running batch")
					txResult, err = a.SignAndSendTransaction(evmContract, batch)
					if err != nil {
						return err
					}
					batch = batch[:0] // reset
				} else {
					continue
				}
			}

			utils.PrettyPrintResponse(txResult)
			status := txResult["status"].(map[string]interface{})
			jsn, err := json.MarshalIndent(status, "", "  ")
			if err != nil {
				return err
			}
			fmt.Println(string(jsn))
			if status["Failure"] != nil {
				return errors.New("replayer: transaction failed")
			}
		} else if tx.Comment != "" {
			fmt.Println(tx.Comment)
		}
	}

	// process last batch, if not empty
	if len(batch) > 0 {
		fmt.Println("running last batch")
		txResult, err := a.SignAndSendTransaction(evmContract, batch)
		if err != nil {
			return err
		}
		utils.PrettyPrintResponse(txResult)
		status := txResult["status"].(map[string]interface{})
		jsn, err := json.MarshalIndent(status, "", "  ")
		if err != nil {
			return err
		}
		fmt.Println(string(jsn))
		if status["Failure"] != nil {
			return errors.New("replayer: transaction failed")
		}
	}
	return nil
}
