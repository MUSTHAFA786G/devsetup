// Package runner starts the development server for a detected stack and
// forwards OS signals so Ctrl-C cleanly terminates the child process.
package runner

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"

	"github.com/devsetup/devsetup/internal/detector"
	"github.com/devsetup/devsetup/internal/logger"
)

// Runner launches the dev server for a project.
type Runner struct {
	log *logger.Logger
}

// New returns a new Runner.
func New(log *logger.Logger) *Runner {
	return &Runner{log: log}
}

// Run starts the development server.
// For Node.js projects it inspects package.json to prefer the "dev" script.
// For all other stacks it uses the first DevCmd in the stack definition.
func (r *Runner) Run(repoPath string, stack *detector.Stack) error {
	if stack.IsUnknown() {
		r.log.Warn("Unknown stack — cannot determine how to start the project")
		r.log.Info("Navigate to %s and start the project manually.", repoPath)
		return nil
	}
	if len(stack.DevCmds) == 0 {
		r.log.Warn("No dev commands defined for %s", stack.Name)
		return nil
	}

	cmd := r.resolveCmd(repoPath, stack)
	r.log.Info("Starting dev server: %s", cmd)
	r.printRunBanner(cmd)
	return r.execInteractive(repoPath, cmd)
}

// resolveCmd picks the best dev command for the stack.
func (r *Runner) resolveCmd(repoPath string, stack *detector.Stack) string {
	if stack.Name == detector.StackNode {
		pkgJSON := readFile(repoPath + "/package.json")
		if strings.Contains(pkgJSON, `"dev"`) {
			return "npm run dev"
		}
		if strings.Contains(pkgJSON, `"start"`) {
			return "npm start"
		}
	}
	return stack.DevCmds[0]
}

// execInteractive starts cmd as a child process, wiring stdin/stdout/stderr
// and forwarding SIGINT / SIGTERM for a clean shutdown.
func (r *Runner) execInteractive(dir, cmdStr string) error {
	parts := strings.Fields(cmdStr)
	if len(parts) == 0 {
		return fmt.Errorf("empty command")
	}

	cmd := exec.Command(parts[0], parts[1:]...)
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Env = os.Environ()

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start %q: %w", cmdStr, err)
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	done := make(chan error, 1)
	go func() { done <- cmd.Wait() }()

	select {
	case sig := <-sigs:
		r.log.Info("Received %v — shutting down...", sig)
		if cmd.Process != nil {
			_ = cmd.Process.Signal(sig)
		}
		return <-done
	case err := <-done:
		return err
	}
}

// readFile reads a file and returns its content as a string (empty on error).
func readFile(path string) string {
	b, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	return string(b)
}

// printRunBanner displays a startup box before the server output begins.
func (r *Runner) printRunBanner(cmd string) {
	fmt.Println()
	fmt.Println("  \033[1;32m╔══════════════════════════════════════════╗\033[0m")
	fmt.Println("  \033[1;32m║\033[0m  🚀 \033[1mDev server starting\033[0m                   \033[1;32m║\033[0m")
	fmt.Println("  \033[1;32m╠══════════════════════════════════════════╣\033[0m")
	fmt.Printf("  \033[1;32m║\033[0m  \033[2mCommand:\033[0m %-32s\033[1;32m║\033[0m\n", cmd)
	fmt.Println("  \033[2m║  Press Ctrl+C to stop                    ║\033[0m")
	fmt.Println("  \033[1;32m╚══════════════════════════════════════════╝\033[0m")
	fmt.Println()
}
