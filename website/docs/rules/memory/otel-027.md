---
id: otel-027
title: "OTEL-027: memory_limiter check_interval is 0 or unset"
sidebar_label: OTEL-027
description: Without a positive check_interval the memory_limiter never runs.
---

# OTEL-027: `memory_limiter` `check_interval` is 0 or unset

**Severity:** <span className="badge badge--danger">deny (blocking)</span>

## Rule Details

`memory_limiter` enforces limits on a polling interval — it reads `runtime.MemStats` every `check_interval` and, if the soft or hard limit is crossed, refuses new work. An interval of zero (or no interval at all) means the check never runs: the processor is in the config and in the pipeline but never fires. That is indistinguishable from having no limiter at all.

This rule fires when a `memory_limiter` processor has `check_interval` unset or equal to a zero duration.

## Options

This rule has no options.

## Examples

:::danger[Incorrect]

```yaml
processors:
  memory_limiter:
    limit_percentage: 80
    spike_limit_percentage: 25
    # check_interval missing
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

Never — an un-polled limiter is a no-op.

## Related Rules

- [OTEL-001](../core/otel-001) — `memory_limiter` processor must be configured
- [OTEL-002](../core/otel-002) — `memory_limiter` must be included in every pipeline
- [OTEL-028](./otel-028) — `spike_limit_mib` >= `limit_mib`
- [OTEL-029](./otel-029) — neither `limit_mib` nor `limit_percentage` set
- [OTEL-030](./otel-030) — `limit_percentage` outside safe range

## Version

Available since augur v0.1.0.

## Further Reading

- [OpenTelemetry Collector — memory_limiter processor](https://github.com/open-telemetry/opentelemetry-collector/blob/main/processor/memorylimiterprocessor/README.md)

## Resources

- Rule source: [`policy/main/memory.rego`](https://github.com/starkross/augur/blob/main/policy/main/memory.rego)
