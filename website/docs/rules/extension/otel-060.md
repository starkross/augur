---
id: otel-060
title: "OTEL-060: zpages endpoint bound to 0.0.0.0"
sidebar_label: OTEL-060
description: zpages exposes pipeline internals — do not serve it to the whole world.
---

# OTEL-060: `zpages` endpoint bound to `0.0.0.0`

**Severity:** <span className="badge badge--warning">warn (advisory)</span>

## Rule Details

`zpages` serves live debugging pages that show spans in flight, recent errors, and pipeline throughput. Binding it to `0.0.0.0` makes that information reachable from any host that can hit the Collector's network, including attackers looking for a quick recon surface. Bind it to `localhost` and tunnel to it (`kubectl port-forward`, SSH) when you actually need it.

This rule fires when a `zpages` extension has an `endpoint` containing `0.0.0.0`.

## Options

This rule has no options.

## Examples

:::warning[Avoid]

```yaml
extensions:
  zpages:
    endpoint: "0.0.0.0:55679"
```

:::

:::tip[Prefer]

```yaml
extensions:
  zpages:
    endpoint: "localhost:55679"
```

:::

## When Not To Use It

Never — there is no legitimate reason to expose `zpages` on every interface.

## Related Rules

- [OTEL-010](../core/otel-010) — receivers should not bind to `0.0.0.0`
- [OTEL-059](./otel-059) — `pprof` extension enabled in production
- [OTEL-070](../lifecycle/otel-070) — telemetry metrics address bound to `0.0.0.0`

## Version

Available since augur v0.1.0.

## Further Reading

- [OpenTelemetry Collector Contrib — zpages extension](https://github.com/open-telemetry/opentelemetry-collector-contrib/blob/main/extension/zpagesextension/README.md)

## Resources

- Rule source: [`policy/main/extension.rego`](https://github.com/starkross/augur/blob/main/policy/main/extension.rego)
