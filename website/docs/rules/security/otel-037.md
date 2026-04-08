---
id: otel-037
title: "OTEL-037: Inline key_pem detected (use key_file instead)"
sidebar_label: OTEL-037
description: Inlining a private key in the config file puts it into source control.
---

# OTEL-037: Inline `key_pem` detected (use `key_file` instead)

**Severity:** <span className="badge badge--warning">warn (advisory)</span>

## Rule Details

`tls.key_pem` takes the contents of a private key directly in the config file — which usually means the PEM block ends up in git, in a container image layer, or in a Kubernetes ConfigMap. Use `tls.key_file` with a path instead and mount the key via a Secret, CSI driver, or a proper secret store.

This rule fires when any exporter or receiver protocol has `tls.key_pem` set.

## Options

This rule has no options.

## Examples

:::warning[Avoid]

```yaml
exporters:
  otlp/backend:
    endpoint: backend:4317
    tls:
      cert_pem: |
        -----BEGIN CERTIFICATE-----
        ...
      key_pem: |                            # private key inlined
        -----BEGIN PRIVATE KEY-----
        ...
```

:::

:::tip[Prefer]

```yaml
exporters:
  otlp/backend:
    endpoint: backend:4317
    tls:
      cert_file: /etc/certs/client.crt
      key_file: /etc/certs/client.key
```

:::

## When Not To Use It

Never. Private key material should never sit in a config file — use an out-of-band secret store.

## Related Rules

- [OTEL-004](../core/otel-004) — no hardcoded secrets in exporters
- [OTEL-005](../core/otel-005) — no hardcoded secrets in receivers
- [OTEL-031](./otel-031) — TLS `min_version` below 1.2
- [OTEL-032](./otel-032) — `insecure_skip_verify` enabled

## Version

Available since augur v0.1.0.

## Further Reading

- [OpenTelemetry Collector — TLS configuration](https://github.com/open-telemetry/opentelemetry-collector/blob/main/config/configtls/README.md)

## Resources

- Rule source: [`policy/main/security.rego`](https://github.com/starkross/augur/blob/main/policy/main/security.rego)
