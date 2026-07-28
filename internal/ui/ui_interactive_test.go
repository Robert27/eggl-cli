package ui

import (
	"io"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/creack/pty"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func openTestTTY(t *testing.T) (io.Writer, *os.File) {
	t.Helper()

	master, tty, err := pty.Open()
	if err != nil {
		t.Skip("pty:", err)
	}
	t.Cleanup(func() {
		_ = master.Close()
		_ = tty.Close()
	})
	return tty, master
}

func readPTY(t *testing.T, master *os.File) string {
	t.Helper()

	deadline := time.Now().Add(200 * time.Millisecond)
	_ = master.SetReadDeadline(deadline)

	buf := make([]byte, 8192)
	n, err := master.Read(buf)
	if err != nil && n == 0 {
		t.Fatalf("read pty: %v", err)
	}
	return string(buf[:n])
}

func renderPTY(t *testing.T, master *os.File, slave io.Writer, render func(io.Writer)) string {
	t.Helper()

	out := make(chan string, 1)
	go func() {
		data, _ := io.ReadAll(master)
		out <- string(data)
	}()

	render(slave)

	if f, ok := slave.(*os.File); ok {
		_ = f.Close()
	}

	select {
	case got := <-out:
		return got
	case <-time.After(5 * time.Second):
		t.Fatal("timeout reading pty output")
		return ""
	}
}

func TestRenderVersionInteractive(t *testing.T) {
	t.Setenv("NO_COLOR", "")

	w, master := openTestTTY(t)
	if !IsInteractive(w) {
		t.Fatal("expected tty writer to be interactive")
	}

	RenderVersion(w, VersionInfo{
		Version: "v1.0.0",
		Commit:  "abc",
		Date:    "2026-01-01",
	})

	got := readPTY(t, master)
	for _, want := range []string{"eggl", "v1.0.0", "abc", "2026-01-01"} {
		if !strings.Contains(got, want) {
			t.Fatalf("expected %q in styled output, got %q", want, got)
		}
	}
}

func TestRenderHelpInteractive(t *testing.T) {
	t.Setenv("NO_COLOR", "")

	w, master := openTestTTY(t)
	flags := pflag.NewFlagSet("test", pflag.ContinueOnError)
	flags.BoolP("verbose", "v", false, "Log details")

	got := renderPTY(t, master, w, func(out io.Writer) {
		RenderHelp(out, "Helper CLI", "", []HelpCommand{
			{Name: "doctor", Description: "Checks"},
		}, flags)
	})
	for _, want := range []string{"eggl", "Available Commands", "doctor", "Global Flags"} {
		if !strings.Contains(got, want) {
			t.Fatalf("expected %q in styled help, got %q", want, got)
		}
	}
}

func TestRenderDoctorInteractive(t *testing.T) {
	t.Setenv("NO_COLOR", "")

	w, master := openTestTTY(t)
	got := renderPTY(t, master, w, func(out io.Writer) {
		RenderDoctor(out, []DoctorCheck{
			{Name: "go", Status: "ok", Detail: "runtime", OK: true},
		})
	})

	for _, want := range []string{"eggl doctor", "go", "All checks passed"} {
		if !strings.Contains(got, want) {
			t.Fatalf("expected %q in styled doctor output, got %q", want, got)
		}
	}
}

func TestRenderDedashInteractiveModified(t *testing.T) {
	t.Setenv("NO_COLOR", "")

	w, master := openTestTTY(t)
	got := renderPTY(t, master, w, func(out io.Writer) {
		RenderDedash(out, DedashSummary{
			Scanned:           2,
			Modified:          1,
			TotalReplacements: 1,
			Changes:           []DedashChange{{Path: "a.md", Replacements: 1}},
		})
	})
	for _, want := range []string{"modified 1", "a.md"} {
		if !strings.Contains(got, want) {
			t.Fatalf("expected %q in styled dedash output, got %q", want, got)
		}
	}
}

func TestRenderEOLInteractiveModified(t *testing.T) {
	t.Setenv("NO_COLOR", "")

	w, master := openTestTTY(t)
	got := renderPTY(t, master, w, func(out io.Writer) {
		RenderEOL(out, EOLSummary{
			Scanned:           2,
			Modified:          1,
			TotalReplacements: 2,
			Changes:           []EOLChange{{Path: "a.md", Replacements: 2}},
		})
	})
	for _, want := range []string{"modified 1", "a.md"} {
		if !strings.Contains(got, want) {
			t.Fatalf("expected %q in styled eol output, got %q", want, got)
		}
	}
}

func TestRenderCommandHelpInteractive(t *testing.T) {
	t.Setenv("NO_COLOR", "")

	w, master := openTestTTY(t)
	cmd := &cobra.Command{
		Use:   "dedash",
		Short: "Replace em-dashes",
		Long:  "Long help text.",
	}
	cmd.Flags().String("path", ".", "scan root")

	got := renderPTY(t, master, w, func(out io.Writer) {
		RenderCommandHelp(out, cmd)
	})
	for _, want := range []string{"dedash", "Usage", "Flags", "--path"} {
		if !strings.Contains(got, want) {
			t.Fatalf("expected %q in styled command help, got %q", want, got)
		}
	}
}
