// Package gnumake contains wrappers around some Make commands.
package gnumake

import (
	"os"
	"os/exec"
)

// Make calls make in the directory by provided path.
func Make(path string, args ...string) error {
	cmd := exec.Command("make", args...)
	cmd.Dir = path
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
