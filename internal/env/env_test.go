package env

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Robert27/eggl-cli/internal/tailscale"
)

type fakeKube struct {
	context string
	useErr  error
}

func (f *fakeKube) CurrentContext(context.Context) (string, error) {
	return f.context, nil
}

func (f *fakeKube) UseContext(_ context.Context, name string) error {
	if f.useErr != nil {
		return f.useErr
	}
	f.context = name
	return nil
}

func (f *fakeKube) PortForward(context.Context, []string) error {
	return nil
}

type fakeTS struct {
	accounts  []tailscale.Account
	switchErr error
}

func (f *fakeTS) ListAccounts(context.Context) ([]tailscale.Account, error) {
	return f.accounts, nil
}

func (f *fakeTS) Switch(_ context.Context, accountID string) error {
	if f.switchErr != nil {
		return f.switchErr
	}
	for i := range f.accounts {
		f.accounts[i].Selected = f.accounts[i].ID == accountID
	}
	return nil
}

func testAccounts() []tailscale.Account {
	return []tailscale.Account{
		{ID: "a7f2", Tailnet: "example-beta.internal", Selected: false},
		{ID: "b3e1", Tailnet: "example-alpha.internal", Selected: true},
	}
}

func writeTestConfig(t *testing.T) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "config.yaml")
	content := `profiles:
  alpha:
    kube_context: ctx-a
    tailscale_account: b3e1
  beta:
    kube_context: ctx-b
    tailscale_account: a7f2
`
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestDefaultOptions(t *testing.T) {
	opts := DefaultOptions("/tmp/config.yaml")
	if opts.ConfigPath != "/tmp/config.yaml" {
		t.Fatalf("ConfigPath = %q", opts.ConfigPath)
	}
	if opts.Kube == nil {
		t.Fatal("expected Kube runner")
	}
	if opts.TS == nil {
		t.Fatal("expected TS runner")
	}
}

func TestShowDetectsProfile(t *testing.T) {
	path := writeTestConfig(t)
	report, err := Show(context.Background(), Options{
		ConfigPath: path,
		Kube:       &fakeKube{context: "ctx-a"},
		TS:         &fakeTS{accounts: testAccounts()},
	})
	if err != nil {
		t.Fatalf("Show() error = %v", err)
	}
	if report.ActiveProfile != "alpha" {
		t.Fatalf("ActiveProfile = %q, want alpha", report.ActiveProfile)
	}
	if report.Unknown {
		t.Fatal("expected known profile")
	}
}

func TestToggleSwitchesOtherProfile(t *testing.T) {
	path := writeTestConfig(t)
	fk := &fakeKube{context: "ctx-a"}
	ft := &fakeTS{accounts: testAccounts()}

	result, err := Toggle(context.Background(), Options{
		ConfigPath: path,
		Kube:       fk,
		TS:         ft,
	})
	if err != nil {
		t.Fatalf("Toggle() error = %v", err)
	}
	if result.ToProfile != "beta" {
		t.Fatalf("ToProfile = %q, want beta", result.ToProfile)
	}
	if fk.context != "ctx-b" {
		t.Fatalf("kube context = %q, want ctx-b", fk.context)
	}
	if !accountSelected(ft.accounts, "a7f2") {
		t.Fatal("expected beta tailscale account selected")
	}
}

func TestToggleUnknownState(t *testing.T) {
	path := writeTestConfig(t)
	_, err := Toggle(context.Background(), Options{
		ConfigPath: path,
		Kube:       &fakeKube{context: "other"},
		TS:         &fakeTS{accounts: testAccounts()},
	})
	if err == nil {
		t.Fatal("expected error for unknown state")
	}
}

func TestUsePartialFailure(t *testing.T) {
	path := writeTestConfig(t)
	_, err := Use(context.Background(), Options{
		ConfigPath: path,
		Kube:       &fakeKube{context: "ctx-a", useErr: errors.New("kube failed")},
		TS:         &fakeTS{accounts: testAccounts()},
	}, "beta")
	if err == nil {
		t.Fatal("expected error")
	}
	if got := err.Error(); got == "" || !strings.Contains(got, "tailscale already switched") {
		t.Fatalf("error = %q, want tailscale partial failure message", got)
	}
}

func accountSelected(accounts []tailscale.Account, id string) bool {
	for _, a := range accounts {
		if a.ID == id && a.Selected {
			return true
		}
	}
	return false
}

func TestToggleRequiresTwoProfiles(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.yaml")
	if err := os.WriteFile(path, []byte(`profiles:
  only:
    kube_context: a
    tailscale_account: b3e1
`), 0o644); err != nil {
		t.Fatal(err)
	}

	_, err := Toggle(context.Background(), Options{
		ConfigPath: path,
		Kube:       &fakeKube{context: "a"},
		TS:         &fakeTS{accounts: testAccounts()},
	})
	if err == nil {
		t.Fatal("expected error")
	}
}
