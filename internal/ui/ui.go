package ui

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"golang.org/x/term"
)

const appName = "eggl"

type Theme struct {
	enabled bool

	accent     lipgloss.Color
	success    lipgloss.Color
	errorColor lipgloss.Color
	mutedColor lipgloss.Color
	text       lipgloss.Color

	title    lipgloss.Style
	subtitle lipgloss.Style
	box      lipgloss.Style
	label    lipgloss.Style
	value    lipgloss.Style
	ok       lipgloss.Style
	fail     lipgloss.Style
	muted    lipgloss.Style
	command  lipgloss.Style
}

func NewTheme(w io.Writer) Theme {
	enabled := IsInteractive(w)

	t := Theme{
		enabled:    enabled,
		accent:     lipgloss.Color("#7C3AED"),
		success:    lipgloss.Color("#22C55E"),
		errorColor: lipgloss.Color("#EF4444"),
		mutedColor: lipgloss.Color("#6B7280"),
		text:       lipgloss.Color("#E5E7EB"),
	}

	if !enabled {
		return t
	}

	t.title = lipgloss.NewStyle().
		Bold(true).
		Foreground(t.accent).
		MarginBottom(1)

	t.subtitle = lipgloss.NewStyle().
		Foreground(t.mutedColor).
		MarginBottom(1)

	t.box = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(t.accent).
		Padding(0, 2).
		MarginBottom(1)

	t.label = lipgloss.NewStyle().
		Foreground(t.mutedColor).
		Width(8)

	t.value = lipgloss.NewStyle().
		Foreground(t.text)

	t.ok = lipgloss.NewStyle().
		Foreground(t.success).
		Bold(true)

	t.fail = lipgloss.NewStyle().
		Foreground(t.errorColor).
		Bold(true)

	t.muted = lipgloss.NewStyle().
		Foreground(t.mutedColor)

	t.command = lipgloss.NewStyle().
		Foreground(t.accent).
		Bold(true)

	return t
}

func IsInteractive(w io.Writer) bool {
	if os.Getenv("NO_COLOR") != "" {
		return false
	}

	f, ok := w.(*os.File)
	if !ok {
		return false
	}

	return term.IsTerminal(int(f.Fd()))
}

func RenderHeader(w io.Writer, title, subtitle string) {
	t := NewTheme(w)
	if !t.enabled {
		fmt.Fprintf(w, "%s\n", title)
		if subtitle != "" {
			fmt.Fprintf(w, "%s\n\n", subtitle)
		}
		return
	}

	fmt.Fprintln(w, t.title.Render(title))
	if subtitle != "" {
		fmt.Fprintln(w, t.subtitle.Render(subtitle))
	}
}

type VersionInfo struct {
	Version string
	Commit  string
	Date    string
}

func RenderVersion(w io.Writer, info VersionInfo) {
	t := NewTheme(w)
	if !t.enabled {
		fmt.Fprintf(w, "eggl version %s\n", info.Version)
		fmt.Fprintf(w, "commit: %s\n", info.Commit)
		fmt.Fprintf(w, "built:  %s\n", info.Date)
		return
	}

	lines := []string{
		renderKV(t, "version", info.Version),
		renderKV(t, "commit", info.Commit),
		renderKV(t, "built", info.Date),
	}

	body := strings.Join(lines, "\n")
	fmt.Fprintln(w, t.box.Render(t.title.Render(appName)+"\n"+body))
}

func renderKV(t Theme, key, value string) string {
	return t.label.Render(key) + t.value.Render(value)
}

type DoctorCheck struct {
	Name   string
	Status string
	Detail string
	OK     bool
}

func RenderDoctor(w io.Writer, checks []DoctorCheck) {
	t := NewTheme(w)
	if !t.enabled {
		for _, check := range checks {
			marker := "ok"
			if !check.OK {
				marker = "fail"
			}
			fmt.Fprintf(w, "[%s] %s: %s (%s)\n", marker, check.Name, check.Status, check.Detail)
		}
		renderDoctorSummary(w, t, checks)
		return
	}

	RenderHeader(w, "eggl doctor", "Environment and dependency checks")

	rows := make([]string, 0, len(checks))
	for _, check := range checks {
		icon := t.ok.Render("✓")
		statusStyle := t.value
		if !check.OK {
			icon = t.fail.Render("✗")
			statusStyle = t.fail
		}

		row := fmt.Sprintf("%s  %-8s %s",
			icon,
			t.command.Render(check.Name),
			statusStyle.Render(check.Status),
		)
		rows = append(rows, row)
		rows = append(rows, t.muted.Render("   "+check.Detail))
	}

	fmt.Fprintln(w, t.box.Render(strings.Join(rows, "\n")))
	renderDoctorSummary(w, t, checks)
}

func renderDoctorSummary(w io.Writer, t Theme, checks []DoctorCheck) {
	failures := 0
	for _, check := range checks {
		if !check.OK {
			failures++
		}
	}

	switch {
	case failures == 0:
		if t.enabled {
			fmt.Fprintln(w, t.ok.Render("All checks passed"))
		} else {
			fmt.Fprintln(w, "All checks passed")
		}
	default:
		msg := fmt.Sprintf("%d check(s) failed", failures)
		if t.enabled {
			fmt.Fprintln(w, t.fail.Render(msg))
		} else {
			fmt.Fprintln(w, msg)
		}
	}
}

type HelpCommand struct {
	Name        string
	Description string
}

func RenderHelp(w io.Writer, summary string, commands []HelpCommand, globalFlags *pflag.FlagSet) {
	t := NewTheme(w)
	sections := helpSectionsForCommands(commands, globalFlags)

	if !t.enabled {
		fmt.Fprintf(w, "%s\n\n", summary)
		writeHelpSectionsPlain(w, sections)
		return
	}

	RenderHeader(w, appName, summary)
	fmt.Fprintln(w, t.box.Render(renderHelpSectionsStyled(t, sections)))
	fmt.Fprintln(w, t.muted.Render("Tip: run `eggl <command> --help` for details"))
}

func RenderCommandHelp(w io.Writer, cmd *cobra.Command) {
	t := NewTheme(w)
	title := cmd.CommandPath()
	subtitle := commandSubtitle(cmd)
	sections := helpSectionsForCommand(cmd)

	if !t.enabled {
		fmt.Fprintf(w, "%s — %s\n\n", title, subtitle)
		writeHelpSectionsPlain(w, sections)
		return
	}

	RenderHeader(w, title, subtitle)
	if len(sections) > 0 {
		fmt.Fprintln(w, t.box.Render(renderHelpSectionsStyled(t, sections)))
	}
}

type helpSection struct {
	title string
	lines []string
}

func commandSubtitle(cmd *cobra.Command) string {
	if cmd.Short != "" {
		return cmd.Short
	}
	return firstLine(cmd.Long)
}

func firstLine(s string) string {
	if s == "" {
		return ""
	}
	line, _, _ := strings.Cut(strings.TrimSpace(s), "\n")
	return line
}

func commandUsage(cmd *cobra.Command) string {
	if cmd.DisableFlagsInUseLine {
		return cmd.Use
	}
	return cmd.UseLine()
}

func helpSectionsForCommands(commands []HelpCommand, globalFlags *pflag.FlagSet) []helpSection {
	sections := make([]helpSection, 0, 2)

	if len(commands) > 0 {
		lines := make([]string, 0, len(commands))
		for _, command := range commands {
			lines = append(lines, fmt.Sprintf("  %-12s %s", command.Name, command.Description))
		}
		sections = append(sections, helpSection{title: "Available Commands", lines: lines})
	}

	if flagLines := collectFlagLines(globalFlags); len(flagLines) > 0 {
		sections = append(sections, helpSection{title: "Global Flags", lines: flagLines})
	}

	return sections
}

func helpSectionsForCommand(cmd *cobra.Command) []helpSection {
	sections := make([]helpSection, 0, 4)

	if usage := commandUsage(cmd); usage != "" {
		sections = append(sections, helpSection{
			title: "Usage",
			lines: []string{"  " + usage},
		})
	}

	if len(cmd.ValidArgs) > 0 {
		sections = append(sections, helpSection{
			title: "Valid Args",
			lines: []string{"  " + strings.Join(cmd.ValidArgs, ", ")},
		})
	}

	if flagLines := collectFlagLines(cmd.NonInheritedFlags()); len(flagLines) > 0 {
		sections = append(sections, helpSection{title: "Flags", lines: flagLines})
	}

	if flagLines := collectFlagLines(cmd.InheritedFlags()); len(flagLines) > 0 {
		sections = append(sections, helpSection{title: "Global Flags", lines: flagLines})
	}

	if example := strings.TrimSpace(cmd.Example); example != "" {
		lines := strings.Split(example, "\n")
		for i, line := range lines {
			lines[i] = "  " + line
		}
		sections = append(sections, helpSection{title: "Examples", lines: lines})
	}

	return sections
}

func collectFlagLines(flags *pflag.FlagSet) []string {
	if flags == nil {
		return nil
	}

	lines := make([]string, 0)
	flags.VisitAll(func(f *pflag.Flag) {
		if f.Hidden || f.Name == "help" {
			return
		}
		lines = append(lines, formatFlagLine(f))
	})
	return lines
}

func formatFlagLine(f *pflag.Flag) string {
	name := flagNames(f)
	if shouldShowDefault(f) {
		name = name + " " + f.DefValue
	}
	return fmt.Sprintf("  %-22s %s", name, f.Usage)
}

func flagNames(f *pflag.Flag) string {
	if f.Shorthand != "" && len(f.Shorthand) == 1 {
		return fmt.Sprintf("-%s, --%s", f.Shorthand, f.Name)
	}
	return "--" + f.Name
}

func shouldShowDefault(f *pflag.Flag) bool {
	if f.DefValue == "" {
		return false
	}
	switch f.Value.Type() {
	case "bool":
		return false
	default:
		return true
	}
}

func writeHelpSectionsPlain(w io.Writer, sections []helpSection) {
	for _, section := range sections {
		fmt.Fprintf(w, "%s:\n", section.title)
		for _, line := range section.lines {
			if strings.HasPrefix(line, "  ") {
				fmt.Fprintln(w, line)
			} else {
				fmt.Fprintf(w, "  %s\n", line)
			}
		}
		fmt.Fprintln(w)
	}
}

func renderHelpSectionsStyled(t Theme, sections []helpSection) string {
	parts := make([]string, 0, len(sections))
	for i, section := range sections {
		if i > 0 {
			parts = append(parts, "")
		}
		parts = append(parts, t.command.Render(section.title))
		for _, line := range section.lines {
			parts = append(parts, renderHelpLineStyled(t, line))
		}
	}
	return strings.Join(parts, "\n")
}

func renderHelpLineStyled(t Theme, line string) string {
	trimmed := strings.TrimPrefix(line, "  ")
	if idx := strings.Index(trimmed, "  "); idx > 0 {
		left := trimmed[:idx]
		right := strings.TrimSpace(trimmed[idx:])
		return "  " + t.command.Render(left) + "  " + t.muted.Render(right)
	}
	return "  " + t.muted.Render(trimmed)
}
