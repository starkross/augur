---
id: otel-073
title: "OTEL-073: memory_limiter ballast_size_mib (deprecated, use GOMEMLIMIT)"
sidebar_label: OTEL-073
description: ballast_size_mib is obsolete — set GOMEMLIMIT on the process instead.
---

# OTEL-073: `memory_limiter` `ballast_size_mib` (deprecated, use `GOMEMLIMIT`)

**Severity:** <span className="badge badge--warning">warn (advisory)</span>

## Rule Details

The `ballast_size_mib` knob on `memory_limiter` was a Go GC workaround for a problem that Go 1.19+ solved with the `GOMEMLIMIT` environment variable. Using `ballast_size_mib` now runs the workaround *and* the new solution side by side, which makes GC behavior harder to reason about. Remove the field from the processor config and set `GOMEMLIMIT` on the process.

This rule fires when a `memory_limiter` processor has `ballast_size_mib` set.

## Options

This rule has no options.

## Examples

:::warning[Avoid]

```yaml
processors:
  memory_limiter:
    check_interval: 5s
    limit_percentage: 80
    ballast_size_mib: 256          # deprecated
```

:::

:::tip[Prefer]

```yaml
# remove ballast_size_mib and set GOMEMLIMIT on the process:
#   env:
#     GOMEMLIMIT: 1GiB
processors:
  memory_limiter:
    check_interval: 5s
    limit_percentage: 80
    spike_limit_percentage: 25
```

:::

## When Not To Use It

Collectors running on Go versions older than 1.19 where `GOMEMLIMIT` is not available. In that case, upgrade Go.

## Related Rules

- [OTEL-001](../core/otel-001) — `memory_limiter` processor must be configured
- [OTEL-061](../extension/otel-061) — `memory_ballast` extension deprecated

## Version

Available since augur v0.1.0.

## Further Reading

- [Go — GOMEMLIMIT](https://pkg.go.dev/runtime#hdr-Environment_Variables)
- [OpenTelemetry Collector — memory_limiter processor](https://github.com/open-telemetry/opentelemetry-collector/blob/main/processor/memorylimiterprocessor/README.md)

## Resources

- Rule source: [`policy/main/lifecycle.rego`](https://github.com/starkross/augur/blob/main/policy/main/lifecycle.rego)
