---
id: otel-059
title: "OTEL-059: pprof extension enabled in production"
sidebar_label: OTEL-059
description: pprof exposes runtime profiling — useful in dev, a disclosure risk in prod.
---

# OTEL-059: `pprof` extension enabled in production

**Severity:** <span className="badge badge--warning">warn (advisory)</span>

## Rule Details

The `pprof` extension exposes Go runtime profiling endpoints (`/debug/pprof/heap`, `/debug/pprof/goroutine`, etc.). In development it is invaluable; in production it can leak environment-specific details (goroutine stacks mention internal package paths, heap dumps can contain secrets from recently allocated objects) and is easy to miss in a port scan. If you need it in production, gate it behind auth and a private interface.

This rule fires when an extension whose base name is `pprof` is configured.

## Options

This rule has no options.

## Examples

:::warning[Avoid]

```yaml
extensions:
  pprof:
    endpoint: "0.0.0.0:1777"

service:
  extensions: [pprof]
```

:::

:::tip[Prefer]

```yaml
extensions:
  # pprof disabled in production; enable via a separate config for debugging
  health_check:
    endpoint: "localhost:13133"

service:
  extensions: [health_check]
```

:::

## When Not To Use It

Brief debugging windows on a pre-production Collector, or production Collectors with a strict port firewall and authenticated access. Remove or disable afterwards.

## Related Rules

- [OTEL-060](./otel-060) — `zpages` endpoint bound to `0.0.0.0`
- [OTEL-062](./otel-062) — extension in `service.extensions` but not defined

## Version

Available since augur v0.1.0.

## Further Reading

- [OpenTelemetry Collector Contrib — pprof extension](https://github.com/open-telemetry/opentelemetry-collector-contrib/blob/main/extension/pprofextension/README.md)

## Resources

- Rule source: [`policy/main/extension.rego`](https://github.com/starkross/augur/blob/main/policy/main/extension.rego)
