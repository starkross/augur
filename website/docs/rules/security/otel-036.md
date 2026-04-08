---
id: otel-036
title: "OTEL-036: gRPC max_recv_msg_size_mib > 128 (decompression bomb risk)"
sidebar_label: OTEL-036
description: A very large receive size turns any bad client into a memory attack.
---

# OTEL-036: gRPC `max_recv_msg_size_mib` > 128 (decompression bomb risk)

**Severity:** <span className="badge badge--warning">warn (advisory)</span>

## Rule Details

`max_recv_msg_size_mib` caps the size of an incoming gRPC message after decompression. Raising it to hundreds of megabytes makes the receiver trivially exploitable with a decompression bomb: a tiny highly-compressed payload that blows up into gigabytes of in-memory protobuf. 128 MiB is already generous; anything larger should be justified by a specific client constraint and paired with an upstream rate limiter.

This rule fires when a receiver protocol has `max_recv_msg_size_mib > 128`.

## Options

| Field | Constraint |
|------|------------|
| `max_recv_msg_size_mib` | Should be ≤ 128 |

## Examples

:::warning[Avoid]

```yaml
receivers:
  otlp:
    protocols:
      grpc:
        endpoint: localhost:4317
        max_recv_msg_size_mib: 1024     # 1 GiB — decompression bomb surface
```

:::

:::tip[Prefer]

```yaml
receivers:
  otlp:
    protocols:
      grpc:
        endpoint: localhost:4317
        max_recv_msg_size_mib: 64
```

:::

## When Not To Use It

Trusted, authenticated clients that genuinely need to push very large batches (e.g., offline log backfills). Gate the receiver behind auth so the large size is not reachable by strangers.

## Related Rules

- [OTEL-001](../core/otel-001) — `memory_limiter` processor must be configured
- [OTEL-033](./otel-033) — receiver on non-localhost endpoint without TLS

## Version

Available since augur v0.1.0.

## Further Reading

- [OpenTelemetry Collector — OTLP receiver](https://github.com/open-telemetry/opentelemetry-collector/blob/main/receiver/otlpreceiver/README.md)

## Resources

- Rule source: [`policy/main/security.rego`](https://github.com/starkross/augur/blob/main/policy/main/security.rego)
