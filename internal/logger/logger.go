// Package logger provides colorized, structured terminal output using
// only ANSI escape codes — no external dependencies required.
package logger

import (
	"fmt"
	"io"
	"os"
	"runtime"
	"time"
)

// ANSI color codes
const (
	reset     = "\033[0m"
	bold      = "\033[1m"
	faint     = "\033[2m"
	red       = "\033[31m"
	green     = "\033[32m"
	yellow    = "\033[33m"
	blue      = "\033[34m"
	magenta   = "\033[35m"
	cyan      = "\033[36m"
	white     = "\033[37m"
	boldGreen = "\033[1;32m"
	boldCyan  = "\033[1;36m"
	boldRed   = "\033[1;31m"
)

// Logger is a structured, colorized console logger.
type Logger struct {
	out     io.Writer
	verbose bool
	noColor bool
}

// New creates a new Logger. Color is disabled on Windows without a VT100 terminal.
func New(verbose bool) *Logger {
	noColor := runtime.GOOS == "windows" && os.Getenv("TERM") == ""
	return &Logger{
		out:     os.Stdout,
		verbose: verbose,
		noColor: noColor,
	}
}

func (l *Logger) colorize(code, text string) string {
	if l.noColor {
		return text
	}
	return code + text + reset
}

// Banner prints the devsetup ASCII welcome banner.
func (l *Logger) Banner() {
	fmt.Fprintln(l.out)
	fmt.Fprintln(l.out, l.colorize(boldCyan, "  ╔══════════════════════════════════════════════════╗"))
	fmt.Fprintln(l.out, l.colorize(boldCyan, "  ║")+l.colorize(bold, "   ██████╗ ███████╗██╗   ██╗███████╗███████╗████████╗██╗   ██╗██████╗  ")+l.colorize(boldCyan, "║"))
	fmt.Fprintln(l.out, l.colorize(boldCyan, "  ║")+l.colorize(bold, "   ██╔══██╗██╔════╝██║   ██║██╔════╝██╔════╝╚══██╔══╝██║   ██║██╔══██╗ ")+l.colorize(boldCyan, "║"))
	fmt.Fprintln(l.out, l.colorize(boldCyan, "  ║")+l.colorize(bold, "   ██║  ██║█████╗  ██║   ██║███████╗█████╗     ██║   ██║   ██║██████╔╝ ")+l.colorize(boldCyan, "║"))
	fmt.Fprintln(l.out, l.colorize(boldCyan, "  ║")+l.colorize(bold, "   ██║  ██║██╔══╝  ╚██╗ ██╔╝╚════██║██╔══╝     ██║   ██║   ██║██╔═══╝  ")+l.colorize(boldCyan, "║"))
	fmt.Fprintln(l.out, l.colorize(boldCyan, "  ║")+l.colorize(bold, "   ██████╔╝███████╗ ╚████╔╝ ███████║███████╗   ██║   ╚██████╔╝██║      ")+l.colorize(boldCyan, "║"))
	fmt.Fprintln(l.out, l.colorize(boldCyan, "  ║")+l.colorize(faint, "   ╚═════╝ ╚══════╝  ╚═══╝  ╚══════╝╚══════╝   ╚═╝    ╚═════╝ ╚═╝      ")+l.colorize(boldCyan, "║"))
	fmt.Fprintln(l.out, l.colorize(boldCyan, "  ╠══════════════════════════════════════════════════╣"))
	fmt.Fprintf(l.out, "  %s  Instant Dev Environment Setup  •  %s%s\n",
		l.colorize(boldCyan, "║"),
		l.colorize(faint, time.Now().Format("2006-01-02 15:04:05")),
		l.colorize(boldCyan, "          ║"))
	fmt.Fprintln(l.out, l.colorize(boldCyan, "  ╚══════════════════════════════════════════════════╝"))
	fmt.Fprintln(l.out)
}

// Step prints a numbered pipeline step header.
func (l *Logger) Step(current, total int, msg string) {
	fmt.Fprintln(l.out)
	fmt.Fprintf(l.out, "  %s  %s\n",
		l.colorize(boldCyan, fmt.Sprintf("[%d/%d]", current, total)),
		l.colorize(bold, msg),
	)
	fmt.Fprintln(l.out, l.colorize(faint, "  ──────────────────────────────────────────────────"))
}

// Info prints an informational message.
func (l *Logger) Info(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	fmt.Fprintf(l.out, "  %s  %s\n", l.colorize(blue, "ℹ"), msg)
}

// Success prints a success message.
func (l *Logger) Success(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	fmt.Fprintf(l.out, "  %s  %s\n", l.colorize(boldGreen, "✓"), l.colorize(boldGreen, msg))
}

// Warn prints a warning message.
func (l *Logger) Warn(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	fmt.Fprintf(l.out, "  %s  %s\n", l.colorize(yellow, "⚠"), l.colorize(yellow, msg))
}

// Error prints an error message.
func (l *Logger) Error(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	fmt.Fprintf(l.out, "  %s  %s\n", l.colorize(boldRed, "✗"), l.colorize(boldRed, msg))
}

// Debug prints a verbose/debug message (only shown when verbose=true).
func (l *Logger) Debug(format string, args ...interface{}) {
	if l.verbose {
		msg := fmt.Sprintf(format, args...)
		fmt.Fprintf(l.out, "  %s  %s\n", l.colorize(faint, "·"), l.colorize(faint, "[debug] "+msg))
	}
}

// Command prints the shell command being executed.
func (l *Logger) Command(cmd string) {
	fmt.Fprintf(l.out, "  %s %s\n", l.colorize(faint, "$"), l.colorize(faint, cmd))
}

// Print writes a raw formatted string.
func (l *Logger) Print(format string, args ...interface{}) {
	fmt.Fprintf(l.out, format, args...)
}

// Println writes a raw line.
func (l *Logger) Println(s string) {
	fmt.Fprintln(l.out, s)
}

// Colorize exposes colorization for use by other packages.
func (l *Logger) Colorize(code, text string) string {
	return l.colorize(code, text)
}

// Color constants exposed for use by sibling packages.
const (
	ColorReset     = reset
	ColorBold      = bold
	ColorFaint     = faint
	ColorRed       = red
	ColorGreen     = green
	ColorYellow    = yellow
	ColorBlue      = blue
	ColorMagenta   = magenta
	ColorCyan      = cyan
	ColorBoldGreen = boldGreen
	ColorBoldCyan  = boldCyan
	ColorBoldRed   = boldRed
)
