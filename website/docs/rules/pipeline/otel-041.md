---
id: otel-041
title: "OTEL-041: Routing connector without default_pipelines"
sidebar_label: OTEL-041
description: Unmatched routes are silently dropped when there is no default.
---

# OTEL-041: Routing connector without `default_pipelines`

**Severity:** <span className="badge badge--warning">warn (advisory)</span>

## Rule Details

The `routing` connector forwards telemetry to different pipelines based on attribute values. If no `default_pipelines` is configured, any item that does not match a route is silently discarded. That is extremely easy to miss — the config passes validation, traffic flows for every route you tested, and the one route you forgot disappears. Always set a default (even if it is a `debug` pipeline you can grep).

This rule fires when a `routing` connector has no `default_pipelines` field.

## Options

This rule has no options.

## Examples

:::warning[Avoid]

```yaml
connectors:
  routing:
    table:
      - statement: route() where attributes["env"] == "prod"
        pipelines: [traces/prod]
    # no default_pipelines
```

:::

:::tip[Prefer]

```yaml
connectors:
  routing:
    default_pipelines: [traces/fallback]
    table:
      - statement: route() where attributes["env"] == "prod"
        pipelines: [traces/prod]
```

:::

## When Not To Use It

You deliberately want unmatched telemetry dropped. In that case set `default_pipelines: []` explicitly so the intent is documented.

## Related Rules

- [OTEL-040](./otel-040) — circular pipeline dependency via connectors
- [OTEL-021](../core/otel-021) — unused exporter

## Version

Available since augur v0.1.0.

## Further Reading

- [OpenTelemetry Collector Contrib — routing connector](https://github.com/open-telemetry/opentelemetry-collector-contrib/blob/main/connector/routingconnector/README.md)

## Resources

- Rule source: [`policy/main/pipeline.rego`](https://github.com/starkross/augur/blob/main/policy/main/pipeline.rego)
