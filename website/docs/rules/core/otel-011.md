---
id: otel-011
title: "OTEL-011: health_check extension recommended"
sidebar_label: OTEL-011
description: Kubernetes probes and load balancers need a health endpoint.
---

# OTEL-011: `health_check` extension recommended

**Severity:** <span className="badge badge--warning">warn (advisory)</span>

## Rule Details

Without the `health_check` extension there is no way for a liveness probe, readiness probe, or external load balancer to know whether the Collector is actually running. Pods hang on to traffic they cannot export, crash loops go undetected, and deploys look green when they are not. The extension is small, cheap, and universally useful.

This rule fires when the `extensions.health_check` block is missing.

## Options

This rule has no options.

## Examples

:::warning[Avoid]

```yaml
extensions:
  pprof:
    endpoint: localhost:1777
# no health_check
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

Throwaway local runs where you are not using Kubernetes probes or any external health checker. In production, always enable it.

## Related Rules

- [OTEL-012](./otel-012) — `health_check` configured but not listed in `service.extensions`
- [OTEL-062](../extension/otel-062) — extension in `service.extensions` but not defined

## Version

Available since augur v0.1.0.

## Further Reading

- [OpenTelemetry Collector Contrib — health_check extension](https://github.com/open-telemetry/opentelemetry-collector-contrib/blob/main/extension/healthcheckextension/README.md)

## Resources

- Rule source: [`policy/main/main.rego`](https://github.com/starkross/augur/blob/main/policy/main/main.rego)
