---
id: otel-071
title: "OTEL-071: logging exporter deprecated (renamed to debug)"
sidebar_label: OTEL-071
description: The logging exporter was renamed to debug in v0.111.0.
---

# OTEL-071: `logging` exporter deprecated (renamed to `debug` in v0.111.0)

**Severity:** <span className="badge badge--warning">warn (advisory)</span>

## Rule Details

The exporter previously named `logging` was renamed to `debug` in OpenTelemetry Collector v0.111.0. Both names currently work but `logging` is slated for removal. Renaming is a mechanical find-and-replace and buys you one fewer deprecation warning in the Collector logs.

This rule fires when an exporter whose base name is `logging` is declared.

## Options

This rule has no options.

## Examples

:::warning[Avoid]

```yaml
exporters:
  logging: {}
```

:::

:::tip[Prefer]

```yaml
exporters:
  debug: {}
```

:::

## When Not To Use It

Collectors pinned to a version older than v0.111.0 where the `debug` name does not yet exist. Upgrade and rename.

## Related Rules

- [OTEL-015](../core/otel-015) — `debug`/`logging` exporter detected
- [OTEL-072](./otel-072) — OpenCensus receiver/exporter deprecated

## Version

Available since augur v0.1.0.

## Further Reading

- [OpenTelemetry Collector — debug exporter](https://github.com/open-telemetry/opentelemetry-collector/blob/main/exporter/debugexporter/README.md)

## Resources

- Rule source: [`policy/main/lifecycle.rego`](https://github.com/starkross/augur/blob/main/policy/main/lifecycle.rego)
