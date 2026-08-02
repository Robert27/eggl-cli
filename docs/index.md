---
pageType: home
title: eggl CLI
titleSuffix: Small tools for the edges of your workflow
hero:
  badge: Go / terminal / workflow
  name: eggl CLI
  text: Small tools for the edges of your workflow
  tagline: A practical helper CLI for switching environments, opening Kubernetes tunnels, cleaning repository text, and taking care of the local tasks around development.
  actions:
    - theme: brand
      text: Install
      link: /install
    - theme: alt
      text: Quickstart
      link: /quickstart
features:
  - title: Switch context and account together
    details: Pair a kubectl context with a Tailscale account and move between named environments without repeating the sequence by hand.
    icon: "01"
    span: 6
  - title: Turn services into local tools
    details: Keep named Kubernetes port-forwards in one YAML file, then start them by name with optional browser opening.
    icon: "02"
    span: 6
  - title: Keep repositories predictable
    details: Preview and normalize em-dashes or line endings with file filters, Git-aware scopes, and explicit write confirmation.
    icon: "03"
    span: 6
  - title: Clean up the local edge
    details: Diagnose dependencies, free stale ports, create CI-triggering commits, and generate completions for your shell.
    icon: "04"
    span: 6
---

## What eggl CLI does

eggl CLI is intentionally small. It connects the local tools that sit between a terminal, a Kubernetes cluster, and a working repository without trying to replace any of them.

```mermaid
flowchart LR
  terminal[Terminal] --> eggl[eggl CLI]
  eggl --> environment[kubectl + Tailscale]
  eggl --> cluster[Kubernetes services]
  eggl --> repository[Git repository]
  eggl --> local[Local processes]
```

## Start here

- [Install eggl CLI](/install)
- [Configure your first profile](/quickstart)
- [Read the command reference](/commands)
