// Package hashcache implements a cache file for block hashes.
package hashcache

import (
	"bufio"
	"os"
	"path/filepath"

	"github.com/ethereum/go-ethereum/common"
	"github.com/frankbraun/codechain/util/file"
	"github.com/frankbraun/codechain/util/lockfile"
)

const (
	defaultFilename = "hashcache.txt"
)

// Load block hashes from hash cache file in cacheDir.
func Load(cacheDir string) ([]common.Hash, error) {
	filename := filepath.Join(cacheDir, defaultFilename)
	l, err := lockfile.Create(filename)
	if err != nil {
		return nil, err
	}
	defer l.Release()
	exists, err := file.Exists(filename)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, nil
	}
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	s := bufio.NewScanner(f)
	var blocks []common.Hash
	for s.Scan() {
		blocks = append(blocks, common.HexToHash(s.Text()))
	}
	if err := s.Err(); err != nil {
		return nil, err
	}
	return blocks, nil

}

// Save hashes from blocks in hash cache file in cacheDir.
func Save(cacheDir string, blocks []common.Hash) error {
	filename := filepath.Join(cacheDir, defaultFilename)
	l, err := lockfile.Create(filename)
	if err != nil {
		return err
	}
	defer l.Release()
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	for _, block := range blocks {
		if _, err := f.WriteString(block.Hex() + "\n"); err != nil {
			return err
		}
	}
	return nil
}
