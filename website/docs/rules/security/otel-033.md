---
id: otel-033
title: "OTEL-033: Receiver on non-localhost endpoint without TLS"
sidebar_label: OTEL-033
description: A remote receiver without TLS accepts cleartext telemetry from anywhere.
---

# OTEL-033: Receiver on non-localhost endpoint without TLS

**Severity:** <span className="badge badge--warning">warn (advisory)</span>

## Rule Details

If a receiver is bound to a non-localhost endpoint (an interface other than `127.0.0.1`/`localhost`), it is almost certainly reachable from another host on the network. Without a `tls` block, the receiver will happily accept cleartext OTLP/HTTP from any caller that can reach the port — so anyone between the client and the Collector can read the traffic or inject it.

This rule fires when a receiver protocol has an `endpoint` that does not contain `localhost` or `127.0.0.1` and no `tls` block is configured.

## Options

This rule has no options.

## Examples

:::warning[Avoid]

```yaml
receivers:
  otlp:
    protocols:
      grpc:
        endpoint: "0.0.0.0:4317"
        # no tls block
```

:::

:::tip[Prefer]

```yaml
receivers:
  otlp:
    protocols:
      grpc:
        endpoint: "10.0.0.5:4317"
        tls:
          min_version: "1.3"
          cert_file: /etc/certs/server.crt
          key_file: /etc/certs/server.key
```

:::

## When Not To Use It

The receiver sits behind a service mesh sidecar that terminates mTLS on its behalf. In that case the mesh provides the encryption, and you should still pin the receiver to a local interface (e.g. `127.0.0.1`) so only the sidecar can talk to it.

## Related Rules

- [OTEL-010](../core/otel-010) — receivers should not bind to `0.0.0.0`
- [OTEL-018](../core/otel-018) — OTLP exporter without TLS on non-local endpoint
- [OTEL-031](./otel-031) — TLS `min_version` below 1.2

## Version

Available since augur v0.1.0.

## Further Reading

- [OpenTelemetry Collector — TLS configuration](https://github.com/open-telemetry/opentelemetry-collector/blob/main/config/configtls/README.md)

## Resources

- Rule source: [`policy/main/security.rego`](https://github.com/starkross/augur/blob/main/policy/main/security.rego)
