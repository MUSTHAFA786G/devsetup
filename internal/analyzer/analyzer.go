// Package analyzer inspects a cloned repository and produces a human-readable
// architecture summary: directory tree, file-type statistics, and detected
// entry points / config files.
package analyzer

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/devsetup/devsetup/internal/detector"
	"github.com/devsetup/devsetup/internal/logger"
)

// Report is the result of analyzing a repository.
type Report struct {
	RepoName   string
	RepoPath   string
	Stack      *detector.Stack
	Tree       []string          // Formatted tree lines
	FileCounts map[string]int    // language → file count
	EntryPoints []string         // Detected entry-point files
	ConfigFiles []string         // Detected configuration files
	TotalFiles  int
}

// Analyzer builds a Report for a local repository.
type Analyzer struct {
	log *logger.Logger
}

// New returns a new Analyzer.
func New(log *logger.Logger) *Analyzer {
	return &Analyzer{log: log}
}

// Analyze inspects repoPath and returns a Report.
func (a *Analyzer) Analyze(repoPath, repoName string, stack *detector.Stack) (*Report, error) {
	a.log.Debug("Analyzing repository at %s", repoPath)

	tree := buildTree(repoPath, 2)
	counts, total := countFiles(repoPath)
	entries := detectEntryPoints(repoPath, stack)
	configs := detectConfigs(repoPath)

	return &Report{
		RepoName:    repoName,
		RepoPath:    repoPath,
		Stack:       stack,
		Tree:        tree,
		FileCounts:  counts,
		EntryPoints: entries,
		ConfigFiles: configs,
		TotalFiles:  total,
	}, nil
}

// Display prints the Report to stdout in a formatted layout.
func (a *Analyzer) Display(r *Report) {
	c := func(code, text string) string { return code + text + "\033[0m" }

	fmt.Println()
	fmt.Println(c("\033[1;36m", "  ┌─────────────────────────────────────────────────┐"))
	fmt.Printf(  "  │  📋  %-43s│\n", "Project Architecture: "+r.RepoName)
	fmt.Println(c("\033[1;36m", "  └─────────────────────────────────────────────────┘"))
	fmt.Println()

	// ── Stack overview ───────────────────────────────────────────────────────
	printField := func(label, value string) {
		fmt.Printf("  %s%-18s%s%s\n",
			"\033[1m", label+":", "\033[0m", value)
	}
	printField("Stack", fmt.Sprintf("%s  %s", r.Stack.Icon, r.Stack.Name))
	printField("Description", r.Stack.Description)
	printField("Detected via", r.Stack.Marker)
	printField("Total files", fmt.Sprintf("%d", r.TotalFiles))
	printField("Location", r.RepoPath)

	// ── Directory tree ───────────────────────────────────────────────────────
	fmt.Println()
	fmt.Println(c("\033[1;36m", "  Directory Structure:"))
	for _, line := range r.Tree {
		fmt.Printf("  %s\n", c("\033[2m", line))
	}

	// ── File statistics ──────────────────────────────────────────────────────
	if len(r.FileCounts) > 0 {
		fmt.Println()
		fmt.Println(c("\033[1;36m", "  File Statistics:"))
		langs := sortedKeys(r.FileCounts)
		maxCount := 0
		for _, l := range langs {
			if r.FileCounts[l] > maxCount {
				maxCount = r.FileCounts[l]
			}
		}
		for _, lang := range langs {
			n := r.FileCounts[lang]
			barLen := 0
			if maxCount > 0 {
				barLen = (n * 20) / maxCount
				if barLen == 0 {
					barLen = 1
				}
			}
			bar := strings.Repeat("█", barLen)
			fmt.Printf("    \033[32m%-14s\033[33m%-22s\033[2m %d file",
				lang, bar, n)
			if n != 1 {
				fmt.Print("s")
			}
			fmt.Println("\033[0m")
		}
	}

	// ── Entry points ─────────────────────────────────────────────────────────
	if len(r.EntryPoints) > 0 {
		fmt.Println()
		fmt.Println(c("\033[1;36m", "  Detected Entry Points:"))
		for _, ep := range r.EntryPoints {
			fmt.Printf("  \033[32m  ▸  \033[0m%s\n", ep)
		}
	}

	// ── Config files ─────────────────────────────────────────────────────────
	if len(r.ConfigFiles) > 0 {
		fmt.Println()
		fmt.Println(c("\033[1;36m", "  Config Files:"))
		for _, cf := range r.ConfigFiles {
			fmt.Printf("  \033[34m  ▸  \033[0m%s\n", cf)
		}
	}

	// ── Commands ─────────────────────────────────────────────────────────────
	fmt.Println()
	if len(r.Stack.InstallCmds) > 0 {
		fmt.Println(c("\033[1;36m", "  Install:"))
		for _, cmd := range r.Stack.InstallCmds {
			fmt.Printf("    \033[2m$ %s\033[0m\n", cmd)
		}
	}
	if len(r.Stack.DevCmds) > 0 {
		fmt.Println(c("\033[1;36m", "  Run:"))
		for _, cmd := range r.Stack.DevCmds {
			fmt.Printf("    \033[2m$ %s\033[0m\n", cmd)
		}
	}
	fmt.Println()
}

// ── Helpers ───────────────────────────────────────────────────────────────────

var skipDirs = map[string]bool{
	"node_modules": true, ".git": true, "vendor": true,
	"dist": true, "build": true, "target": true,
	".next": true, "__pycache__": true, ".idea": true, ".vscode": true,
}

func buildTree(root string, maxDepth int) []string {
	var lines []string
	walkTree(root, 0, maxDepth, "", &lines)
	return lines
}

func walkTree(dir string, depth, maxDepth int, prefix string, lines *[]string) {
	if depth > maxDepth {
		return
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		return
	}
	visible := filterVisible(entries)
	for i, e := range visible {
		isLast := i == len(visible)-1
		connector := "├── "
		childPrefix := prefix + "│   "
		if isLast {
			connector = "└── "
			childPrefix = prefix + "    "
		}
		name := e.Name()
		if e.IsDir() {
			name += "/"
		}
		*lines = append(*lines, prefix+connector+name)
		if e.IsDir() && depth < maxDepth {
			walkTree(filepath.Join(dir, e.Name()), depth+1, maxDepth, childPrefix, lines)
		}
	}
}

func filterVisible(entries []os.DirEntry) []os.DirEntry {
	var out []os.DirEntry
	for _, e := range entries {
		if strings.HasPrefix(e.Name(), ".") {
			continue
		}
		if skipDirs[e.Name()] {
			continue
		}
		out = append(out, e)
	}
	return out
}

func countFiles(root string) (map[string]int, int) {
	counts := make(map[string]int)
	total := 0
	_ = filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if d.IsDir() && skipDirs[d.Name()] {
			return filepath.SkipDir
		}
		if !d.IsDir() {
			total++
			if lang := extToLang(strings.ToLower(filepath.Ext(d.Name()))); lang != "" {
				counts[lang]++
			}
		}
		return nil
	})
	return counts, total
}

func extToLang(ext string) string {
	m := map[string]string{
		".go": "Go", ".js": "JavaScript", ".ts": "TypeScript",
		".jsx": "JSX", ".tsx": "TSX", ".py": "Python",
		".java": "Java", ".rb": "Ruby", ".rs": "Rust",
		".html": "HTML", ".css": "CSS", ".scss": "SCSS",
		".md": "Markdown", ".yaml": "YAML", ".yml": "YAML",
		".json": "JSON", ".toml": "TOML", ".sh": "Shell",
		".c": "C", ".cpp": "C++", ".h": "C/C++ Header",
	}
	return m[ext]
}

var entryPointCandidates = []string{
	"main.go", "main.py", "app.py", "index.js", "index.ts",
	"server.js", "server.ts", "app.js", "app.ts",
	"src/main.go", "src/index.js", "src/index.ts", "src/app.ts",
	"src/main.py", "cmd/main.go",
}

func detectEntryPoints(root string, _ *detector.Stack) []string {
	var found []string
	for _, candidate := range entryPointCandidates {
		if _, err := os.Stat(filepath.Join(root, candidate)); err == nil {
			found = append(found, candidate)
		}
	}
	return found
}

var configCandidates = []string{
	"Dockerfile", "docker-compose.yml", "docker-compose.yaml",
	".env.example", ".env.sample", "config.yaml", "config.yml",
	"config.json", "tsconfig.json", "vite.config.ts", "vite.config.js",
	"webpack.config.js", ".eslintrc.js", ".prettierrc",
	"pyproject.toml", "setup.py", "Makefile",
}

func detectConfigs(root string) []string {
	var found []string
	for _, candidate := range configCandidates {
		if _, err := os.Stat(filepath.Join(root, candidate)); err == nil {
			found = append(found, candidate)
		}
	}
	return found
}

func sortedKeys(m map[string]int) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		return m[keys[i]] > m[keys[j]] // descending by count
	})
	return keys
}
