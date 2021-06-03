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
	"github.com/ethereum/go-ethereum/log"
	"github.com/frankbraun/codechain/util/file"
)

func buildAuroraEngine(head string) error {
	fmt.Println("build aurora-engine")
	// get cwd
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	// switch to aurora-engine directory
	nearDir := filepath.Join(cwd, "..", "aurora-engine")
	if err := os.Chdir(nearDir); err != nil {
		return err
	}
	// checkout
	if err := git.Checkout(head); err != nil {
		return err
	}
	// build
	if err := gnumake.Make("evm-bully=yes"); err != nil {
		return err
	}
	// switch back to original directory
	if err := os.Chdir(cwd); err != nil {
		return err
	}
	return nil
}

func buildNearcore(head string, release bool) error {
	fmt.Println("build nearcore")
	// get cwd
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	// switch to nearcore directory
	nearDir := filepath.Join(cwd, "..", "nearcore")
	if err := os.Chdir(nearDir); err != nil {
		return err
	}
	// checkout
	if err := git.Checkout(head); err != nil {
		return err
	}
	// build
	if err := neard.Build(release); err != nil {
		return err
	}
	// switch back to original directory
	if err := os.Chdir(cwd); err != nil {
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

	// -build -> aurora-engine and neard
	if build {
		if err := buildAuroraEngine(bp.AuroraEngineHead); err != nil {
			return err
		}
		if err := buildNearcore(bp.NearcoreHead, release); err != nil {
			return err
		}
	}

	// copy neard directory
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	localDir := filepath.Join(home, ".near", "local")
	if err := os.RemoveAll(localDir); err != nil {
		return err
	}
	err = file.CopyDir(filepath.Join(breakpointDir, "local"), localDir)
	if err != nil {
		return err
	}

	// start neard
	n, err := neard.Start(release)
	if err != nil {
		return err
	}
	defer n.Stop()

	log.Info("sleep")
	time.Sleep(5 * time.Second)

	// copy credentials file
	dst := filepath.Join(home, ".near-credentials", "local", bp.AccountID+".json")
	if err := os.RemoveAll(dst); err != nil {
		return err
	}
	err = file.Copy(filepath.Join(breakpointDir, bp.AccountID+".json"), dst)
	if err != nil {
		return err
	}

	// upgrade contract before replaying tx, if necessary
	if contract != "" {
		if err := aurora.Upgrade(bp.AccountID, contract); err != nil {
			return err
		}
	}

	// run transaction
	zeroAmount := big.NewInt(0)
	rlp, err := hex.DecodeString(bp.Transaction)
	if err != nil {
		return err
	}

	if err := os.Setenv("NEAR_ENV", "local"); err != nil {
		return err
	}
	cfg := near.GetConfig()
	c := near.NewConnection(cfg.NodeURL)
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
