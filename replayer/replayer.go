// Package replayer implements an Ethereum transaction replayer.
package replayer

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/aurora-is-near/evm-bully/db"
	"github.com/aurora-is-near/evm-bully/replayer/neard"
	"github.com/aurora-is-near/evm-bully/util/aurora"
	"github.com/aurora-is-near/evm-bully/util/tar"
	"github.com/aurora-is-near/near-api-go"
	"github.com/aurora-is-near/near-api-go/utils"
	"github.com/ethereum/go-ethereum/log"
	"github.com/frankbraun/codechain/util/file"
)

// A Replayer replays transactions.
type Replayer struct {
	Config         *near.Config
	Timeout        time.Duration
	ChainID        uint8
	Gas            uint64
	DataDir        string
	Testnet        string
	Defrost        bool
	Skip           bool   // skip empty blocks
	Batch          bool   // batch transactions
	BatchSize      int    // batch size when batching transactions
	StartBlock     int    // start replaying at this block height
	StartTx        int    // start replaying at this transaction (in block given by StartBlock)
	Autobreak      bool   // automatically repeat with break point after error
	BreakBlock     int    // break replaying at this block height
	BreakTx        int    // break replaying at this transaction (in block given by BreakBlock)
	Release        bool   // run release version of neard
	Setup          bool   // setup and run neard before replaying
	NeardPath      string // path to neard binary
	NeardHead      string // git hash of neard
	InitialBalance string
	Contract       string
	Breakpoint     Breakpoint
}

// Breakpoint defines a break point.
type Breakpoint struct {
	ChainID          uint8  `json:"chain-id"`
	AccountID        string `json:"account-id"`
	NearcoreHead     string `json:"nearcore"`
	AuroraEngineHead string `json:"aurora-engine"`
	Transaction      string `json:"transaction"`
	tx               *db.Transaction
}

// startGenerator starts a goroutine that feeds transactions into the returned tx channel.
func (r *Replayer) startTxGenerator() chan *Tx {
	c := make(chan *Tx, 10*r.BatchSize)

	go func() {
		// process genesis block
		genesisBlock := getGenesisBlock(r.Testnet)
		c <- r.beginChainTx(genesisBlock)

		reader, err := db.NewReader(r.Testnet)
		if err != nil {
			c <- &Tx{
				BlockNum: -1,
				Error:    err,
			}
			return
		}
		defer reader.Close()

		emptyRangeStart, emptyRangeEnd := -2, -2
		flushEmptyRange := func() {
			if emptyRangeEnd < 0 {
				return
			}
			c <- &Tx{
				BlockNum: -1,
				Comment: fmt.Sprintf(
					"begin_block() skipped for empty blocks [%d;%d]",
					emptyRangeStart,
					emptyRangeEnd,
				),
			}
			emptyRangeStart, emptyRangeEnd = -2, -2
		}

	outer:
		for blockHeight := 0; true; blockHeight++ {
			b, err := reader.Next()
			if err != nil {
				flushEmptyRange()
				c <- &Tx{
					BlockNum: -1,
					Error:    err,
				}
				return
			}
			if b == nil {
				break
			}

			if blockHeight < r.StartBlock {
				c <- &Tx{
					BlockNum: -1,
					Comment:  fmt.Sprintf("skipping block %d", blockHeight),
				}
				continue
			}

			// early break, if necessary
			if r.BreakBlock != -1 && r.BreakTx == 0 && blockHeight == r.BreakBlock {
				flushEmptyRange()
				c <- &Tx{
					BlockNum: -1,
					Comment:  fmt.Sprintf("breaking block %d", blockHeight),
				}
				log.Info("sleep")
				time.Sleep(5 * time.Second)
				txs := b.Transactions
				if txs != nil && len(txs) > 0 {
					r.Breakpoint.tx = txs[0]
				}
				break
			}

			// block context
			ctx, err := getBlockContext(b)
			if err != nil {
				flushEmptyRange()
				c <- &Tx{
					BlockNum: -1,
					Error:    err,
				}
				return
			}

			if len(b.Transactions) == 0 && r.Skip {
				if emptyRangeEnd != blockHeight-1 {
					emptyRangeStart = blockHeight
				}
				emptyRangeEnd = blockHeight
				continue
			}

			flushEmptyRange()
			c <- beginBlockTx(r.Gas, ctx)

			// actual transactions
			for i, tx := range b.Transactions {
				// early break, if necessary
				if r.BreakBlock != -1 && blockHeight == r.BreakBlock && i == r.BreakTx {
					c <- &Tx{
						BlockNum: -1,
						Comment:  fmt.Sprintf("breaking at transaction %d (in block %d)", i, blockHeight),
					}
					log.Info("sleep")
					time.Sleep(5 * time.Second)
					r.Breakpoint.tx = tx
					break outer
				}
				if blockHeight == r.StartBlock && i < r.StartTx {
					c <- &Tx{
						BlockNum: -1,
						Comment:  fmt.Sprintf("skipping transaction %d (in block %d)", i, blockHeight),
					}
					continue
				}
				amount, err := utils.FormatNearAmount(strconv.FormatUint(r.Gas/uint64(r.BatchSize), 10))
				if err != nil {
					c <- &Tx{
						BlockNum: -1,
						Error:    err,
					}
					return
				}
				c <- &Tx{
					BlockNum: blockHeight,
					TxNum:    i,
					Comment: fmt.Sprintf("submit(%d, tx=%d, tx_size=%d, gas=%sⓃ)",
						blockHeight, i, len(tx.RLP), amount),
					MethodName: "submit",
					Args:       tx.RLP,
					EthTx:      tx,
				}
			}
		}
		flushEmptyRange()
		close(c)
	}()

	return c
}

func (r *Replayer) replay(
	evmContract string,
) (blockNum int, txNum int, errormsg []byte, err error) {
	// setup, if necessary
	if r.Setup {
		// setup neard
		log.Info("setup neard")

		var nearDaemon *neard.NEARDaemon
		var err error
		if r.NeardPath != "" {
			nearDaemon, err = neard.LoadFromBinary(r.NeardPath, r.NeardHead)
		} else {
			nearDaemon, err = neard.LoadFromRepo(filepath.Join("..", "nearcore"), r.Release, true)
		}
		if err != nil {
			return -1, -1, nil, err
		}

		if err := nearDaemon.SetupLocalData(); err != nil {
			return -1, -1, nil, err
		}
		if err := nearDaemon.Start(); err != nil {
			return -1, -1, nil, err
		}
		defer nearDaemon.Stop()
		r.Breakpoint.NearcoreHead = nearDaemon.Head

		log.Info("sleep")
		time.Sleep(5 * time.Second)

		// create account
		log.Info("create account")
		ca := CreateAccount{
			Config:         r.Config,
			InitialBalance: r.InitialBalance,
			MasterAccount:  strings.Join(strings.Split(r.Breakpoint.AccountID, ".")[1:], "."),
		}
		if err := ca.Create(r.Breakpoint.AccountID); err != nil {
			return -1, -1, nil, err
		}

		// install EVM contract
		log.Info("install EVM contract")
		err = aurora.Install(r.Breakpoint.AccountID, r.ChainID, r.Contract)
		if err != nil {
			return -1, -1, nil, err
		}

		// reset key path
		r.Config.KeyPath = ""
	}

	// load account
	conn := near.NewConnectionWithTimeout(r.Config.NodeURL, r.Timeout)
	a, err := near.LoadAccount(conn, r.Config, r.Breakpoint.AccountID)
	if err != nil {
		return -1, -1, nil, err
	}

	// process transactions
	batch := make([]near.Action, 0, r.BatchSize)
	zeroAmount := big.NewInt(0)
	c := r.startTxGenerator()
	for tx := range c {
		if tx.Error != nil {
			return -1, -1, nil, tx.Error
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
					return -1, -1, nil, err
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
						return -1, -1, nil, err
					}
					batch = batch[:0] // reset
				} else {
					continue // batch no full yet
				}
			}
			if errormsg, err := procTxResult(r.Batch, tx.EthTx, txResult); err != nil {
				return tx.BlockNum, tx.TxNum, errormsg, err
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
			return -1, -1, nil, err
		}
		if errormsg, err := procTxResult(r.Batch, nil, txResult); err != nil {
			return -1, -1, errormsg, err
		}
	}
	return -1, -1, nil, nil
}

// Replay transactions with evmContract.
func (r *Replayer) Replay(evmContract string) error {
	keyPath := r.Config.KeyPath
	blockNum, txNum, errormsg, err := r.replay(evmContract)
	if err != nil {
		if r.Autobreak && blockNum != -1 {
			// restore key path
			r.Config.KeyPath = keyPath
			// replay again with breakpoint set
			r.BreakBlock = blockNum
			r.BreakTx = txNum
			log.Info("replay again with breakpoint set")
			if _, _, _, err := r.replay(evmContract); err != nil {
				return err
			}
		} else {
			return err
		}
	}

	if r.BreakBlock != -1 {
		return r.saveBreakpoint(errormsg)
	}
	return nil
}

func showTx(tx *db.Transaction) {
	fmt.Println("transaction:")
	fmt.Println("0x" + hex.EncodeToString(tx.RLP))
	fmt.Printf("nonce: %d\n", tx.Nonce)
	fmt.Printf("gasPrice: %s\n", tx.GasPrice.String())
	fmt.Printf("gasLimit: %d\n", tx.GasLimit)
	if tx.To != nil {
		fmt.Printf("to: 0x%s\n", hex.EncodeToString(tx.To[:]))
	} else {
		fmt.Println("to: contract creation")
	}
	fmt.Printf("value: %s\n", tx.Value.String())
	if len(tx.Data) > 0 {
		fmt.Println("data:")
		fmt.Println("0x" + hex.EncodeToString(tx.Data))
	}
}

func procTxResult(
	batch bool,
	tx *db.Transaction,
	txResult map[string]interface{},
) ([]byte, error) {
	utils.PrettyPrintResponse(txResult)
	status := txResult["status"].(map[string]interface{})
	jsn, err := json.MarshalIndent(status, "", "  ")
	if err != nil {
		return nil, err
	}
	fmt.Println(string(jsn))
	if status["Failure"] != nil {
		if !batch && tx != nil {
			// print last failing transaction if possible
			showTx(tx)
		}
		return jsn, errors.New("replayer: transaction failed")
	}
	return nil, nil
}

// saveBreakpoint saves replayer break point for evmContract.
func (r *Replayer) saveBreakpoint(errormsg []byte) error {
	var err error
	dir := fmt.Sprintf("%s-block-%d-tx-%d", r.Testnet, r.BreakBlock, r.BreakTx)
	log.Info(fmt.Sprintf("save breakpoint %s", dir))

	// set chainID
	r.Breakpoint.ChainID = r.ChainID

	// get HEAD of aurora-engine
	r.Breakpoint.AuroraEngineHead, err = auroraEngineHead(r.Contract)
	if err != nil {
		return err
	}

	// encode transaction
	if r.Breakpoint.tx != nil {
		r.Breakpoint.Transaction = hex.EncodeToString(r.Breakpoint.tx.RLP)
	}

	// remove output dir
	if err := os.RemoveAll(dir); err != nil {
		return err
	}
	if err := os.Mkdir(dir, 0755); err != nil {
		return err
	}

	// marshal breakpoint data structure
	jsn, err := json.MarshalIndent(&r.Breakpoint, "", "  ")
	if err != nil {
		return err
	}
	filename := filepath.Join(dir, "breakpoint.json")
	if err := os.WriteFile(filename, jsn, 0644); err != nil {
		return err
	}
	log.Info(fmt.Sprintf("'%s' written", filename))

	// write error message, if defined (-autobreak was used)
	if errormsg != nil {
		filename := filepath.Join(dir, "errormsg.json")
		if err := os.WriteFile(filename, errormsg, 0644); err != nil {
			return err
		}
		log.Info(fmt.Sprintf("'%s' written", filename))
	}

	// copy key file
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	filename = r.Breakpoint.AccountID + ".json"
	path := filepath.Join(home, ".near-credentials", r.Config.NetworkID, filename)
	if err := file.Copy(path, filepath.Join(dir, filename)); err != nil {
		return err
	}

	// copy local nearcore directory
	localDir := filepath.Join(home, ".near", "local")
	if err := file.CopyDir(localDir, filepath.Join(dir, "local")); err != nil {
		return err
	}

	// tar everything up
	return tar.Create(dir)
}
