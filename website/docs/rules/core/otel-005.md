---
id: otel-005
title: "OTEL-005: No hardcoded secrets in receivers"
sidebar_label: OTEL-005
description: Receiver credentials must come from environment variables, not literal strings.
---

# OTEL-005: No hardcoded secrets in receivers

**Severity:** <span className="badge badge--danger">deny (blocking)</span>

## Rule Details

Receivers that scrape authenticated endpoints (Prometheus with basic auth, Kafka with SASL, cloud provider APIs) often carry credentials directly in the config. Any literal string under a secret-like key becomes a leak once the file is checked in or copied into a pod spec. augur blocks any secret-shaped value that is not an environment-variable reference.

This rule fires when a receiver has a secret-like field (e.g. `password`, `token`, `api_key`) whose value is a plain string rather than `${env:VAR_NAME}`.

## Options

This rule has no options.

## Examples

:::danger[Incorrect]

```yaml
receivers:
  prometheus:
    config:
      scrape_configs:
        - job_name: app
          basic_auth:
            username: metrics
            password: "hunter2"            # literal secret
```

:::

:::tip[Correct]

```yaml
receivers:
  prometheus:
    config:
      scrape_configs:
        - job_name: app
          basic_auth:
            username: metrics
            password: "${env:SCRAPE_PASSWORD}"
```

:::

## When Not To Use It

Never. Treat receiver credentials the same as exporter credentials — always source them from the environment.

## Related Rules

- [OTEL-004](./otel-004) — no hardcoded secrets in exporters
- [OTEL-035](../security/otel-035) — no hardcoded secrets in extensions

## Version

Available since augur v0.1.0.

## Further Reading

- [OpenTelemetry Collector — configuration environment variables](https://opentelemetry.io/docs/collector/configuration/#environment-variables)

## Resources

- Rule source: [`policy/main/main.rego`](https://github.com/starkross/augur/blob/main/policy/main/main.rego)
