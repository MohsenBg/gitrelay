package ui

import (
	"fmt"
	"io"
	"os"

	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
	"golang.org/x/term"
)

const prefixWidth = 7

type Level int

const (
	LevelError Level = iota
	LevelInfo
	LevelWarn
	LevelSuccess
)

var (
	// Writers can be overridden for testing
	outWriter io.Writer = os.Stdout
	errWriter io.Writer = os.Stderr

	// Styles for different log levels
	styles = map[Level]lipgloss.Style{
		LevelError:   lipgloss.NewStyle().Foreground(lipgloss.Color("9")).Bold(true),
		LevelInfo:    lipgloss.NewStyle().Foreground(lipgloss.Color("12")).Bold(true),
		LevelWarn:    lipgloss.NewStyle().Foreground(lipgloss.Color("11")).Bold(true),
		LevelSuccess: lipgloss.NewStyle().Foreground(lipgloss.Color("10")).Bold(true),
	}

	nameStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("12")).
			Bold(true).
			Width(flagNameWidth)
)

// init configures lipgloss based on terminal capabilities and NO_COLOR env var.
func Init() {
	if os.Getenv("NO_COLOR") != "" || !term.IsTerminal(int(os.Stdout.Fd())) {
		lipgloss.SetColorProfile(termenv.Ascii)
	}
}

// Error prints an error message to stderr.
func Error(format string, a ...any) {
	printMessage(LevelError, format, a...)
}

// Info prints an informational message to stdout.
func Info(format string, a ...any) {
	printMessage(LevelInfo, format, a...)
}

// Warning prints a warning message to stderr.
func Warning(format string, a ...any) {
	printMessage(LevelWarn, format, a...)
}

// Success prints a success message to stdout.
func Success(format string, a ...any) {
	printMessage(LevelSuccess, format, a...)
}

// printMessage handles the actual printing logic for different levels.
func printMessage(level Level, format string, a ...any) {
	prefix := levelPrefix(level)

	prefixStyled := styles[level].
		PaddingRight(max(0, prefixWidth-len(prefix))).
		Render(prefix)

	msg := fmt.Sprintf(format, a...)

	w := outWriter
	if level == LevelError || level == LevelWarn {
		w = errWriter
	}

	fmt.Fprintf(w, "%s %s\n", prefixStyled, msg)
}

// levelPrefix returns the string prefix for a given log level.
func levelPrefix(level Level) string {
	switch level {
	case LevelError:
		return "ERROR"
	case LevelInfo:
		return "INFO"
	case LevelWarn:
		return "WARN"
	case LevelSuccess:
		return "OK"
	default:
		return "LOG"
	}
}

// SetWriters allows overriding the default stdout/stderr writers, useful for testing.
func SetWriters(out, err io.Writer) {
	outWriter = out
	errWriter = err
}
