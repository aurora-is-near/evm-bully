// Package neard implements NEAR daemon related functionality.
package neard

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/aurora-is-near/evm-bully/util/git"
	"github.com/ethereum/go-ethereum/log"
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
