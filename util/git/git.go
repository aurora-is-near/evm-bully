// Package git contains wrappers around some Git commands.
package git

import (
	"bytes"
	"os"
	"os/exec"
	"strings"
)

// Head returns the current Git HEAD in the repository by given path.
func Head(repoPath string) (string, error) {
	cmd := exec.Command("git", "rev-parse", "HEAD")
	cmd.Dir = repoPath
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return "", err
	}
	return strings.TrimSpace(stdout.String()), nil
}

// Checkout checks out the given head in the repository by given path.
func Checkout(repoPath string, head string) error {
	cmd := exec.Command("git", "checkout", head)
	cmd.Dir = repoPath
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
