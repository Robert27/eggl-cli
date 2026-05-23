package tailscale

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"sort"
	"strings"
)

type Account struct {
	ID       string
	Tailnet  string
	Account  string
	Nickname string
	Selected bool
}

type Runner interface {
	ListAccounts(ctx context.Context) ([]Account, error)
	Switch(ctx context.Context, accountID string) error
}

type CLI struct {
	Bin string
}

func (c CLI) bin() string {
	if c.Bin != "" {
		return c.Bin
	}
	return "tailscale"
}

func (c CLI) ListAccounts(ctx context.Context) ([]Account, error) {
	cmd := exec.CommandContext(ctx, c.bin(), "switch", "--list", "--json")
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("tailscale switch --list --json: %w", err)
	}
	return ParseAccountsJSON(out)
}

func (c CLI) Switch(ctx context.Context, accountID string) error {
	cmd := exec.CommandContext(ctx, c.bin(), "switch", accountID)
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("tailscale switch %s: %w: %s", accountID, err, strings.TrimSpace(string(out)))
	}
	return nil
}

type accountJSON struct {
	ID       string `json:"id"`
	Tailnet  string `json:"tailnet"`
	Account  string `json:"account"`
	Nickname string `json:"nickname"`
	Selected bool   `json:"selected"`
}

func ParseAccountsJSON(data []byte) ([]Account, error) {
	var raw []accountJSON
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("parse tailscale accounts: %w", err)
	}

	accounts := make([]Account, len(raw))
	for i, a := range raw {
		accounts[i] = Account{
			ID:       a.ID,
			Tailnet:  a.Tailnet,
			Account:  a.Account,
			Nickname: a.Nickname,
			Selected: a.Selected,
		}
	}
	return accounts, nil
}

func ResolveAccount(ref string, accounts []Account) (Account, error) {
	ref = strings.TrimSpace(ref)
	if ref == "" {
		return Account{}, fmt.Errorf("empty tailscale account reference")
	}

	var matches []Account
	for _, a := range accounts {
		if strings.EqualFold(a.ID, ref) ||
			strings.EqualFold(a.Tailnet, ref) ||
			strings.EqualFold(a.Account, ref) ||
			strings.EqualFold(a.Nickname, ref) {
			matches = append(matches, a)
		}
	}

	switch len(matches) {
	case 0:
		return Account{}, fmt.Errorf("tailscale account %q not found in switch --list", ref)
	case 1:
		return matches[0], nil
	default:
		ids := make([]string, len(matches))
		for i, a := range matches {
			ids[i] = a.ID
		}
		sort.Strings(ids)
		return Account{}, fmt.Errorf("ambiguous tailscale account reference %q: matches %s", ref, strings.Join(ids, ", "))
	}
}

func CurrentAccount(accounts []Account) (Account, error) {
	for _, a := range accounts {
		if a.Selected {
			return a, nil
		}
	}
	return Account{}, fmt.Errorf("no selected tailscale account")
}

func FormatAccount(a Account) string {
	if a.Tailnet != "" {
		return fmt.Sprintf("%s (%s)", a.ID, a.Tailnet)
	}
	return a.ID
}
