---
id: otel-017
title: "OTEL-017: Exporter missing retry_on_failure/sending_queue"
sidebar_label: OTEL-017
description: Network exporters need retry and/or queueing to tolerate brief outages.
---

# OTEL-017: Exporter missing `retry_on_failure`/`sending_queue`

**Severity:** <span className="badge badge--warning">warn (advisory)</span>

## Rule Details

Any exporter that pushes data over the network will eventually see a transient failure — TLS handshake timeout, backend restart, rate limit, DNS blip. Without `retry_on_failure` or `sending_queue` the failed batch is dropped on the floor. Pull-based exporters (`prometheus`, `prometheusremotewrite`) and diagnostic ones (`debug`, `logging`) do not need this; everything else does.

This rule fires when an exporter whose base type is not pull-based has neither `retry_on_failure` nor `sending_queue` configured.

### AWS exporter alternative retry (`max_retries`)

Certain AWS exporters implement their own retry logic via `max_retries` instead of the standard `retry_on_failure`/`sending_queue` fields. The rule recognises this for the following exporter types:

| Exporter | Alternative field | Notes |
|----------|------------------|-------|
| `awsemf` | `max_retries` | AWS CloudWatch EMF exporter |
| `awscloudwatchlogs` | `max_retries` | AWS CloudWatch Logs exporter |
| `awsxray` | `max_retries` | AWS X-Ray exporter |
| `awss3` | `max_retries` | AWS S3 exporter |

When `max_retries` is set on one of these exporters, the rule does **not** fire.

> **Note on `sending_queue`:** The `max_retries` exemption intentionally covers the absence of `sending_queue` as well. `max_retries` only provides retry — it does not offer durable queueing like `sending_queue`. This is by design: AWS exporters manage back-pressure and transient failures through their own SDK-level retry logic, and adding a Collector-level `sending_queue` on top is unnecessary for most use cases. If you need durable queueing (e.g. surviving Collector restarts), configure `sending_queue` with persistent storage explicitly.

The `max_retries` exemption applies **only** to the AWS exporters listed above. Other exporters with a `max_retries` field will still trigger this rule — use `retry_on_failure` and/or `sending_queue` for those.

## Options

This rule has no options. The set of pull-based exporter types (`debug`, `logging`, `prometheus`, `prometheusremotewrite`) and the AWS exporter allowlist (`awsemf`, `awscloudwatchlogs`, `awsxray`, `awss3`) are defined inside the policy.

## Examples

:::warning[Avoid]

```yaml
exporters:
  otlp/backend:
    endpoint: backend:4317
    # no retry_on_failure or sending_queue
```

:::

:::tip[Prefer]

```yaml
exporters:
  otlp/backend:
    endpoint: backend:4317
    retry_on_failure:
      enabled: true
      max_elapsed_time: 300s
    sending_queue:
      enabled: true
      num_consumers: 10
      queue_size: 5000
```

:::

### AWS exporter with `max_retries`

:::tip[Prefer]

```yaml
exporters:
  awsemf/prod:
    log_retention: 30
    max_retries: 3
```

:::

:::warning[Avoid]

```yaml
exporters:
  awsemf/prod:
    log_retention: 30
    # no max_retries, retry_on_failure, or sending_queue
```

:::

## When Not To Use It

Never for production. Leave enabled so any forgotten retry/queue configuration is flagged.

## Related Rules

- [OTEL-048](../exporter/otel-048) — `sending_queue` explicitly disabled
- [OTEL-049](../exporter/otel-049) — `sending_queue.queue_size` below 10
- [OTEL-050](../exporter/otel-050) — `sending_queue.queue_size` above 50000
- [OTEL-053](../exporter/otel-053) — retry `max_elapsed_time` set to 0
- [OTEL-065](../reliability/otel-065) — `sending_queue` without persistent storage

## Version

Available since augur v0.1.0.

## Further Reading

- [OpenTelemetry Collector — exporter helper configuration](https://github.com/open-telemetry/opentelemetry-collector/blob/main/exporter/exporterhelper/README.md)

## Resources

- Rule source: [`policy/main/main.rego`](https://github.com/starkross/augur/blob/main/policy/main/main.rego)
