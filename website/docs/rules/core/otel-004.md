---
id: otel-004
title: "OTEL-004: No hardcoded secrets in exporters"
sidebar_label: OTEL-004
description: Exporter credentials must come from environment variables, not literal strings.
---

# OTEL-004: No hardcoded secrets in exporters

**Severity:** <span className="badge badge--danger">deny (blocking)</span>

## Rule Details

Hardcoded API keys, bearer tokens, passwords, or signed URLs inside exporter config end up in version control, in CI logs, and in shared configmaps. augur scans every string value under `exporters` whose key name looks secret-like (`api_key`, `token`, `password`, `authorization`, etc.) and blocks any value that is not a `${env:VAR_NAME}` reference.

This rule fires when an exporter has a secret-like field whose value is a plain string rather than an environment variable reference.

## Options

This rule has no options. The set of "secret-like" keys is maintained inside [`policy/lib`](https://github.com/starkross/augur/tree/main/policy/lib).

## Examples

:::danger[Incorrect]

```yaml
exporters:
  otlp/vendor:
    endpoint: api.vendor.com:4317
    headers:
      api_key: "sk-hardcoded-secret"       # literal secret
```

:::

:::tip[Correct]

```yaml
exporters:
  otlp/vendor:
    endpoint: api.vendor.com:4317
    headers:
      api_key: "${env:VENDOR_API_KEY}"
```

:::

## When Not To Use It

Never. Even for local development, use a `.env` file and `${env:...}` references so the habit scales and so the config file can be safely committed.

## Related Rules

- [OTEL-005](./otel-005) — no hardcoded secrets in receivers
- [OTEL-035](../security/otel-035) — no hardcoded secrets in extensions

## Version

Available since augur v0.1.0.

## Further Reading

- [OpenTelemetry Collector — configuration environment variables](https://opentelemetry.io/docs/collector/configuration/#environment-variables)

## Resources

- Rule source: [`policy/main/main.rego`](https://github.com/starkross/augur/blob/main/policy/main/main.rego)
