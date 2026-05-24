package netbird

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const sampleList = `Found 3 profiles:
✓ work
✗ default
✗ home
`

func TestParseProfileList(t *testing.T) {
	profiles, err := ParseProfileList([]byte(sampleList))
	if err != nil {
		t.Fatalf("ParseProfileList() error = %v", err)
	}
	if len(profiles) != 3 {
		t.Fatalf("len = %d, want 3", len(profiles))
	}
	if profiles[0].Name != "work" || !profiles[0].Active {
		t.Fatalf("profiles[0] = %+v", profiles[0])
	}
	if profiles[1].Name != "default" || profiles[1].Active {
		t.Fatalf("profiles[1] = %+v", profiles[1])
	}
}

func TestResolveProfile(t *testing.T) {
	profiles, err := ParseProfileList([]byte(sampleList))
	if err != nil {
		t.Fatal(err)
	}

	p, err := ResolveProfile("HOME", profiles)
	if err != nil {
		t.Fatalf("ResolveProfile() error = %v", err)
	}
	if p.Name != "home" {
		t.Fatalf("name = %q, want home", p.Name)
	}
}

func TestCurrentProfile(t *testing.T) {
	profiles, err := ParseProfileList([]byte(sampleList))
	if err != nil {
		t.Fatal(err)
	}

	p, err := CurrentProfile(profiles)
	if err != nil {
		t.Fatalf("CurrentProfile() error = %v", err)
	}
	if p.Name != "work" {
		t.Fatalf("name = %q, want work", p.Name)
	}
}

func writeFakeNetbird(t *testing.T, dir string, listOutput string) string {
	t.Helper()

	listPath := filepath.Join(dir, "profiles.txt")
	if err := os.WriteFile(listPath, []byte(listOutput), 0o644); err != nil {
		t.Fatal(err)
	}

	path := filepath.Join(dir, "netbird")
	script := fmt.Sprintf(`#!/bin/sh
list=%q
if [ "$1" = profile ] && [ "$2" = list ]; then
  cat "$list"
  exit 0
fi
if [ "$1" = profile ] && [ "$2" = select ]; then
  exit 0
fi
exit 1
`, listPath)
	if err := os.WriteFile(path, []byte(script), 0o755); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestCLIBinDefault(t *testing.T) {
	c := CLI{}
	if got := c.bin(); got != "netbird" {
		t.Fatalf("bin() = %q, want netbird", got)
	}
}

func TestListProfilesCLI(t *testing.T) {
	dir := t.TempDir()
	bin := writeFakeNetbird(t, dir, sampleList)

	profiles, err := CLI{Bin: bin}.ListProfiles(context.Background())
	if err != nil {
		t.Fatalf("ListProfiles() error = %v", err)
	}
	if len(profiles) != 3 {
		t.Fatalf("len = %d, want 3", len(profiles))
	}
}

func TestListProfilesCLIError(t *testing.T) {
	dir := t.TempDir()
	bin := filepath.Join(dir, "netbird")
	if err := os.WriteFile(bin, []byte("#!/bin/sh\nexit 1\n"), 0o755); err != nil {
		t.Fatal(err)
	}

	_, err := CLI{Bin: bin}.ListProfiles(context.Background())
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "netbird profile list") {
		t.Fatalf("error = %v", err)
	}
}

func TestSelectCLI(t *testing.T) {
	dir := t.TempDir()
	bin := writeFakeNetbird(t, dir, sampleList)

	cli := CLI{Bin: bin}
	if err := cli.SelectProfile(context.Background(), "home"); err != nil {
		t.Fatalf("SelectProfile() error = %v", err)
	}
}

func TestSelectCLIError(t *testing.T) {
	dir := t.TempDir()
	bin := filepath.Join(dir, "netbird")
	script := "#!/bin/sh\nif [ \"$1\" = profile ] && [ \"$2\" = select ]; then echo oops >&2; exit 1; fi\nexit 1\n"
	if err := os.WriteFile(bin, []byte(script), 0o755); err != nil {
		t.Fatal(err)
	}

	err := CLI{Bin: bin}.SelectProfile(context.Background(), "home")
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "netbird profile select home") {
		t.Fatalf("error = %v", err)
	}
}

func TestResolveProfileNotFound(t *testing.T) {
	profiles, err := ParseProfileList([]byte(sampleList))
	if err != nil {
		t.Fatal(err)
	}
	_, err = ResolveProfile("missing", profiles)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestCurrentProfileNoneActive(t *testing.T) {
	profiles := []Profile{{Name: "a", Active: false}}
	_, err := CurrentProfile(profiles)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestFormatProfile(t *testing.T) {
	if got := FormatProfile(Profile{Name: "work"}); got != "work" {
		t.Fatalf("FormatProfile() = %q", got)
	}
}
