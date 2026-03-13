// Package utils provides shared helper functions used across devsetup modules.
package utils

import (
	"fmt"
	"os/exec"
)

// RequireCommand checks whether a CLI tool exists in PATH.
// Returns a descriptive error when it is missing.
func RequireCommand(name string) error {
	_, err := exec.LookPath(name)
	if err != nil {
		return fmt.Errorf("'%s' not found in PATH — please install it and try again", name)
	}
	return nil
}

// CommandExists returns true when the named binary is available in PATH.
func CommandExists(name string) bool {
	return RequireCommand(name) == nil
}
