// Package installer runs the dependency installation commands for a
// detected technology stack.
package installer

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/devsetup/devsetup/internal/detector"
	"github.com/devsetup/devsetup/internal/logger"
	"github.com/devsetup/devsetup/pkg/utils"
)

// Installer executes dependency installation for a given stack.
type Installer struct {
	log *logger.Logger
}

// New returns a new Installer.
func New(log *logger.Logger) *Installer {
	return &Installer{log: log}
}

// Install runs all InstallCmds defined on the stack inside repoPath.
func (i *Installer) Install(repoPath string, stack *detector.Stack) error {
	if stack.IsUnknown() {
		i.log.Warn("Unknown stack — skipping dependency installation")
		return nil
	}
	if len(stack.InstallCmds) == 0 {
		i.log.Warn("No install commands defined for %s", stack.Name)
		return nil
	}

	// Verify the primary tool is available before doing any work.
	if tool := primaryTool(stack); tool != "" {
		if err := utils.RequireCommand(tool); err != nil {
			return fmt.Errorf(
				"required tool not found — %w\n  "+
					"Install %s and re-run devsetup", err, tool)
		}
		i.log.Debug("Verified tool in PATH: %s", tool)
	}

	for _, cmdStr := range stack.InstallCmds {
		i.log.Info("Running: %s", cmdStr)
		i.log.Command(cmdStr)
		if err := runInDir(repoPath, cmdStr); err != nil {
			return fmt.Errorf("install command %q failed: %w", cmdStr, err)
		}
	}
	return nil
}

// runInDir executes a space-separated command string inside dir.
func runInDir(dir, cmdStr string) error {
	parts := strings.Fields(cmdStr)
	if len(parts) == 0 {
		return fmt.Errorf("empty command")
	}
	cmd := exec.Command(parts[0], parts[1:]...)
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = os.Environ()
	return cmd.Run()
}

// primaryTool returns the main executable name for a stack.
func primaryTool(stack *detector.Stack) string {
	switch stack.Name {
	case detector.StackNode:
		return "npm"
	case detector.StackPython:
		return "pip"
	case detector.StackGo:
		return "go"
	case detector.StackJava:
		return "mvn"
	case detector.StackRuby:
		return "bundle"
	case detector.StackRust:
		return "cargo"
	default:
		return ""
	}
}
