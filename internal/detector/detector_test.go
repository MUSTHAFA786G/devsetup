package detector_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/devsetup/devsetup/internal/detector"
	"github.com/devsetup/devsetup/internal/logger"
)

func log() *logger.Logger { return logger.New(false) }

func touch(t *testing.T, dir, name, content string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0o644); err != nil {
		t.Fatalf("touch %s: %v", name, err)
	}
}

func TestDetect_Node(t *testing.T) {
	dir := t.TempDir()
	touch(t, dir, "package.json", `{"name":"test"}`)
	s, err := detector.New(log()).Detect(dir)
	assertStack(t, err, s, detector.StackNode)
}

func TestDetect_Python(t *testing.T) {
	dir := t.TempDir()
	touch(t, dir, "requirements.txt", "flask\n")
	s, err := detector.New(log()).Detect(dir)
	assertStack(t, err, s, detector.StackPython)
}

func TestDetect_Go(t *testing.T) {
	dir := t.TempDir()
	touch(t, dir, "go.mod", "module example.com/x\ngo 1.21\n")
	s, err := detector.New(log()).Detect(dir)
	assertStack(t, err, s, detector.StackGo)
}

func TestDetect_Java(t *testing.T) {
	dir := t.TempDir()
	touch(t, dir, "pom.xml", "<project/>")
	s, err := detector.New(log()).Detect(dir)
	assertStack(t, err, s, detector.StackJava)
}

func TestDetect_Ruby(t *testing.T) {
	dir := t.TempDir()
	touch(t, dir, "Gemfile", `source "https://rubygems.org"`)
	s, err := detector.New(log()).Detect(dir)
	assertStack(t, err, s, detector.StackRuby)
}

func TestDetect_Rust(t *testing.T) {
	dir := t.TempDir()
	touch(t, dir, "Cargo.toml", "[package]\nname=\"x\"")
	s, err := detector.New(log()).Detect(dir)
	assertStack(t, err, s, detector.StackRust)
}

func TestDetect_Unknown(t *testing.T) {
	dir := t.TempDir()
	s, err := detector.New(log()).Detect(dir)
	assertStack(t, err, s, detector.StackUnknown)
	if !s.IsUnknown() {
		t.Error("IsUnknown() should return true")
	}
}

func assertStack(t *testing.T, err error, s *detector.Stack, want detector.StackType) {
	t.Helper()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.Name != want {
		t.Errorf("got %q, want %q", s.Name, want)
	}
}
