// Package detector identifies the technology stack of a local repository
// by inspecting well-known marker files.
package detector

import (
	"os"
	"path/filepath"

	"github.com/MUSTHAFA786G/devsetup/internal/logger"
)

// StackType is the name of a detected technology stack.
type StackType string

const (
	StackNode    StackType = "Node.js"
	StackPython  StackType = "Python"
	StackGo      StackType = "Go"
	StackJava    StackType = "Java"
	StackRuby    StackType = "Ruby"
	StackRust    StackType = "Rust"
	StackUnknown StackType = "Unknown"
)

// Stack holds everything devsetup needs to know about a project's tech stack.
type Stack struct {
	// Name is the human-readable stack identifier.
	Name StackType
	// Marker is the file whose presence triggered detection.
	Marker string
	// Icon is an emoji used in terminal output.
	Icon string
	// Description is a one-line human summary.
	Description string
	// InstallCmds are the commands run to install dependencies.
	InstallCmds []string
	// DevCmds are the candidates for starting the dev server (tried in order).
	DevCmds []string
}

// IsUnknown returns true when no stack could be identified.
func (s *Stack) IsUnknown() bool { return s.Name == StackUnknown }

// rule maps a marker filename to a Stack.
type rule struct {
	file  string
	stack Stack
}

// detectionRules are evaluated in order; the first match wins.
var detectionRules = []rule{
	{
		file: "package.json",
		stack: Stack{
			Name:        StackNode,
			Marker:      "package.json",
			Icon:        "🟩",
			Description: "JavaScript / TypeScript project managed by npm",
			InstallCmds: []string{"npm install"},
			DevCmds:     []string{"npm run dev", "npm start"},
		},
	},
	{
		file: "requirements.txt",
		stack: Stack{
			Name:        StackPython,
			Marker:      "requirements.txt",
			Icon:        "🐍",
			Description: "Python project with pip dependencies",
			InstallCmds: []string{"pip install -r requirements.txt"},
			DevCmds:     []string{"python main.py", "python app.py", "python manage.py runserver"},
		},
	},
	{
		file: "go.mod",
		stack: Stack{
			Name:        StackGo,
			Marker:      "go.mod",
			Icon:        "🐹",
			Description: "Go module project",
			InstallCmds: []string{"go mod download"},
			DevCmds:     []string{"go run ."},
		},
	},
	{
		file: "pom.xml",
		stack: Stack{
			Name:        StackJava,
			Marker:      "pom.xml",
			Icon:        "☕",
			Description: "Java project managed by Maven",
			InstallCmds: []string{"mvn install -DskipTests"},
			DevCmds:     []string{"mvn spring-boot:run", "java -jar target/*.jar"},
		},
	},
	{
		file: "Gemfile",
		stack: Stack{
			Name:        StackRuby,
			Marker:      "Gemfile",
			Icon:        "💎",
			Description: "Ruby project managed by Bundler",
			InstallCmds: []string{"bundle install"},
			DevCmds:     []string{"bundle exec rails server", "ruby app.rb"},
		},
	},
	{
		file: "Cargo.toml",
		stack: Stack{
			Name:        StackRust,
			Marker:      "Cargo.toml",
			Icon:        "🦀",
			Description: "Rust project managed by Cargo",
			InstallCmds: []string{"cargo build"},
			DevCmds:     []string{"cargo run"},
		},
	},
}

// Detector inspects a repository path and identifies its technology stack.
type Detector struct {
	log *logger.Logger
}

// New returns a new Detector.
func New(log *logger.Logger) *Detector {
	return &Detector{log: log}
}

// Detect walks the detection rules and returns the first matching Stack.
// When no rule matches it returns StackUnknown (not an error).
func (d *Detector) Detect(repoPath string) (*Stack, error) {
	d.log.Debug("Scanning %s for stack markers", repoPath)

	for _, r := range detectionRules {
		p := filepath.Join(repoPath, r.file)
		if _, err := os.Stat(p); err == nil {
			d.log.Debug("Found marker: %s", r.file)
			s := r.stack // value copy — callers may not mutate the global
			return &s, nil
		}
	}

	d.log.Warn("No recognized stack marker found — treating project as Unknown")
	return &Stack{
		Name:        StackUnknown,
		Marker:      "—",
		Icon:        "📦",
		Description: "No recognized dependency file found",
	}, nil
}
