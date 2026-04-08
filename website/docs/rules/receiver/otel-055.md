---
id: otel-055
title: "OTEL-055: hostmetrics collection_interval below 10s"
sidebar_label: OTEL-055
description: Sub-10s host metric collection burns CPU on the collector host.
---

# OTEL-055: `hostmetrics` `collection_interval` below 10s

**Severity:** <span className="badge badge--warning">warn (advisory)</span>

## Rule Details

`hostmetrics` reads `/proc`, `/sys`, and platform APIs on every tick to produce CPU, memory, disk, and network metrics. At intervals below 10s the cost of the collection itself starts showing up in the host metrics you are collecting — a self-measurement feedback loop that burns CPU without giving you meaningfully better resolution.

This rule fires when a `hostmetrics` receiver has `collection_interval < 10s`.

## Options

| Field | Constraint |
|------|------------|
| `collection_interval` | Should be ≥ 10s |

## Examples

:::warning[Avoid]

```yaml
receivers:
  hostmetrics:
    collection_interval: 1s           # too aggressive
    scrapers:
      cpu:
      memory:
```

:::

:::tip[Prefer]

```yaml
receivers:
  hostmetrics:
    collection_interval: 30s
    scrapers:
      cpu:
      memory:
      disk:
      network:
```

:::

## When Not To Use It

A short forensic run during incident response where you deliberately want sub-10s granularity. Revert after the investigation.

## Related Rules

- [OTEL-054](./otel-054) — Prometheus `scrape_interval` below 10s

## Version

Available since augur v0.1.0.

## Further Reading

- [OpenTelemetry Collector Contrib — hostmetrics receiver](https://github.com/open-telemetry/opentelemetry-collector-contrib/blob/main/receiver/hostmetricsreceiver/README.md)

## Resources

- Rule source: [`policy/main/receiver.rego`](https://github.com/starkross/augur/blob/main/policy/main/receiver.rego)
