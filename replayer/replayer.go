// Package replayer implements an Ethereum transaction replayer.
package replayer

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"strconv"

	"github.com/aurora-is-near/near-api-go"
	"github.com/aurora-is-near/near-api-go/utils"
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
	a *near.Account,
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
			if blockHeight < r.StartBlock {
				c <- &Tx{Comment: fmt.Sprintf("skipping block %d", blockHeight)}
				continue
			}
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
				amount, err := utils.FormatNearAmount(strconv.FormatUint(r.Gas/uint64(r.BatchSize), 10))
				if err != nil {
					c <- &Tx{Error: err}
					return
				}
				c <- &Tx{
					Comment: fmt.Sprintf("submit(%d, tx=%d, tx_size=%d, gas=%sⓃ)",
						blockHeight, i, len(rlp), amount),
					MethodName: "submit",
					Args:       rlp,
					EthTx:      tx,
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
	StartBlock  int  // start replaying at this block height
}

// Replay transactions with evmContract owned by account a.
func (r *Replayer) Replay(a *near.Account, evmContract string) error {
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
	batch := make([]near.Action, 0, r.BatchSize)
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
				batch = append(batch, near.Action{
					Enum: 2,
					FunctionCall: near.FunctionCall{
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
					continue // batch no full yet
				}
			}
			if err := procTxResult(r.Batch, tx.EthTx, txResult); err != nil {
				return err
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
		if err := procTxResult(r.Batch, nil, txResult); err != nil {
			return err
		}
	}
	return nil
}

func showTx(tx *types.Transaction) {
	rlp, err := tx.MarshalBinary()
	if err != nil {
		panic(err) // marshalled before
	}
	fmt.Println("transaction:")
	fmt.Println("0x" + hex.EncodeToString(rlp))
	fmt.Printf("nonce: %d\n", tx.Nonce())
	fmt.Printf("gasPrice: %s\n", tx.GasPrice().String())
	fmt.Printf("gasLimit: %d\n", tx.Gas())
	if tx.To() != nil {
		fmt.Printf("to: 0x%s\n", hex.EncodeToString(tx.To()[:]))
	} else {
		fmt.Println("to: contract creation")
	}
	fmt.Printf("value: %s\n", tx.Value().String())
	if len(tx.Data()) > 0 {
		fmt.Println("data:")
		fmt.Println("0x" + hex.EncodeToString(tx.Data()))
	}
}

func procTxResult(batch bool, tx *types.Transaction, txResult map[string]interface{}) error {
	utils.PrettyPrintResponse(txResult)
	status := txResult["status"].(map[string]interface{})
	jsn, err := json.MarshalIndent(status, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(jsn))
	if status["Failure"] != nil {
		if !batch && tx != nil {
			// print last failing transaction if possible
			showTx(tx)
		}
		return errors.New("replayer: transaction failed")
	}
	return nil
}
