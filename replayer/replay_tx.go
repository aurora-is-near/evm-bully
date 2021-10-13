package replayer

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"time"

	"github.com/aurora-is-near/evm-bully/replayer/neard"
	"github.com/aurora-is-near/evm-bully/util/aurora"
	"github.com/aurora-is-near/evm-bully/util/git"
	"github.com/aurora-is-near/evm-bully/util/gnumake"
	"github.com/aurora-is-near/near-api-go"
	"github.com/aurora-is-near/near-api-go/utils"
	"github.com/frankbraun/codechain/util/file"
)

func buildAuroraEngine(head string) error {
	fmt.Println("build aurora-engine")
	engineDir := filepath.Join("..", "aurora-engine")
	// checkout
	if err := git.Checkout(engineDir, head); err != nil {
		return err
	}
	// build
	if err := gnumake.Make(engineDir, "evm-bully=yes"); err != nil {
		return err
	}
	return nil
}

// ReplayTx replays transaction from breakpointDir.
func ReplayTx(
	breakpointDir string,
	build bool,
	contract string,
	release bool,
	gas uint64,
) error {
	// parse breakpoint.json
	filename := filepath.Join(breakpointDir, "breakpoint.json")
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	var bp Breakpoint
	if err := json.Unmarshal(data, &bp); err != nil {
		return err
	}

	if build {
		if err := buildAuroraEngine(bp.AuroraEngineHead); err != nil {
			return err
		}
	}

	nearDir := filepath.Join("..", "nearcore")
	if build {
		if err := git.Checkout(nearDir, bp.NearcoreHead); err != nil {
			return err
		}
	}
	nearDaemon, err := neard.LoadFromRepo(nearDir, release, build)
	if err != nil {
		return err
	}
	if err := nearDaemon.RestoreLocalData(filepath.Join(breakpointDir, "local")); err != nil {
		return err
	}
	if err := nearDaemon.Start(); err != nil {
		return err
	}
	defer nearDaemon.Stop()

	if err := os.Setenv("NEAR_ENV", "local"); err != nil {
		return err
	}
	cfg := near.GetConfig()
	c := near.NewConnection(cfg.NodeURL)
	nearStarted := checkUntilTrue(time.Second*100, "near node is not up yet", func() bool {
		_, err := c.GetNodeStatus()
		return err == nil
	})
	if !nearStarted {
		return fmt.Errorf("replayer: near node is not reachable after 100 seconds")
	}

	// copy credentials file
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	credDir := filepath.Join(home, ".near-credentials", "local")
	if err := os.MkdirAll(credDir, 0700); err != nil {
		return err
	}
	dst := filepath.Join(credDir, bp.AccountID+".json")
	if err := os.RemoveAll(dst); err != nil {
		return err
	}
	err = file.Copy(filepath.Join(breakpointDir, bp.AccountID+".json"), dst)
	if err != nil {
		return err
	}

	// upgrade contract before replaying tx, if necessary
	if contract != "" {
		err := aurora.Upgrade(bp.AccountID, bp.ChainID, contract)
		if err != nil {
			return err
		}
	}

	// run transaction
	zeroAmount := big.NewInt(0)
	rlp, err := hex.DecodeString(bp.Transaction)
	if err != nil {
		return err
	}

	// TODO: why validator_key.json and test.near here?
	cfg.KeyPath = filepath.Join(home, ".near", "local", "validator_key.json")
	a, err := near.LoadAccount(c, cfg, "test.near")
	if err != nil {
		return err
	}

	txResult, err := a.FunctionCall(bp.AccountID, "submit", rlp, gas, *zeroAmount)
	if err != nil {
		return err
	}
	utils.PrettyPrintResponse(txResult)
	res, err := near.GetTransactionLastResult(txResult)
	if err != nil {
		return err
	}
	if res != nil {
		jsn, err := json.MarshalIndent(res, "", "  ")
		if err != nil {
			return err
		}
		fmt.Println(string(jsn))
	}
	return nil
}
