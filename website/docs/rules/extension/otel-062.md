---
id: otel-062
title: "OTEL-062: Extension in service.extensions but not defined"
sidebar_label: OTEL-062
description: Referencing an extension in service.extensions without defining it fails at startup.
---

# OTEL-062: Extension in `service.extensions` but not defined

**Severity:** <span className="badge badge--warning">warn (advisory)</span>

## Rule Details

`service.extensions` is the list of enabled extensions. Each entry must refer to a block that actually exists in the top-level `extensions:` section. If it does not, the Collector aborts startup with a cryptic "unknown extension" error. This rule catches stray references — usually left over from renaming or removing an extension without updating the service block.

This rule fires when a name in `service.extensions` does not have a matching entry in `extensions:`.

## Options

This rule has no options.

## Examples

:::warning[Avoid]

```yaml
extensions:
  health_check:
    endpoint: "localhost:13133"

service:
  extensions: [health_check, pprof]    # pprof not defined
```

:::

:::tip[Prefer]

```yaml
extensions:
  health_check:
    endpoint: "localhost:13133"

service:
  extensions: [health_check]
```

:::

## When Not To Use It

Never — this configuration cannot start.

## Related Rules

- [OTEL-011](../core/otel-011) — `health_check` extension recommended
- [OTEL-012](../core/otel-012) — `health_check` configured but not listed in `service.extensions`

## Version

Available since augur v0.1.0.

## Further Reading

- [OpenTelemetry Collector — extensions](https://opentelemetry.io/docs/collector/configuration/#extensions)

## Resources

- Rule source: [`policy/main/extension.rego`](https://github.com/starkross/augur/blob/main/policy/main/extension.rego)
