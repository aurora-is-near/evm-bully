// Package git contains wrappers around some Git commands.
package git

import (
	"bytes"
	"os"
	"os/exec"
	"strings"
)

// Head returns the current Git HEAD in the current working directory.
func Head() (string, error) {
	cmd := exec.Command("git", "rev-parse", "HEAD")
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return "", err
	}
	return strings.TrimSpace(stdout.String()), nil
}

// Checkout checks out the given head in the current working directory.
func Checkout(head string) error {
	cmd := exec.Command("git", "checkout", head)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
