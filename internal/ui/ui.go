package ui

import (
	"bufio"
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
	heading  lipgloss.Style
}

func NewTheme(w io.Writer) Theme {
	enabled := IsInteractive(w)

	t := Theme{
		enabled:    enabled,
		accent:     lipgloss.Color("#F97316"),
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

	t.heading = lipgloss.NewStyle().
		Foreground(t.text).
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

func IsInteractiveInput(r io.Reader) bool {
	f, ok := r.(*os.File)
	if !ok {
		return false
	}
	return term.IsTerminal(int(f.Fd()))
}

func ConfirmPrompt(w io.Writer, in io.Reader, prompt string) (bool, error) {
	if _, err := fmt.Fprint(w, prompt); err != nil {
		return false, err
	}

	line, err := bufio.NewReader(in).ReadString('\n')
	if err != nil {
		if err == io.EOF && strings.TrimSpace(line) == "" {
			return false, err
		}
		if err != io.EOF {
			return false, err
		}
	}

	line = strings.TrimSpace(strings.ToLower(line))
	return line == "y" || line == "yes", nil
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

type DedashChange struct {
	Path         string
	Replacements int
}

type DedashSummary struct {
	Scanned           int
	Modified          int
	Skipped           int
	TotalReplacements int
	Changes           []DedashChange
	DryRun            bool
}

func RenderDedash(w io.Writer, summary DedashSummary) {
	t := NewTheme(w)
	line := dedashSummaryLine(summary)

	switch {
	case !t.enabled:
		fmt.Fprintln(w, line)
	case summary.Modified > 0 && !summary.DryRun:
		fmt.Fprintln(w, t.ok.Render(line))
	case summary.DryRun:
		fmt.Fprintln(w, t.muted.Render(line))
	default:
		fmt.Fprintln(w, t.muted.Render(line))
	}

	for _, change := range summary.Changes {
		if t.enabled {
			fmt.Fprintf(w, "  %s %s\n",
				t.command.Render(change.Path),
				t.muted.Render(fmt.Sprintf("(%d)", change.Replacements)),
			)
			continue
		}
		fmt.Fprintf(w, "  %s (%d)\n", change.Path, change.Replacements)
	}
}

func dedashSummaryLine(summary DedashSummary) string {
	skipped := ""
	if summary.Skipped > 0 {
		skipped = fmt.Sprintf(", skipped %d", summary.Skipped)
	}

	if summary.Modified == 0 {
		return fmt.Sprintf("scanned %d files%s, no em-dashes found", summary.Scanned, skipped)
	}

	if summary.DryRun {
		return fmt.Sprintf("dry-run: scanned %d files%s, would modify %d (%d replacements)",
			summary.Scanned, skipped, summary.Modified, summary.TotalReplacements)
	}

	return fmt.Sprintf("scanned %d files%s, modified %d (%d replacements)",
		summary.Scanned, skipped, summary.Modified, summary.TotalReplacements)
}

type EOLChange struct {
	Path         string
	Replacements int
}

type EOLSummary struct {
	Scanned           int
	Modified          int
	Skipped           int
	TotalReplacements int
	Changes           []EOLChange
	DryRun            bool
}

func RenderEOL(w io.Writer, summary EOLSummary) {
	t := NewTheme(w)
	line := eolSummaryLine(summary)

	switch {
	case !t.enabled:
		fmt.Fprintln(w, line)
	case summary.Modified > 0 && !summary.DryRun:
		fmt.Fprintln(w, t.ok.Render(line))
	case summary.DryRun:
		fmt.Fprintln(w, t.muted.Render(line))
	default:
		fmt.Fprintln(w, t.muted.Render(line))
	}

	for _, change := range summary.Changes {
		if t.enabled {
			fmt.Fprintf(w, "  %s %s\n",
				t.command.Render(change.Path),
				t.muted.Render(fmt.Sprintf("(%d)", change.Replacements)),
			)
			continue
		}
		fmt.Fprintf(w, "  %s (%d)\n", change.Path, change.Replacements)
	}
}

func eolSummaryLine(summary EOLSummary) string {
	skipped := ""
	if summary.Skipped > 0 {
		skipped = fmt.Sprintf(", skipped %d", summary.Skipped)
	}

	if summary.Modified == 0 {
		return fmt.Sprintf("scanned %d files%s, no line ending fixes needed", summary.Scanned, skipped)
	}

	if summary.DryRun {
		return fmt.Sprintf("dry-run: scanned %d files%s, would modify %d (%d line endings)",
			summary.Scanned, skipped, summary.Modified, summary.TotalReplacements)
	}

	return fmt.Sprintf("scanned %d files%s, modified %d (%d line endings)",
		summary.Scanned, skipped, summary.Modified, summary.TotalReplacements)
}

type EnvProfile struct {
	Name             string
	KubeContext      string
	TailscaleAccount string
}

type EnvShowReport struct {
	ActiveProfile string
	Unknown       bool
	KubeContext   string
	Tailscale     string
	ConfigPath    string
	Profiles      []EnvProfile
}

type EnvSwitchResult struct {
	FromProfile string
	ToProfile   string
	FromKube    string
	ToKube      string
	FromTS      string
	ToTS        string
}

func RenderEnvShow(w io.Writer, report EnvShowReport) {
	t := NewTheme(w)
	if !t.enabled {
		if report.Unknown {
			fmt.Fprintf(w, "profile: unknown\n")
		} else {
			fmt.Fprintf(w, "profile: %s\n", report.ActiveProfile)
		}
		fmt.Fprintf(w, "kube: %s\n", report.KubeContext)
		fmt.Fprintf(w, "tailscale: %s\n", report.Tailscale)
		fmt.Fprintf(w, "config: %s\n", report.ConfigPath)
		for _, p := range report.Profiles {
			fmt.Fprintf(w, "  %s: kube=%s tailscale=%s\n", p.Name, p.KubeContext, p.TailscaleAccount)
		}
		return
	}

	RenderHeader(w, "eggl env", "Environment profile status")

	profileVal := report.ActiveProfile
	if report.Unknown {
		profileVal = "unknown"
	}
	lines := []string{
		renderKV(t, "profile", profileVal),
		renderKV(t, "kube", report.KubeContext),
		renderKV(t, "tailscale", report.Tailscale),
		renderKV(t, "config", report.ConfigPath),
	}
	fmt.Fprintln(w, t.box.Render(strings.Join(lines, "\n")))

	if len(report.Profiles) > 0 {
		fmt.Fprintln(w, t.muted.Render("Configured profiles:"))
		for _, p := range report.Profiles {
			fmt.Fprintf(w, "  %s %s\n",
				t.command.Render(p.Name),
				t.muted.Render(fmt.Sprintf("kube=%s tailscale=%s", p.KubeContext, p.TailscaleAccount)),
			)
		}
	}
}

func RenderEnvSwitch(w io.Writer, result EnvSwitchResult) {
	t := NewTheme(w)
	if !t.enabled {
		fmt.Fprintf(w, "profile: %s → %s\n", result.FromProfile, result.ToProfile)
		fmt.Fprintf(w, "kube: %s → %s\n", result.FromKube, result.ToKube)
		fmt.Fprintf(w, "tailscale: %s → %s\n", result.FromTS, result.ToTS)
		return
	}

	RenderHeader(w, "eggl env", "Switched environment profile")
	lines := []string{
		renderKV(t, "profile", fmt.Sprintf("%s → %s", emptyDash(result.FromProfile), result.ToProfile)),
		renderKV(t, "kube", fmt.Sprintf("%s → %s", result.FromKube, result.ToKube)),
		renderKV(t, "tailscale", fmt.Sprintf("%s → %s", result.FromTS, result.ToTS)),
	}
	fmt.Fprintln(w, t.ok.Render(t.box.Render(strings.Join(lines, "\n"))))
}

func emptyDash(s string) string {
	if s == "" {
		return "—"
	}
	return s
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

func RenderHelp(w io.Writer, summary, description string, commands []HelpCommand, globalFlags *pflag.FlagSet) {
	t := NewTheme(w)
	sections := helpSectionsForCommands(commands, globalFlags)
	if description != "" {
		lines := strings.Split(description, "\n")
		for i, line := range lines {
			lines[i] = "  " + line
		}
		sections = append([]helpSection{{title: "Description", lines: lines}}, sections...)
	}

	if !t.enabled {
		fmt.Fprintf(w, "%s\n\n", summary)
		writeHelpSectionsPlain(w, sections)
		return
	}

	RenderHeader(w, appName, summary)
	writeHelpSectionsStyled(w, t, sections)
	fmt.Fprintln(w, t.muted.Render("Tip: run `eggl <command> --help` for details"))
}

func RenderCommandHelp(w io.Writer, cmd *cobra.Command) {
	t := NewTheme(w)
	title := cmd.CommandPath()
	subtitle := commandSubtitle(cmd)
	sections := helpSectionsForCommand(cmd)

	if !t.enabled {
		fmt.Fprintf(w, "%s - %s\n\n", title, subtitle)
		writeHelpSectionsPlain(w, sections)
		return
	}

	RenderHeader(w, title, subtitle)
	writeHelpSectionsStyled(w, t, sections)
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

func CommandDescription(cmd *cobra.Command) string {
	return commandDescription(cmd)
}

func commandDescription(cmd *cobra.Command) string {
	long := strings.TrimSpace(cmd.Long)
	if long == "" {
		return ""
	}
	short := strings.TrimSpace(cmd.Short)
	if short != "" && firstLine(long) == short {
		rest, _, found := strings.Cut(long, "\n")
		if !found || strings.TrimSpace(rest) == "" {
			return ""
		}
		_, remainder, _ := strings.Cut(long, "\n")
		return strings.TrimSpace(remainder)
	}
	return long
}

func commandUsage(cmd *cobra.Command) string {
	if cmd.DisableFlagsInUseLine {
		use := cmd.Use
		if idx := strings.Index(use, " "); idx > 0 {
			return cmd.CommandPath() + use[idx:]
		}
		return cmd.CommandPath()
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

	if desc := commandDescription(cmd); desc != "" {
		lines := strings.Split(desc, "\n")
		for i, line := range lines {
			lines[i] = "  " + line
		}
		sections = append(sections, helpSection{title: "Description", lines: lines})
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
		lines := make([]string, 0)
		for _, line := range strings.Split(example, "\n") {
			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}
			lines = append(lines, "  "+line)
		}
		if len(lines) > 0 {
			sections = append(sections, helpSection{title: "Examples", lines: lines})
		}
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

func helpContentWidth(w io.Writer) int {
	const (
		defaultWidth = 72
		minWidth     = 40
		maxWidth     = 100
	)
	f, ok := w.(*os.File)
	if !ok {
		return defaultWidth
	}
	width, _, err := term.GetSize(int(f.Fd()))
	if err != nil || width < minWidth {
		return defaultWidth
	}
	width -= 8
	if width > maxWidth {
		width = maxWidth
	}
	if width < minWidth {
		width = minWidth
	}
	return width
}

func writeHelpSectionsStyled(w io.Writer, t Theme, sections []helpSection) {
	width := helpContentWidth(w)
	box := t.box
	if width > 0 {
		box = box.Width(width)
	}

	blocks := make([]string, 0, len(sections))
	for _, section := range sections {
		lines := []string{t.heading.Render(section.title)}
		for _, line := range section.lines {
			lines = append(lines, renderHelpLineStyled(t, line))
		}
		blocks = append(blocks, box.Render(strings.Join(lines, "\n")))
	}
	fmt.Fprintln(w, strings.Join(blocks, "\n"))
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
