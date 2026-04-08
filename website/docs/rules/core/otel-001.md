---
id: otel-001
title: "OTEL-001: memory_limiter processor must be configured"
sidebar_label: OTEL-001
description: memory_limiter is required to prevent the Collector from being OOM-killed.
---

# OTEL-001: memory_limiter processor must be configured

**Severity:** <span className="badge badge--danger">deny (blocking)</span>

## Rule Details

The `memory_limiter` processor is the Collector's safety valve against out-of-memory kills. Without it, a slow exporter, a traffic burst, or a stuck sending queue can push the process past its cgroup limit and get SIGKILL'd — dropping every batch in flight and taking observability down at the worst possible moment.

This rule fires when no `memory_limiter` processor is declared in the top-level `processors` section.

## Options

This rule has no options.

## Examples

:::danger[Incorrect]

```yaml
processors:
  batch:
    timeout: 1s
```

:::

:::tip[Correct]

```yaml
processors:
  memory_limiter:
    check_interval: 5s
    limit_percentage: 80
    spike_limit_percentage: 25
  batch:
    timeout: 1s
```

:::

## When Not To Use It

Short-lived test runs, one-shot config imports, or throwaway sidecars where OOM restart is acceptable. For any long-running production Collector you should leave this rule enabled.

## Related Rules

- [OTEL-002](./otel-002) — `memory_limiter` must be included in every pipeline
- [OTEL-014](./otel-014) — `memory_limiter` should be first processor
- [OTEL-027](../memory/otel-027) — `memory_limiter` `check_interval` is 0 or unset
- [OTEL-028](../memory/otel-028) — `spike_limit_mib` >= `limit_mib`
- [OTEL-029](../memory/otel-029) — neither `limit_mib` nor `limit_percentage` set

## Version

Available since augur v0.1.0.

## Further Reading

- [OpenTelemetry Collector — memory_limiter processor](https://github.com/open-telemetry/opentelemetry-collector/blob/main/processor/memorylimiterprocessor/README.md)
- [OpenTelemetry Collector — sizing](https://opentelemetry.io/docs/collector/scaling/)

## Resources

- Rule source: [`policy/main/main.rego`](https://github.com/starkross/augur/blob/main/policy/main/main.rego)
