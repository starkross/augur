---
id: quick-start
title: Quick start
sidebar_position: 3
---

# Quick start

Point `augur` at a collector config and it'll surface any problems:

```sh
augur otel-collector-config.yaml
```

Sample output:

```text
otel-collector-config.yaml
  FAIL OTEL-001: memory_limiter processor is not configured. Required to prevent OOM in production.
  FAIL OTEL-003: batch processor is not configured. Required for efficient data export.
  WARN OTEL-011: health_check extension is not configured. Recommended for k8s liveness/readiness probes.

✗ 2 failure(s), 1 warning(s)
```

Exit code `1` on any failure. Warnings are informational by default — pass `--strict` to promote them to failures.

## Next steps

- Review every [rule](./rules) augur can enforce.
- Read [Usage](./usage) for flag reference and integration examples.
- Write [custom policies](./custom-policies) tuned to your platform.
