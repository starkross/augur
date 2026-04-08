---
id: otel-028
title: "OTEL-028: spike_limit_mib >= limit_mib"
sidebar_label: OTEL-028
description: The soft limit must sit below the hard limit.
---

# OTEL-028: `spike_limit_mib` >= `limit_mib` (soft limit zero or negative)

**Severity:** <span className="badge badge--danger">deny (blocking)</span>

## Rule Details

`memory_limiter` uses two thresholds: `limit_mib` is the hard ceiling and `spike_limit_mib` is the headroom subtracted from it to form a soft ceiling (`limit_mib - spike_limit_mib`). If `spike_limit_mib >= limit_mib` the soft ceiling is zero or negative, which means the limiter jumps straight to "refuse everything" the moment the Collector allocates anything at all. That is not how you want back-pressure to behave.

This rule fires when `spike_limit_mib >= limit_mib` on any `memory_limiter` processor.

## Options

| Field | Constraint |
|------|------------|
| `spike_limit_mib` | Must be < `limit_mib` |

## Examples

:::danger[Incorrect]

```yaml
processors:
  memory_limiter:
    check_interval: 5s
    limit_mib: 512
    spike_limit_mib: 768           # larger than limit_mib
```

:::

:::tip[Correct]

```yaml
processors:
  memory_limiter:
    check_interval: 5s
    limit_mib: 512
    spike_limit_mib: 128
```

:::

## When Not To Use It

Never — the configuration cannot do anything useful.

## Related Rules

- [OTEL-001](../core/otel-001) — `memory_limiter` processor must be configured
- [OTEL-027](./otel-027) — `memory_limiter` `check_interval` is 0 or unset
- [OTEL-029](./otel-029) — neither `limit_mib` nor `limit_percentage` set
- [OTEL-030](./otel-030) — `limit_percentage` outside safe range

## Version

Available since augur v0.1.0.

## Further Reading

- [OpenTelemetry Collector — memory_limiter processor](https://github.com/open-telemetry/opentelemetry-collector/blob/main/processor/memorylimiterprocessor/README.md)

## Resources

- Rule source: [`policy/main/memory.rego`](https://github.com/starkross/augur/blob/main/policy/main/memory.rego)
