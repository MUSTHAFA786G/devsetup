package analyzer_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/MUSTHAFA786G/devsetup/internal/analyzer"
	"github.com/MUSTHAFA786G/devsetup/internal/detector"
	"github.com/MUSTHAFA786G/devsetup/internal/logger"
)

func TestAnalyze_BasicReport(t *testing.T) {
	dir := t.TempDir()
	// Create a small fake Go project
	os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module x\ngo 1.21"), 0o644)
	os.WriteFile(filepath.Join(dir, "main.go"), []byte("package main"), 0o644)
	os.WriteFile(filepath.Join(dir, "README.md"), []byte("# hi"), 0o644)

	stack := &detector.Stack{
		Name:   detector.StackGo,
		Marker: "go.mod",
		Icon:   "🐹",
	}
	a := analyzer.New(logger.New(false))
	report, err := a.Analyze(dir, "myrepo", stack)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if report.RepoName != "myrepo" {
		t.Errorf("RepoName = %q, want %q", report.RepoName, "myrepo")
	}
	if report.TotalFiles < 3 {
		t.Errorf("TotalFiles = %d, want ≥ 3", report.TotalFiles)
	}
	if report.FileCounts["Go"] == 0 {
		t.Error("expected Go files to be counted")
	}
	if len(report.EntryPoints) == 0 {
		t.Error("expected main.go to be detected as entry point")
	}
}
