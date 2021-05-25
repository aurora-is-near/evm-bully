// Package neard implements NEAR daemon related functionality.
package neard

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/aurora-is-near/evm-bully/util/git"
	"github.com/ethereum/go-ethereum/log"
	"github.com/frankbraun/codechain/util/file"
)

type NEARDaemon struct {
	head string
}

func build(release bool) error {
	args := []string{
		"build",
		"--package", "neard",
		"--features", "protocol_feature_evm,nightly_protocol_features",
	}
	if release {
		args = append(args, "--release")
	}
	cmd := exec.Command("cargo", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func Setup(release bool) (*NEARDaemon, error) {
	var n NEARDaemon
	log.Info("setup neard")
	// get cwd
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	// switch to nearcore directory
	nearDir := filepath.Join(cwd, "..", "nearcore")
	if err := os.Chdir(nearDir); err != nil {
		return nil, err
	}
	// get current HEAD
	n.head, err = git.Head()
	if err != nil {
		return nil, err
	}
	log.Info(fmt.Sprintf("head=%s", n.head))
	// backup .near/local, if it exists
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	local := filepath.Join(home, ".near", "local")
	exists, err := file.Exists(local)
	if err != nil {
		return nil, err
	}
	if exists {
		localOld := local + "_old"
		log.Info(fmt.Sprintf("mv %s %s", local, localOld))
		// remove old backup directory
		if err := os.RemoveAll(localOld); err != nil {
			return nil, err
		}
		// move
		if err := os.Rename(local, localOld); err != nil {
			return nil, err
		}
	} else {
		log.Info(fmt.Sprintf("directory '%s' does not exist", local))
	}
	// make sure neard is build
	if err := build(release); err != nil {
		return nil, err
	}
	// switch back to original directory
	if err := os.Chdir(cwd); err != nil {
		return nil, err
	}
	return &n, nil
}

func (n *NEARDaemon) Stop() {
	log.Info("stop neard")
}
