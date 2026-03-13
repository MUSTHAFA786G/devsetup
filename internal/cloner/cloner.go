// Package cloner handles cloning of remote Git repositories.
package cloner

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/MUSTHAFA786G/devsetup/internal/logger"
	"github.com/MUSTHAFA786G/devsetup/pkg/utils"
)

// Cloner clones a remote Git repository to a local path.
type Cloner struct {
	log *logger.Logger
}

// Result holds the outcome of a successful clone operation.
type Result struct {
	// RepoPath is the absolute path of the cloned directory.
	RepoPath string
	// RepoName is the repository name extracted from the URL.
	RepoName string
	// AlreadyExisted is true when the target directory was already present.
	AlreadyExisted bool
}

// New returns a new Cloner.
func New(log *logger.Logger) *Cloner {
	return &Cloner{log: log}
}

// Clone clones repoURL into targetDir and returns a Result.
func (c *Cloner) Clone(repoURL, targetDir string) (*Result, error) {
	if err := utils.RequireCommand("git"); err != nil {
		return nil, fmt.Errorf("git is required: %w", err)
	}

	repoURL = normalizeURL(repoURL)
	repoName := extractRepoName(repoURL)
	if repoName == "" {
		return nil, fmt.Errorf("could not extract repository name from URL: %q", repoURL)
	}

	absTarget, err := filepath.Abs(targetDir)
	if err != nil {
		return nil, fmt.Errorf("invalid target directory %q: %w", targetDir, err)
	}
	repoPath := filepath.Join(absTarget, repoName)

	// Destination already exists — skip clone.
	if _, statErr := os.Stat(repoPath); statErr == nil {
		c.log.Warn("Directory %q already exists — skipping clone", repoPath)
		return &Result{RepoPath: repoPath, RepoName: repoName, AlreadyExisted: true}, nil
	}

	if err := os.MkdirAll(absTarget, 0o755); err != nil {
		return nil, fmt.Errorf("cannot create target directory: %w", err)
	}

	c.log.Info("Cloning from %s", repoURL)
	c.log.Command(fmt.Sprintf("git clone %s %s", repoURL, repoPath))

	cmd := exec.Command("git", "clone", "--progress", repoURL, repoPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("git clone failed: %w", err)
	}

	return &Result{RepoPath: repoPath, RepoName: repoName}, nil
}

// normalizeURL converts SSH URLs to HTTPS and ensures a .git suffix.
func normalizeURL(url string) string {
	// SSH → HTTPS
	if strings.HasPrefix(url, "git@github.com:") {
		url = strings.Replace(url, "git@github.com:", "https://github.com/", 1)
	}
	url = strings.TrimRight(url, "/")
	if !strings.HasSuffix(url, ".git") {
		url += ".git"
	}
	return url
}

// extractRepoName parses the repository name from a URL.
//
//	https://github.com/user/project.git  →  project
func extractRepoName(url string) string {
	parts := strings.Split(url, "/")
	if len(parts) == 0 {
		return ""
	}
	return strings.TrimSuffix(parts[len(parts)-1], ".git")
}
