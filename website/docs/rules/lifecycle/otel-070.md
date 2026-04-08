---
id: otel-070
title: "OTEL-070: Telemetry metrics address bound to 0.0.0.0"
sidebar_label: OTEL-070
description: Exposing internal Collector metrics on every interface leaks operational details.
---

# OTEL-070: Telemetry metrics address bound to `0.0.0.0`

**Severity:** <span className="badge badge--warning">warn (advisory)</span>

## Rule Details

The Collector's internal metrics endpoint exposes operational details (pipeline names, exporter queue lengths, retry counts, internal version). Binding it to `0.0.0.0` makes that reachable from every network the host is on. Bind it to `localhost` and let a side-loaded scrape job (a Prometheus pod on the same node, a sidecar) read it privately.

This rule fires when `service.telemetry.metrics.address` contains `0.0.0.0`.

## Options

This rule has no options.

## Examples

:::warning[Avoid]

```yaml
service:
  telemetry:
    metrics:
      address: "0.0.0.0:8888"
```

:::

:::tip[Prefer]

```yaml
service:
  telemetry:
    metrics:
      address: "localhost:8888"
```

:::

## When Not To Use It

Never — there is a safer alternative.

## Related Rules

- [OTEL-010](../core/otel-010) — receivers should not bind to `0.0.0.0`
- [OTEL-060](../extension/otel-060) — `zpages` endpoint bound to `0.0.0.0`
- [OTEL-069](./otel-069) — telemetry metrics level set to `none`
- [OTEL-074](./otel-074) — `service.telemetry.metrics.address` deprecated

## Version

Available since augur v0.1.0.

## Further Reading

- [OpenTelemetry Collector — internal telemetry](https://opentelemetry.io/docs/collector/internal-telemetry/)

## Resources

- Rule source: [`policy/main/lifecycle.rego`](https://github.com/starkross/augur/blob/main/policy/main/lifecycle.rego)
