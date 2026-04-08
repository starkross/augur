---
id: otel-054
title: "OTEL-054: Prometheus scrape_interval below 10s"
sidebar_label: OTEL-054
description: Sub-10s scrapes hammer the target without meaningful resolution gain.
---

# OTEL-054: Prometheus `scrape_interval` below 10s

**Severity:** <span className="badge badge--warning">warn (advisory)</span>

## Rule Details

Most Prometheus metrics have sample jitter that makes sub-10s scraping statistically noisy. What you do get is a lot more load on the target (Go runtime metrics in particular are expensive to serialize) and a lot more storage in the backend. 10s is a safe baseline; go below it only with a clear need and data to back it up.

This rule fires when any `prometheus` scrape config has `scrape_interval < 10s`.

## Options

| Field | Constraint |
|------|------------|
| `scrape_interval` | Should be ≥ 10s |

## Examples

:::warning[Avoid]

```yaml
receivers:
  prometheus:
    config:
      scrape_configs:
        - job_name: app
          scrape_interval: 1s       # too aggressive
          static_configs:
            - targets: [app:9090]
```

:::

:::tip[Prefer]

```yaml
receivers:
  prometheus:
    config:
      scrape_configs:
        - job_name: app
          scrape_interval: 15s
          static_configs:
            - targets: [app:9090]
```

:::

## When Not To Use It

Short-lived diagnostic runs where you want high-frequency samples for a few minutes. Use a separate temporary job in that case and leave the rule enabled globally.

## Related Rules

- [OTEL-055](./otel-055) — `hostmetrics` `collection_interval` below 10s

## Version

Available since augur v0.1.0.

## Further Reading

- [OpenTelemetry Collector Contrib — prometheus receiver](https://github.com/open-telemetry/opentelemetry-collector-contrib/blob/main/receiver/prometheusreceiver/README.md)

## Resources

- Rule source: [`policy/main/receiver.rego`](https://github.com/starkross/augur/blob/main/policy/main/receiver.rego)
