package tailscale

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestParseAccountsJSON(t *testing.T) {
	data, err := os.ReadFile(filepath.Join("testdata", "accounts.json"))
	if err != nil {
		t.Fatal(err)
	}

	accounts, err := ParseAccountsJSON(data)
	if err != nil {
		t.Fatalf("ParseAccountsJSON() error = %v", err)
	}
	if len(accounts) != 3 {
		t.Fatalf("len = %d, want 3", len(accounts))
	}
}

func TestResolveAccount(t *testing.T) {
	accounts, err := ParseAccountsJSON(mustRead(t, "testdata/accounts.json"))
	if err != nil {
		t.Fatal(err)
	}

	cases := []struct {
		ref  string
		want string
	}{
		{"b3e1", "b3e1"},
		{"example-alpha.internal", "b3e1"},
		{"user-a@example.com", "b3e1"},
		{"a7f2", "a7f2"},
		{"example-beta.internal", "a7f2"},
	}

	for _, tc := range cases {
		a, err := ResolveAccount(tc.ref, accounts)
		if err != nil {
			t.Fatalf("ResolveAccount(%q) error = %v", tc.ref, err)
		}
		if a.ID != tc.want {
			t.Fatalf("ResolveAccount(%q) id = %q, want %q", tc.ref, a.ID, tc.want)
		}
	}
}

func TestCurrentAccount(t *testing.T) {
	accounts, err := ParseAccountsJSON(mustRead(t, "testdata/accounts.json"))
	if err != nil {
		t.Fatal(err)
	}

	a, err := CurrentAccount(accounts)
	if err != nil {
		t.Fatalf("CurrentAccount() error = %v", err)
	}
	if a.ID != "b3e1" {
		t.Fatalf("id = %q, want b3e1", a.ID)
	}
}

func mustRead(t *testing.T, path string) []byte {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return data
}

func writeFakeTailscale(t *testing.T, dir string, listJSON []byte) string {
	t.Helper()

	jsonPath := filepath.Join(dir, "accounts.json")
	if err := os.WriteFile(jsonPath, listJSON, 0o644); err != nil {
		t.Fatal(err)
	}

	path := filepath.Join(dir, "tailscale")
	script := fmt.Sprintf(`#!/bin/sh
json=%q
if [ "$1" = switch ] && [ "$2" = --list ] && [ "$3" = --json ]; then
  cat "$json"
  exit 0
fi
if [ "$1" = switch ]; then
  exit 0
fi
exit 1
`, jsonPath)
	if err := os.WriteFile(path, []byte(script), 0o755); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestCLIBinDefault(t *testing.T) {
	c := CLI{}
	if got := c.bin(); got != "tailscale" {
		t.Fatalf("bin() = %q, want tailscale", got)
	}
}

func TestCLIBinCustom(t *testing.T) {
	c := CLI{Bin: "/custom/tailscale"}
	if got := c.bin(); got != "/custom/tailscale" {
		t.Fatalf("bin() = %q", got)
	}
}

func TestListAccountsCLI(t *testing.T) {
	dir := t.TempDir()
	jsonData := mustRead(t, "testdata/accounts.json")
	bin := writeFakeTailscale(t, dir, jsonData)

	accounts, err := CLI{Bin: bin}.ListAccounts(context.Background())
	if err != nil {
		t.Fatalf("ListAccounts() error = %v", err)
	}
	if len(accounts) != 3 {
		t.Fatalf("len = %d, want 3", len(accounts))
	}
}

func TestListAccountsCLIError(t *testing.T) {
	dir := t.TempDir()
	bin := filepath.Join(dir, "tailscale")
	if err := os.WriteFile(bin, []byte("#!/bin/sh\nexit 1\n"), 0o755); err != nil {
		t.Fatal(err)
	}

	_, err := CLI{Bin: bin}.ListAccounts(context.Background())
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "tailscale switch --list --json") {
		t.Fatalf("error = %v", err)
	}
}

func TestSwitchCLI(t *testing.T) {
	dir := t.TempDir()
	bin := writeFakeTailscale(t, dir, []byte("[]"))

	cli := CLI{Bin: bin}
	if err := cli.Switch(context.Background(), "b3e1"); err != nil {
		t.Fatalf("Switch() error = %v", err)
	}
}

func TestSwitchCLIError(t *testing.T) {
	dir := t.TempDir()
	bin := filepath.Join(dir, "tailscale")
	script := "#!/bin/sh\nif [ \"$1\" = switch ]; then echo oops >&2; exit 1; fi\nexit 1\n"
	if err := os.WriteFile(bin, []byte(script), 0o755); err != nil {
		t.Fatal(err)
	}

	err := CLI{Bin: bin}.Switch(context.Background(), "b3e1")
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "tailscale switch b3e1") {
		t.Fatalf("error = %v", err)
	}
}

func TestFormatAccount(t *testing.T) {
	if got := FormatAccount(Account{ID: "abc"}); got != "abc" {
		t.Fatalf("FormatAccount() = %q, want abc", got)
	}
	if got := FormatAccount(Account{ID: "abc", Tailnet: "example.internal"}); got != "abc (example.internal)" {
		t.Fatalf("FormatAccount() = %q", got)
	}
}

func TestResolveAccountNotFound(t *testing.T) {
	accounts, err := ParseAccountsJSON(mustRead(t, "testdata/accounts.json"))
	if err != nil {
		t.Fatal(err)
	}
	_, err = ResolveAccount("missing", accounts)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestResolveAccountAmbiguous(t *testing.T) {
	accounts := []Account{
		{ID: "aaa1", Nickname: "shared"},
		{ID: "bbb2", Nickname: "shared"},
	}

	_, err := ResolveAccount("shared", accounts)
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "ambiguous tailscale account reference") {
		t.Fatalf("error = %v", err)
	}
	if !strings.Contains(err.Error(), "aaa1") || !strings.Contains(err.Error(), "bbb2") {
		t.Fatalf("error should list matching ids, got %v", err)
	}
}

func TestCurrentAccountNoneSelected(t *testing.T) {
	accounts := []Account{{ID: "a", Selected: false}}
	_, err := CurrentAccount(accounts)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestParseAccountsJSONInvalid(t *testing.T) {
	_, err := ParseAccountsJSON([]byte("not json"))
	if err == nil {
		t.Fatal("expected error")
	}
}
