---
id: otel-030
title: "OTEL-030: memory_limiter limit_percentage outside safe range"
sidebar_label: OTEL-030
description: A percentage under 20% wastes capacity; over 90% risks OOM before the limiter fires.
---

# OTEL-030: `memory_limiter` `limit_percentage` outside safe range (20–90%)

**Severity:** <span className="badge badge--warning">warn (advisory)</span>

## Rule Details

Below 20%, the limiter refuses work while the vast majority of the allocated memory sits idle — you are paying for capacity you never use. Above 90%, there is so little headroom between the soft ceiling and the container's cgroup limit that a brief allocation spike (Go's garbage collector lagging, a large incoming batch) can push the process into OOM-kill territory before the next `check_interval` tick fires.

This rule fires when `limit_percentage < 20` or `limit_percentage > 90`.

## Options

| Field | Constraint |
|------|------------|
| `limit_percentage` | Should be between 20 and 90 inclusive |

## Examples

:::warning[Avoid]

```yaml
processors:
  memory_limiter:
    check_interval: 5s
    limit_percentage: 95             # too close to the hard limit
```

:::

:::tip[Prefer]

```yaml
processors:
  memory_limiter:
    check_interval: 5s
    limit_percentage: 80
    spike_limit_percentage: 25
```

:::

## When Not To Use It

A Collector with tightly measured allocations and a very short `check_interval` may run safely above 90%. Verify with load tests before disabling.

## Related Rules

- [OTEL-001](../core/otel-001) — `memory_limiter` processor must be configured
- [OTEL-027](./otel-027) — `memory_limiter` `check_interval` is 0 or unset
- [OTEL-028](./otel-028) — `spike_limit_mib` >= `limit_mib`
- [OTEL-029](./otel-029) — neither `limit_mib` nor `limit_percentage` set

## Version

Available since augur v0.1.0.

## Further Reading

- [OpenTelemetry Collector — memory_limiter processor](https://github.com/open-telemetry/opentelemetry-collector/blob/main/processor/memorylimiterprocessor/README.md)

## Resources

- Rule source: [`policy/main/memory.rego`](https://github.com/starkross/augur/blob/main/policy/main/memory.rego)
