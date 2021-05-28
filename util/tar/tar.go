// Package tar contains wrappers around some tar commands.
package tar

import (
	"os"
	"os/exec"
)

// Create creates a tar archive for the given dir.
func Create(dir string) error {
	cmd := exec.Command("tar", "cvzf", dir+".tar.gz", dir)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
