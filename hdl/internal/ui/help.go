package ui

import (
	"flag"
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	helpTitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("13"))

	flagNameStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("10"))

	defaultStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("12"))

	usageStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("7"))

	flagNameWidth = 20
	indentUsage   = 4
)

func PrintDefaults() {
	// Title
	fmt.Println(helpTitleStyle.Render("Flags"))
	// Directly print flags, no border/panel
	fmt.Println(renderFlags())
}

func renderFlags() string {
	var b strings.Builder

	first := true
	flag.CommandLine.VisitAll(func(f *flag.Flag) {
		if !first {
			b.WriteString("\n")
		}
		first = false

		name := fmt.Sprintf("--%s", f.Name)
		namePadded := flagNameStyle.Width(flagNameWidth).Render(name)

		def := fmt.Sprintf("(default: %s)", formatDefault(f.DefValue))
		line1 := fmt.Sprintf("%s %s", namePadded, defaultStyle.Render(def))

		usage := strings.TrimSpace(f.Usage)
		if usage == "" {
			usage = "-"
		}

		usageBlock := usageStyle.Render(indentText(usage, indentUsage))

		b.WriteString("  " + line1 + "\n  " + usageBlock)
	})

	if b.Len() == 0 {
		b.WriteString("  (no flags)")
	}

	return b.String()
}

func indentText(s string, n int) string {
	prefix := strings.Repeat(" ", n)
	s = strings.ReplaceAll(s, "\n", "\n"+prefix)
	return prefix + s
}

func formatDefault(def string) string {
	if def == "" {
		return `""`
	}
	return def
}

