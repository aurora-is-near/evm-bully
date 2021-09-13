// Package neard implements NEAR daemon related functionality.
package neard

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/aurora-is-near/evm-bully/util/git"
	"github.com/ethereum/go-ethereum/log"
	"github.com/frankbraun/codechain/util/file"
)

// NEARDaemon wraps a neard.
type NEARDaemon struct {
	Head       string
	binaryPath string
	localDir   string
	nearDaemon *exec.Cmd
}

func getDefaultLocalDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".near", "local"), nil
}

// LoadFromBinary loads NEARDaemon from existing binary.
func LoadFromBinary(binaryPath string, head string) (*NEARDaemon, error) {
	log.Info("load neard binary")

	if _, err := os.Stat(binaryPath); err != nil {
		return nil, fmt.Errorf("can't access neard binary on path %v: %v", binaryPath, err)
	}

	localDir, err := getDefaultLocalDir()
	if err != nil {
		return nil, err
	}

	return &NEARDaemon{
		Head:       head,
		binaryPath: binaryPath,
		localDir:   localDir,
	}, nil
}

func buildBinary(repoPath string, release bool) error {
	log.Info("build neard")

	args := []string{
		"build",
		"--package", "neard",
		"--features", "nightly_protocol_features",
	}
	if release {
		args = append(args, "--release")
	}

	cmd := exec.Command("cargo", args...)
	cmd.Dir = repoPath
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// LoadFromRepo loads NEARDaemon from repo by provided path (and build if requested).
func LoadFromRepo(repoPath string, release bool, build bool) (*NEARDaemon, error) {
	log.Info("load neard repo")

	head, err := git.Head(repoPath)
	if err != nil {
		return nil, err
	}

	if build {
		if err := buildBinary(repoPath, release); err != nil {
			return nil, err
		}
	}

	binaryPath := filepath.Join(repoPath, "target", "debug", "neard")
	if release {
		binaryPath = filepath.Join(repoPath, "target", "release", "neard")
	}

	return LoadFromBinary(binaryPath, head)
}

// Backup local data if it exists
func (daemon *NEARDaemon) backupLocalData() error {
	exists, err := file.Exists(daemon.localDir)
	if err != nil {
		return err
	}
	if exists {
		localOld := strings.TrimSuffix(daemon.localDir, "/") + "_old"
		log.Info(fmt.Sprintf("mv %s %s", daemon.localDir, localOld))
		// remove old backup directory
		if err := os.RemoveAll(localOld); err != nil {
			return err
		}
		// move
		if err := os.Rename(daemon.localDir, localOld); err != nil {
			return err
		}
	} else {
		log.Info(fmt.Sprintf("directory '%s' does not exist", daemon.localDir))
	}
	return nil
}

func (daemon *NEARDaemon) init() error {
	cmd := exec.Command(daemon.binaryPath, "--home="+daemon.localDir, "--verbose=true", "init")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func (daemon *NEARDaemon) editGenesis() error {
	filename := filepath.Join(daemon.localDir, "genesis.json")
	backup := filepath.Join(daemon.localDir, "genesis_old.json")
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

// SetupLocalData initializes local data of a NEARDaemon.
func (daemon *NEARDaemon) SetupLocalData() error {
	log.Info("setup neard local data")

	if err := daemon.backupLocalData(); err != nil {
		return err
	}

	// initialize neard
	if err := daemon.init(); err != nil {
		return err
	}

	// edit genesis.json
	if err := daemon.editGenesis(); err != nil {
		return err
	}

	return nil
}

// RestoreLocalData restores local data of a NEARDaemon from given directory.
func (daemon *NEARDaemon) RestoreLocalData(source string) error {
	log.Info("restore neard local data")
	return file.CopyDir(source, daemon.localDir)
}

// Start NEARDaemon.
func (daemon *NEARDaemon) Start() error {
	log.Info("start neard")
	daemon.nearDaemon = exec.Command(daemon.binaryPath, "--home="+daemon.localDir, "--verbose=true", "run")
	daemon.nearDaemon.Stdout = os.Stdout
	daemon.nearDaemon.Stderr = os.Stderr
	return daemon.nearDaemon.Start()
}

// Stop NEARDaemon.
func (daemon *NEARDaemon) Stop() error {
	log.Info("stop neard")
	return daemon.nearDaemon.Process.Kill()
}
