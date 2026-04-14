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

### Exporter-specific alternative retry mechanisms

Some exporters provide their own retry mechanism instead of the standard `retry_on_failure`/`sending_queue` fields. The rule recognises these alternatives:

| Exporter | Alternative field | Notes |
|----------|------------------|-------|
| `awsemf` | `max_retries` | AWS CloudWatch EMF exporter uses its own retry logic |

When one of these alternative fields is configured, the rule does **not** fire.

## Options

This rule has no options. The set of pull-based exporter types (`debug`, `logging`, `prometheus`, `prometheusremotewrite`) is exempted inside the policy.

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

:::tip[Prefer — AWS CloudWatch EMF]

The `awsemf` exporter uses `max_retries` as its own retry mechanism. Configuring it satisfies the rule:

```yaml
exporters:
  awsemf/production:
    region: us-east-1
    log_group_name: /ecs/otel
    max_retries: 5
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
