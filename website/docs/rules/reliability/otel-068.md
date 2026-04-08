---
id: otel-068
title: "OTEL-068: K8s environment without resourcedetection processor"
sidebar_label: OTEL-068
description: Without resourcedetection you lose host, cloud, and node identity on every item.
---

# OTEL-068: K8s environment without `resourcedetection` processor

**Severity:** <span className="badge badge--warning">warn (advisory)</span>

## Rule Details

`resourcedetection` adds resource attributes describing the environment the Collector runs in — cloud provider, region, availability zone, node name, container ID. Without it, telemetry from a Kubernetes receiver lacks the context you need to correlate issues to a specific node or availability zone, which is typically the first question during an incident.

This rule fires when a Kubernetes receiver is configured but no `resourcedetection` processor exists.

## Options

This rule has no options.

## Examples

:::warning[Avoid]

```yaml
receivers:
  k8s_cluster:
    auth_type: serviceAccount

processors:
  memory_limiter:
    check_interval: 5s
    limit_percentage: 80
  batch:
    timeout: 1s
# no resourcedetection
```

:::

:::tip[Prefer]

```yaml
receivers:
  k8s_cluster:
    auth_type: serviceAccount

processors:
  memory_limiter:
    check_interval: 5s
    limit_percentage: 80
  resourcedetection:
    detectors: [env, system, gcp, aws]
    timeout: 2s
  batch:
    timeout: 1s
```

:::

## When Not To Use It

Collectors running on bare-metal hosts with externally-injected resource attributes (via `OTEL_RESOURCE_ATTRIBUTES`). Document this so the missing processor is intentional.

## Related Rules

- [OTEL-067](./otel-067) — K8s environment without `k8sattributes` processor

## Version

Available since augur v0.1.0.

## Further Reading

- [OpenTelemetry Collector Contrib — resourcedetection processor](https://github.com/open-telemetry/opentelemetry-collector-contrib/blob/main/processor/resourcedetectionprocessor/README.md)

## Resources

- Rule source: [`policy/main/reliability.rego`](https://github.com/starkross/augur/blob/main/policy/main/reliability.rego)
