package utils_test

import (
	"testing"

	"github.com/devsetup/devsetup/pkg/utils"
)

func TestCommandExists_GoPresent(t *testing.T) {
	// "go" is guaranteed present when running Go tests
	if !utils.CommandExists("go") {
		t.Skip("go binary not in PATH — skipping")
	}
}

func TestCommandExists_Missing(t *testing.T) {
	if utils.CommandExists("this-tool-should-never-exist-xyz123") {
		t.Error("expected false for non-existent command")
	}
}

func TestRequireCommand_Error(t *testing.T) {
	if err := utils.RequireCommand("this-tool-should-never-exist-xyz123"); err == nil {
		t.Error("expected error for missing command")
	}
}
