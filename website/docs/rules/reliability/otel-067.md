---
id: otel-067
title: "OTEL-067: K8s environment without k8sattributes processor"
sidebar_label: OTEL-067
description: Kubernetes receivers produce telemetry with no pod/namespace metadata unless k8sattributes is configured.
---

# OTEL-067: K8s environment without `k8sattributes` processor

**Severity:** <span className="badge badge--warning">warn (advisory)</span>

## Rule Details

Receivers like `kubeletstats`, `k8s_cluster`, `k8s_events`, and `k8sobjects` pull telemetry from Kubernetes APIs — but on their own they tag each item with the source component, not the pod, namespace, deployment, or node the metric actually describes. The `k8sattributes` processor enriches each item with that metadata so you can actually slice dashboards by workload.

This rule fires when a Kubernetes receiver is configured but no `k8sattributes` processor exists.

## Options

This rule has no options.

## Examples

:::warning[Avoid]

```yaml
receivers:
  kubeletstats:
    collection_interval: 30s
    auth_type: serviceAccount

processors:
  memory_limiter:
    check_interval: 5s
    limit_percentage: 80
  batch:
    timeout: 1s
# no k8sattributes
```

:::

:::tip[Prefer]

```yaml
receivers:
  kubeletstats:
    collection_interval: 30s
    auth_type: serviceAccount

processors:
  memory_limiter:
    check_interval: 5s
    limit_percentage: 80
  k8sattributes:
    auth_type: serviceAccount
    passthrough: false
  batch:
    timeout: 1s
```

:::

## When Not To Use It

A pipeline that feeds a backend which already enriches with Kubernetes metadata on ingest. In that case document the assumption explicitly.

## Related Rules

- [OTEL-068](./otel-068) — K8s environment without `resourcedetection` processor

## Version

Available since augur v0.1.0.

## Further Reading

- [OpenTelemetry Collector Contrib — k8sattributes processor](https://github.com/open-telemetry/opentelemetry-collector-contrib/blob/main/processor/k8sattributesprocessor/README.md)

## Resources

- Rule source: [`policy/main/reliability.rego`](https://github.com/starkross/augur/blob/main/policy/main/reliability.rego)
