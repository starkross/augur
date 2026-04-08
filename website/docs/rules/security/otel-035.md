---
id: otel-035
title: "OTEL-035: Hardcoded secrets in extensions"
sidebar_label: OTEL-035
description: Extensions like basicauth must read credentials from the environment.
---

# OTEL-035: Hardcoded secrets in extensions

**Severity:** <span className="badge badge--danger">deny (blocking)</span>

## Rule Details

Extensions such as `basicauth`, `oauth2client`, `oidc`, and `bearertokenauth` carry credentials. Writing those credentials inline in the config leaks them into git, into pod manifests, and into CI logs. augur scans the `extensions:` block — including one level of nesting for blocks like `client_auth:` — and blocks any secret-like field whose value is not an environment-variable reference.

This rule fires when an extension (or a nested sub-object inside an extension) has a secret-like field whose value is a plain string rather than `${env:VAR_NAME}`.

## Options

This rule has no options.

## Examples

:::danger[Incorrect]

```yaml
extensions:
  basicauth/server:
    htpasswd:
      inline: |
        admin:$2y$05$abcdefg              # literal secret
  oauth2client:
    client_auth:
      client_secret: "literal-client-secret"
```

:::

:::tip[Correct]

```yaml
extensions:
  basicauth/server:
    htpasswd:
      file: /etc/otel/htpasswd
  oauth2client:
    client_auth:
      client_secret: "${env:OAUTH_CLIENT_SECRET}"
```

:::

## When Not To Use It

Never. Treat extension credentials the same as exporter and receiver credentials — environment variables only.

## Related Rules

- [OTEL-004](../core/otel-004) — no hardcoded secrets in exporters
- [OTEL-005](../core/otel-005) — no hardcoded secrets in receivers

## Version

Available since augur v0.1.0.

## Further Reading

- [OpenTelemetry Collector — configuration environment variables](https://opentelemetry.io/docs/collector/configuration/#environment-variables)

## Resources

- Rule source: [`policy/main/security.rego`](https://github.com/starkross/augur/blob/main/policy/main/security.rego)
