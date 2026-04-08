---
id: intro
title: Introduction
slug: /intro
sidebar_position: 1
---

# augur

**augur** is a fast, opinionated linter for [OpenTelemetry Collector](https://opentelemetry.io/docs/collector/) configurations. It catches misconfigurations, security issues, and performance pitfalls before they hit production.

Built with [OPA/Rego](https://www.openpolicyagent.org/) — every rule is a plain `.rego` file you can read, override, or extend.

## Why augur?

The OpenTelemetry Collector is flexible, but that flexibility makes it easy to ship configs that silently drop data, leak secrets, or OOM under load. `augur` encodes hard-won operational knowledge into automated checks:

- **No memory limiter?** You'll OOM in production.
- **Hardcoded API key?** It'll end up in version control.
- **Batch processor in the wrong position?** You're leaving performance on the table.

## Where to go next

- [Install](./install) — get augur on your machine
- [Quick start](./quick-start) — lint your first config
- [Usage](./usage) — flags and examples
- [Rules](./rules) — the full list of built-in checks
- [Custom policies](./custom-policies) — write your own Rego rules
- [Security](./security) — verify releases and report vulnerabilities
