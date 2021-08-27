// Package neard implements NEAR daemon related functionality.
package neard

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/aurora-is-near/evm-bully/util/git"
	"github.com/ethereum/go-ethereum/log"
	"github.com/frankbraun/codechain/util/file"
)

// NEARDaemon wraps a running neard.
type NEARDaemon struct {
	Head       string
	nearDaemon *exec.Cmd
}

// Build builds neard in CWD.
func Build(release bool) error {
	args := []string{
		"build",
		"--package", "neard",
		"--features", "nightly_protocol_features",
	}
	if release {
		args = append(args, "--release")
	}
	cmd := exec.Command("cargo", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func initDaemon(release bool, localDir string) error {
	var name string
	if release {
		name = filepath.Join(".", "target", "release", "neard")
	} else {
		name = filepath.Join(".", "target", "debug", "neard")
	}
	cmd := exec.Command(name, "--home="+localDir, "--verbose=true", "init")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func editGenesis(localDir string) error {
	filename := filepath.Join(localDir, "genesis.json")
	backup := filepath.Join(localDir, "genesis_old.json")
	if err := file.Copy(filename, backup); err != nil {
		return err
	}
	// read file
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	// change default values the brute force way, neard chokes on edited JSON
	data = bytes.Replace(data,
		[]byte("\"max_gas_burnt\": 200000000000000"),
		[]byte("\"max_gas_burnt\": 800000000000000"),
		1)
	data = bytes.Replace(data,
		[]byte("\"max_total_prepaid_gas\": 300000000000000"),
		[]byte("\"max_total_prepaid_gas\": 800000000000000"),
		1)
	// write file
	if err := os.WriteFile(filename, data, 0644); err != nil {
		return err
	}
	return nil
}

func (n *NEARDaemon) start(release bool, localDir string) error {
	var name string
	if release {
		name = filepath.Join(".", "target", "release", "neard")
	} else {
		name = filepath.Join(".", "target", "debug", "neard")
	}
	n.nearDaemon = exec.Command(name, "--home="+localDir, "--verbose=true", "run")
	n.nearDaemon.Stdout = os.Stdout
	n.nearDaemon.Stderr = os.Stderr
	return n.nearDaemon.Start()
}

// Setup initializes and starts a (release) NEARDaemon.
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
	n.Head, err = git.Head()
	if err != nil {
		return nil, err
	}
	log.Info(fmt.Sprintf("head=%s", n.Head))
	// backup .near/local, if it exists
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	localDir := filepath.Join(home, ".near", "local")
	exists, err := file.Exists(localDir)
	if err != nil {
		return nil, err
	}
	if exists {
		localOld := localDir + "_old"
		log.Info(fmt.Sprintf("mv %s %s", localDir, localOld))
		// remove old backup directory
		if err := os.RemoveAll(localOld); err != nil {
			return nil, err
		}
		// move
		if err := os.Rename(localDir, localOld); err != nil {
			return nil, err
		}
	} else {
		log.Info(fmt.Sprintf("directory '%s' does not exist", localDir))
	}
	// make sure neard is build
	if err := Build(release); err != nil {
		return nil, err
	}
	// initialize neard
	if err := initDaemon(release, localDir); err != nil {
		return nil, err
	}
	// edit genesis.json
	if err := editGenesis(localDir); err != nil {
		return nil, err
	}
	// start neard
	if err := n.start(release, localDir); err != nil {
		return nil, err
	}
	// switch back to original directory
	if err := os.Chdir(cwd); err != nil {
		return nil, err
	}
	return &n, nil
}

// Start starts a (release) NEARDaemon.
func Start(release bool) (*NEARDaemon, error) {
	var n NEARDaemon
	log.Info("start neard")
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
	// start neard
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	localDir := filepath.Join(home, ".near", "local")
	if err := n.start(release, localDir); err != nil {
		return nil, err
	}
	// switch back to original directory
	if err := os.Chdir(cwd); err != nil {
		return nil, err
	}
	return &n, nil
}

// Stop NEARDaemon.
func (n *NEARDaemon) Stop() error {
	log.Info("stop neard")
	return n.nearDaemon.Process.Kill()
}
