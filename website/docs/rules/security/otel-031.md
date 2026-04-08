---
id: otel-031
title: "OTEL-031: TLS min_version below 1.2"
sidebar_label: OTEL-031
description: TLS 1.0 and 1.1 are obsolete — do not accept them.
---

# OTEL-031: TLS `min_version` below 1.2

**Severity:** <span className="badge badge--danger">deny (blocking)</span>

## Rule Details

TLS 1.0 and 1.1 have been deprecated by every major browser, by the IETF (RFC 8996), and by most compliance programs (PCI-DSS, HIPAA). They have known cryptographic weaknesses and MITM attacks in the wild. Any Collector receiver or exporter that accepts `min_version: "1.0"` or `"1.1"` is a downgrade surface — drop it to `1.2` at minimum, `1.3` when possible.

This rule fires when any receiver protocol TLS block or exporter TLS block has `min_version` set to `"1.0"` or `"1.1"`.

## Options

| Field | Constraint |
|------|------------|
| `tls.min_version` | Must be `"1.2"` or `"1.3"` |

## Examples

:::danger[Incorrect]

```yaml
receivers:
  otlp:
    protocols:
      grpc:
        endpoint: 0.0.0.0:4317
        tls:
          min_version: "1.0"       # obsolete
          cert_file: /etc/certs/server.crt
          key_file: /etc/certs/server.key
```

:::

:::tip[Correct]

```yaml
receivers:
  otlp:
    protocols:
      grpc:
        endpoint: 0.0.0.0:4317
        tls:
          min_version: "1.3"
          cert_file: /etc/certs/server.crt
          key_file: /etc/certs/server.key
```

:::

## When Not To Use It

Never. If a legacy client needs TLS 1.0/1.1, put a reverse proxy in front that terminates the old protocol and re-encrypts to the Collector on TLS 1.3.

## Related Rules

- [OTEL-032](./otel-032) — `insecure_skip_verify` enabled
- [OTEL-033](./otel-033) — receiver on non-localhost endpoint without TLS
- [OTEL-018](../core/otel-018) — OTLP exporter without TLS on non-local endpoint
- [OTEL-037](./otel-037) — inline `key_pem` detected

## Version

Available since augur v0.1.0.

## Further Reading

- [RFC 8996 — Deprecating TLS 1.0 and 1.1](https://www.rfc-editor.org/rfc/rfc8996.html)
- [OpenTelemetry Collector — TLS configuration](https://github.com/open-telemetry/opentelemetry-collector/blob/main/config/configtls/README.md)

## Resources

- Rule source: [`policy/main/security.rego`](https://github.com/starkross/augur/blob/main/policy/main/security.rego)
