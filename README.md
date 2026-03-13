# devsetup 🚀

> **One command. Any GitHub repo. Running in minutes.**

[![CI](https://github.com/devsetup/devsetup/actions/workflows/ci.yml/badge.svg)](https://github.com/devsetup/devsetup/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/devsetup/devsetup)](https://goreportcard.com/report/github.com/devsetup/devsetup)
[![Release](https://img.shields.io/github/v/release/devsetup/devsetup)](https://github.com/devsetup/devsetup/releases)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8.svg)](https://go.dev)

`devsetup` automatically clones a GitHub repository, identifies its technology
stack, installs dependencies, analyses the project architecture, and launches
the development server — all from a single command.

```bash
devsetup https://github.com/user/project
```

---

## Table of Contents

- [Features](#features)
- [Installation](#installation)
- [Usage](#usage)
- [Example CLI Output](#example-cli-output)
- [Supported Stacks](#supported-stacks)
- [Architecture](#architecture)
- [Contributing](#contributing)
- [License](#license)

---

## Features

| Feature | Description |
|---|---|
| 🔍 **Auto-detection** | Identifies Node.js, Python, Go, Java, Ruby, Rust by marker files |
| 📦 **Dependency install** | Runs the correct package manager automatically |
| 🌳 **Architecture summary** | Prints a directory tree, file-type stats, entry points & config files |
| 🚀 **Dev server launch** | Starts the right dev server with signal-safe process management |
| 🎨 **Colorized output** | ANSI-colored, step-by-step pipeline output |
| 🖥 **Cross-platform** | Statically compiled — Linux, macOS (Intel + Apple Silicon), Windows |
| 🔌 **Zero dependencies** | Pure Go standard library — no external packages required |
| ⚡ **Fast** | Single static binary, instant startup |

---

## Installation

### Option 1 — Go install (recommended)

```bash
go install github.com/devsetup/devsetup@latest
```

### Option 2 — Homebrew (macOS / Linux)

```bash
brew install devsetup/tap/devsetup
```

### Option 3 — Pre-built binary

Download the binary for your platform from the
[Releases page](https://github.com/devsetup/devsetup/releases):

| Platform | File |
|---|---|
| Linux x86-64 | `devsetup_linux_x86_64.tar.gz` |
| Linux ARM64 | `devsetup_linux_arm64.tar.gz` |
| macOS Intel | `devsetup_macOS_x86_64.tar.gz` |
| macOS Apple Silicon | `devsetup_macOS_arm64.tar.gz` |
| Windows x86-64 | `devsetup_Windows_x86_64.zip` |

Extract and move to your `$PATH`:

```bash
tar -xzf devsetup_*.tar.gz
sudo mv devsetup /usr/local/bin/
```

### Option 4 — Build from source

**Prerequisites:** Go 1.21+, Git

```bash
git clone https://github.com/devsetup/devsetup.git
cd devsetup
go build -o devsetup .
sudo mv devsetup /usr/local/bin/   # optional: add to PATH
```

Or use the Makefile:

```bash
make build    # build ./devsetup
make install  # install to GOPATH/bin
```

---

## Usage

```
devsetup [flags] <github-url>
devsetup version

Flags:
  --dir            Target directory to clone into (default: ".")
  --skip-install   Skip dependency installation
  --skip-run       Skip starting the dev server
  --verbose        Enable verbose/debug output
  --help           Show this help
```

### Examples

```bash
# Full setup: clone → detect → install → analyze → run
devsetup https://github.com/user/project

# Clone + analyze only (no install, no run)
devsetup https://github.com/user/project --skip-install --skip-run

# Clone into a custom directory
devsetup https://github.com/user/project --dir ~/projects

# Install deps but don't start the dev server
devsetup https://github.com/user/project --skip-run

# Debug mode: see every command and decision
devsetup https://github.com/user/project --verbose

# SSH URLs work too
devsetup git@github.com:user/project.git

# Print version info
devsetup version
```

---

## Example CLI Output

```
  ╔══════════════════════════════════════════════════╗
  ║   ██████╗ ███████╗██╗   ██╗███████╗███████╗████████╗██╗   ██╗██████╗  ║
  ║   ██╔══██╗██╔════╝██║   ██║██╔════╝██╔════╝╚══██╔══╝██║   ██║██╔══██╗ ║
  ║   ██║  ██║█████╗  ██║   ██║███████╗█████╗     ██║   ██║   ██║██████╔╝ ║
  ║   ██║  ██║██╔══╝  ╚██╗ ██╔╝╚════██║██╔══╝     ██║   ██║   ██║██╔═══╝  ║
  ║   ██████╔╝███████╗ ╚████╔╝ ███████║███████╗   ██║   ╚██████╔╝██║      ║
  ║   ╚═════╝ ╚══════╝  ╚═══╝  ╚══════╝╚══════╝   ╚═╝    ╚═════╝ ╚═╝      ║
  ╠══════════════════════════════════════════════════╣
  ║  Instant Dev Environment Setup  •  2024-01-15 10:32:07          ║
  ╚══════════════════════════════════════════════════╝

  ℹ  Repository : https://github.com/example/project

  [1/5]  Cloning Repository
  ──────────────────────────────────────────────────
  ℹ  Cloning from https://github.com/example/project.git
  $  git clone --progress https://github.com/example/project.git ./project
  Cloning into './project'...
  remote: Enumerating objects: 143, done.
  remote: Counting objects: 100% (143/143), done.
  ✓  Cloned 'project' → /home/user/project

  [2/5]  Detecting Technology Stack
  ──────────────────────────────────────────────────
  ✓  Stack detected: 🟩 Node.js

  [3/5]  Installing Dependencies
  ──────────────────────────────────────────────────
  ℹ  Running: npm install
  $  npm install

  added 342 packages in 8s
  ✓  Dependencies installed

  [4/5]  Analyzing Architecture
  ──────────────────────────────────────────────────
  ✓  Architecture analysis complete

  ┌─────────────────────────────────────────────────┐
  │  📋  Project Architecture: project              │
  └─────────────────────────────────────────────────┘

  Stack:             🟩  Node.js
  Description:       JavaScript / TypeScript project managed by npm
  Detected via:      package.json
  Total files:       87
  Location:          /home/user/project

  Directory Structure:
  ├── src/
  │   ├── components/
  │   ├── pages/
  │   └── utils/
  ├── public/
  ├── package.json
  ├── tsconfig.json
  └── README.md

  File Statistics:
    TypeScript    ████████████████████  41 files
    JSON          ████                   8 files
    Markdown      ██                     4 files
    CSS           ██                     3 files
    YAML          █                      2 files

  Detected Entry Points:
    ▸  src/index.ts

  Config Files:
    ▸  tsconfig.json
    ▸  .env.example
    ▸  Dockerfile

  Install:
    $ npm install
  Run:
    $ npm run dev
    $ npm start

  [5/5]  Starting Development Server
  ──────────────────────────────────────────────────
  ℹ  Starting dev server: npm run dev

  ╔══════════════════════════════════════════╗
  ║  🚀 Dev server starting                   ║
  ╠══════════════════════════════════════════╣
  ║  Command: npm run dev                    ║
  ║  Press Ctrl+C to stop                    ║
  ╚══════════════════════════════════════════╝

  > project@1.0.0 dev
  > vite

    VITE v5.0.0  ready in 324 ms

    ➜  Local:   http://localhost:5173/
    ➜  Network: use --host to expose
```

---

## Supported Stacks

| Stack | Marker File | Install Command | Dev Command |
|---|---|---|---|
| 🟩 Node.js | `package.json` | `npm install` | `npm run dev` / `npm start` |
| 🐍 Python | `requirements.txt` | `pip install -r requirements.txt` | `python main.py` |
| 🐹 Go | `go.mod` | `go mod download` | `go run .` |
| ☕ Java | `pom.xml` | `mvn install -DskipTests` | `mvn spring-boot:run` |
| 💎 Ruby | `Gemfile` | `bundle install` | `bundle exec rails server` |
| 🦀 Rust | `Cargo.toml` | `cargo build` | `cargo run` |

Detection rules are evaluated in order; the first matching marker wins.
Adding a new stack requires a single entry in `internal/detector/detector.go`.

---

## Architecture

```
devsetup/
├── main.go                          Entry point
├── go.mod                           Module definition (zero external deps)
├── Makefile                         Developer workflow
│
├── cmd/devsetup/
│   └── root.go                      CLI flag parsing + 5-step pipeline
│
├── internal/                        Private application packages
│   ├── cloner/
│   │   ├── cloner.go                Git clone + URL normalization
│   │   └── cloner_test.go
│   ├── detector/
│   │   ├── detector.go              Stack detection via marker files
│   │   └── detector_test.go
│   ├── installer/
│   │   └── installer.go             Dependency installation runner
│   ├── runner/
│   │   └── runner.go                Dev server launcher + signal handling
│   ├── analyzer/
│   │   ├── analyzer.go              Architecture analysis + display
│   │   └── analyzer_test.go
│   └── logger/
│       └── logger.go                ANSI colorized structured logger
│
├── pkg/utils/
│   ├── utils.go                     Shared helpers (PATH checks)
│   └── utils_test.go
│
└── .github/workflows/
    ├── ci.yml                       Test matrix (Go 1.21/1.22 × 3 OS)
    └── release.yml                  GoReleaser on tag push
```

### Design Principles

- **Modular** — Each concern lives in its own package with a minimal interface. The pipeline in `root.go` is just five sequential function calls.
- **Extensible** — Adding a new stack is a single `rule{}` struct in `detector.go`. No other code changes.
- **Stdlib-only** — Zero external dependencies. `go build -o devsetup .` works immediately after `git clone`.
- **Fail-fast** — Required tools (git, npm, pip, …) are checked via `exec.LookPath` before any work starts.
- **Signal-safe** — The dev server child process receives `SIGINT`/`SIGTERM` forwarded from the parent so `Ctrl+C` shuts down cleanly without orphaned processes.
- **Cross-platform** — `CGO_ENABLED=0` produces fully static binaries. ANSI colors are suppressed on legacy Windows terminals.

---

## Contributing

Contributions are welcome! Please follow these steps:

1. **Fork** the repository on GitHub.
2. **Clone** your fork: `git clone https://github.com/<you>/devsetup.git`
3. **Create a branch**: `git checkout -b feat/my-feature`
4. **Make your changes** and add/update tests.
5. **Run the test suite**: `make test` — all tests must pass.
6. **Run the linter**: `make lint` — no new warnings.
7. **Commit** with a conventional commit message:
   ```
   feat: add PHP/Composer stack support
   fix: handle repos with no README
   docs: update supported stacks table
   ```
8. **Push** your branch and open a Pull Request against `main`.

### Adding a new stack

Edit `internal/detector/detector.go` and append a rule to `detectionRules`:

```go
{
    file: "composer.json",
    stack: Stack{
        Name:        "PHP",
        Marker:      "composer.json",
        Icon:        "🐘",
        Description: "PHP project managed by Composer",
        InstallCmds: []string{"composer install"},
        DevCmds:     []string{"php -S localhost:8000"},
    },
},
```

Then add the primary tool check to `installer.primaryTool` and a test case
to `internal/detector/detector_test.go`. That's it — nothing else changes.

### Running tests

```bash
make test          # unit tests with -race detector
make test-cover    # tests + HTML coverage report
```

### Reporting bugs

Please open a [GitHub Issue](https://github.com/devsetup/devsetup/issues) with:
- `devsetup version` output
- OS and architecture
- The command you ran
- The full terminal output

---

## License

[MIT](LICENSE) — © 2024 devsetup contributors
