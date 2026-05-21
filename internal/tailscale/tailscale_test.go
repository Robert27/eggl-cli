package tailscale

import (
	"os"
	"path/filepath"
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
