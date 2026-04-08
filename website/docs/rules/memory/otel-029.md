---
id: otel-029
title: "OTEL-029: Neither limit_mib nor limit_percentage set on memory_limiter"
sidebar_label: OTEL-029
description: memory_limiter has to know what it is limiting.
---

# OTEL-029: Neither `limit_mib` nor `limit_percentage` set on `memory_limiter`

**Severity:** <span className="badge badge--danger">deny (blocking)</span>

## Rule Details

`memory_limiter` is configured with either an absolute limit (`limit_mib`) or a percentage of the host/container memory (`limit_percentage`). If neither is set, the limiter has nothing to compare against — the processor is inert. Pick one. Absolute limits are easier to reason about when cgroup limits are fixed; percentages are easier when pods scale vertically.

This rule fires when a `memory_limiter` processor has neither `limit_mib` nor `limit_percentage` set.

## Options

This rule has no options.

## Examples

:::danger[Incorrect]

```yaml
processors:
  memory_limiter:
    check_interval: 5s
    # no limit_mib or limit_percentage
```

:::

:::tip[Correct]

```yaml
processors:
  memory_limiter:
    check_interval: 5s
    limit_percentage: 80
    spike_limit_percentage: 25
```

:::

## When Not To Use It

Never — a limiter with no limit configured does not protect anything.

## Related Rules

- [OTEL-001](../core/otel-001) — `memory_limiter` processor must be configured
- [OTEL-027](./otel-027) — `memory_limiter` `check_interval` is 0 or unset
- [OTEL-028](./otel-028) — `spike_limit_mib` >= `limit_mib`
- [OTEL-030](./otel-030) — `limit_percentage` outside safe range

## Version

Available since augur v0.1.0.

## Further Reading

- [OpenTelemetry Collector — memory_limiter processor](https://github.com/open-telemetry/opentelemetry-collector/blob/main/processor/memorylimiterprocessor/README.md)

## Resources

- Rule source: [`policy/main/memory.rego`](https://github.com/starkross/augur/blob/main/policy/main/memory.rego)
