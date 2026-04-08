---
id: otel-061
title: "OTEL-061: memory_ballast extension (deprecated, use GOMEMLIMIT)"
sidebar_label: OTEL-061
description: memory_ballast is a pre-Go-1.19 workaround that is no longer needed.
---

# OTEL-061: `memory_ballast` extension (deprecated, use `GOMEMLIMIT`)

**Severity:** <span className="badge badge--warning">warn (advisory)</span>

## Rule Details

`memory_ballast` was a Collector workaround for the Go garbage collector's old behavior — it pre-allocated a large unused block so the GC trigger would sit at a predictable fraction of the container limit. Since Go 1.19 the same effect is achieved by setting the `GOMEMLIMIT` environment variable, and the Collector team has deprecated the extension. Remove `memory_ballast` and set `GOMEMLIMIT` on the process instead.

This rule fires when an extension whose base name is `memory_ballast` is configured.

## Options

This rule has no options.

## Examples

:::warning[Avoid]

```yaml
extensions:
  memory_ballast:
    size_mib: 512
```

:::

:::tip[Prefer]

```yaml
# remove the extension — set GOMEMLIMIT on the process:
#   env:
#     GOMEMLIMIT: 1GiB
extensions:
  health_check:
    endpoint: "localhost:13133"
```

:::

## When Not To Use It

Never — the extension is deprecated. Use `GOMEMLIMIT`.

## Related Rules

- [OTEL-073](../lifecycle/otel-073) — `memory_limiter` `ballast_size_mib` (deprecated)
- [OTEL-001](../core/otel-001) — `memory_limiter` processor must be configured

## Version

Available since augur v0.1.0.

## Further Reading

- [Go 1.19 release notes — runtime/debug.SetMemoryLimit](https://go.dev/doc/go1.19#runtime)
- [Go — GOMEMLIMIT](https://pkg.go.dev/runtime#hdr-Environment_Variables)

## Resources

- Rule source: [`policy/main/extension.rego`](https://github.com/starkross/augur/blob/main/policy/main/extension.rego)
