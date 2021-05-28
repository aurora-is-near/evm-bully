// Package gnumake contains wrappers around some Make commands.
package gnumake

import (
	"os"
	"os/exec"
)

// Make calls make in the current working directory.
func Make(args ...string) error {
	cmd := exec.Command("make", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
