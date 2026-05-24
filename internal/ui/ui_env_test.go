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
		ShowTailscale: true,
		ConfigPath:    "/tmp/config.yaml",
		Profiles: []EnvProfile{
			{Name: "alpha", KubeContext: "ctx-a", Mesh: "tailscale:b3e1"},
		},
	})

	got := buf.String()
	for _, want := range []string{
		"profile: alpha",
		"kube: ctx-a",
		"tailscale: b3e1 (example-alpha.internal)",
		"config: /tmp/config.yaml",
		"alpha: kube=ctx-a mesh=tailscale:b3e1",
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("expected %q in output, got %q", want, got)
		}
	}
}

func TestRenderEnvShowPlainNetbird(t *testing.T) {
	t.Setenv("NO_COLOR", "1")

	var buf bytes.Buffer
	RenderEnvShow(&buf, EnvShowReport{
		ActiveProfile: "homelab",
		KubeContext:   "ctx-home",
		Netbird:       "homelab",
		ShowNetbird:   true,
		ConfigPath:    "/tmp/config.yaml",
		Profiles: []EnvProfile{
			{Name: "homelab", KubeContext: "ctx-home", Mesh: "netbird:homelab"},
		},
	})

	got := buf.String()
	for _, want := range []string{
		"netbird: homelab",
		"mesh=netbird:homelab",
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
		FromMesh:    "b3e1 (example-alpha.internal)",
		ToMesh:      "a7f2 (example-beta.internal)",
	})

	got := buf.String()
	for _, want := range []string{
		"profile: alpha → beta",
		"kube: ctx-a → ctx-b",
		"mesh: b3e1 (example-alpha.internal) → a7f2 (example-beta.internal)",
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("expected %q in output, got %q", want, got)
		}
	}
}

func TestEmptyDash(t *testing.T) {
	if got := emptyDash(""); got != "-" {
		t.Fatalf("emptyDash(\"\") = %q", got)
	}
	if got := emptyDash("alpha"); got != "alpha" {
		t.Fatalf("emptyDash(\"alpha\") = %q", got)
	}
}
