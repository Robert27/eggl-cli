package ui

import (
	"bytes"
	"strings"
	"testing"
)

func TestRenderEnvShowPlain(t *testing.T) {
	t.Setenv("NO_COLOR", "1")

	var buf bytes.Buffer
	RenderEnvShow(&buf, EnvShowReport{
		ActiveProfile: "alpha",
		KubeContext:   "ctx-a",
		Tailscale:     "b3e1 (example-alpha.internal)",
		ConfigPath:    "/tmp/config.yaml",
		Profiles: []EnvProfile{
			{Name: "alpha", KubeContext: "ctx-a", TailscaleAccount: "b3e1"},
		},
	})

	got := buf.String()
	for _, want := range []string{
		"profile: alpha",
		"kube: ctx-a",
		"tailscale: b3e1 (example-alpha.internal)",
		"config: /tmp/config.yaml",
		"alpha: kube=ctx-a tailscale=b3e1",
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("expected %q in output, got %q", want, got)
		}
	}
}

func TestRenderEnvShowPlainUnknown(t *testing.T) {
	t.Setenv("NO_COLOR", "1")

	var buf bytes.Buffer
	RenderEnvShow(&buf, EnvShowReport{Unknown: true, KubeContext: "other"})

	if !strings.Contains(buf.String(), "profile: unknown") {
		t.Fatalf("output = %q", buf.String())
	}
}

func TestRenderEnvSwitchPlain(t *testing.T) {
	t.Setenv("NO_COLOR", "1")

	var buf bytes.Buffer
	RenderEnvSwitch(&buf, EnvSwitchResult{
		FromProfile: "alpha",
		ToProfile:   "beta",
		FromKube:    "ctx-a",
		ToKube:      "ctx-b",
		FromTS:      "b3e1 (example-alpha.internal)",
		ToTS:        "a7f2 (example-beta.internal)",
	})

	got := buf.String()
	for _, want := range []string{
		"profile: alpha → beta",
		"kube: ctx-a → ctx-b",
		"tailscale: b3e1 (example-alpha.internal) → a7f2 (example-beta.internal)",
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("expected %q in output, got %q", want, got)
		}
	}
}

func TestEmptyDash(t *testing.T) {
	if got := emptyDash(""); got != "—" {
		t.Fatalf("emptyDash(\"\") = %q", got)
	}
	if got := emptyDash("alpha"); got != "alpha" {
		t.Fatalf("emptyDash(\"alpha\") = %q", got)
	}
}
