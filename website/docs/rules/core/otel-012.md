---
id: otel-012
title: "OTEL-012: health_check configured but not listed in service.extensions"
sidebar_label: OTEL-012
description: Defining an extension does not enable it — it also has to appear in service.extensions.
---

# OTEL-012: `health_check` configured but not listed in `service.extensions`

**Severity:** <span className="badge badge--warning">warn (advisory)</span>

## Rule Details

In the Collector config, an extension has to be declared in `extensions:` *and* referenced in `service.extensions:` before it starts. It is very easy to write the `health_check` block and then forget the second step — the config validates, the extension does nothing, and probes silently fail.

This rule fires when `extensions.health_check` exists but `health_check` is not present in `service.extensions`.

## Options

This rule has no options.

## Examples

:::warning[Avoid]

```yaml
extensions:
  health_check:
    endpoint: "localhost:13133"

service:
  extensions: []            # health_check not enabled
  pipelines:
    traces:
      receivers: [otlp]
      exporters: [otlp/backend]
```

:::

:::tip[Prefer]

```yaml
extensions:
  health_check:
    endpoint: "localhost:13133"

service:
  extensions: [health_check]
  pipelines:
    traces:
      receivers: [otlp]
      exporters: [otlp/backend]
```

:::

## When Not To Use It

Never — this rule catches a near-universal copy/paste mistake. If you do not want `health_check` running, delete the `extensions.health_check` block entirely.

## Related Rules

- [OTEL-011](./otel-011) — `health_check` extension recommended
- [OTEL-062](../extension/otel-062) — extension in `service.extensions` but not defined

## Version

Available since augur v0.1.0.

## Further Reading

- [OpenTelemetry Collector — extensions](https://opentelemetry.io/docs/collector/configuration/#extensions)

## Resources

- Rule source: [`policy/main/main.rego`](https://github.com/starkross/augur/blob/main/policy/main/main.rego)
