---
id: index
title: Rules
sidebar_position: 5
sidebar_label: Rules
---

# Rules

augur ships with two severities: **deny** rules block (exit code `1`), and **warn** rules are advisory by default — promote them with `--strict`. Every rule lives under [`policy/`](https://github.com/starkross/augur/tree/main/policy) as a standalone `.rego` file you can read, override, or extend.

Rules are grouped by category in the sidebar (Core, Memory, Security, Pipeline, Exporter, Receiver, Extension, Reliability, Lifecycle). The tables below give the full index in deny-then-warn order.

## Deny (blocking)

| ID | Description |
|----|-------------|
| [OTEL-001](./core/otel-001.md) | `memory_limiter` processor must be configured |
| [OTEL-002](./core/otel-002.md) | `memory_limiter` must be included in every pipeline |
| [OTEL-003](./core/otel-003.md) | `batch` processor must be configured |
| [OTEL-004](./core/otel-004.md) | No hardcoded secrets in exporters |
| [OTEL-005](./core/otel-005.md) | No hardcoded secrets in receivers |
| [OTEL-006](./core/otel-006.md) | `service.pipelines` must be defined |
| [OTEL-007](./core/otel-007.md) | Every pipeline must have receivers and exporters |
| [OTEL-024](./memory/otel-024.md) | `batch` `send_batch_max_size` < `send_batch_size` |
| [OTEL-027](./memory/otel-027.md) | `memory_limiter` `check_interval` is 0 or unset |
| [OTEL-028](./memory/otel-028.md) | `spike_limit_mib` >= `limit_mib` (soft limit zero or negative) |
| [OTEL-029](./memory/otel-029.md) | Neither `limit_mib` nor `limit_percentage` set on `memory_limiter` |
| [OTEL-031](./security/otel-031.md) | TLS `min_version` below 1.2 |
| [OTEL-034](./security/otel-034.md) | CORS `allowed_origins` contains wildcard `*` |
| [OTEL-035](./security/otel-035.md) | Hardcoded secrets in extensions |
| [OTEL-040](./pipeline/otel-040.md) | Circular pipeline dependency via connectors |
| [OTEL-044](./exporter/otel-044.md) | OTLP gRPC exporter endpoint has `http(s)://` scheme (use bare `host:port`) |
| [OTEL-058](./receiver/otel-058.md) | Multiple receivers bound to the same endpoint |
| [OTEL-066](./extension/otel-066.md) | `sending_queue.storage` references undefined extension |

## Warn (advisory)

| ID | Description |
|----|-------------|
| [OTEL-010](./core/otel-010.md) | Receivers should not bind to `0.0.0.0` |
| [OTEL-011](./core/otel-011.md) | `health_check` extension recommended |
| [OTEL-012](./core/otel-012.md) | `health_check` configured but not listed in `service.extensions` |
| [OTEL-013](./core/otel-013.md) | `batch` processor should be last in pipeline |
| [OTEL-014](./core/otel-014.md) | `memory_limiter` should be first processor in pipeline |
| [OTEL-015](./core/otel-015.md) | `debug`/`logging` exporter detected |
| [OTEL-016](./core/otel-016.md) | Telemetry log level set to `debug` |
| [OTEL-017](./core/otel-017.md) | Exporter missing `retry_on_failure`/`sending_queue` |
| [OTEL-018](./core/otel-018.md) | OTLP exporter without TLS on non-local endpoint |
| [OTEL-020](./core/otel-020.md) | Unused receiver |
| [OTEL-021](./core/otel-021.md) | Unused exporter |
| [OTEL-022](./core/otel-022.md) | Unused processor |
| [OTEL-023](./memory/otel-023.md) | `batch` `send_batch_max_size` unset (unlimited) |
| [OTEL-025](./memory/otel-025.md) | `batch` timeout below 100ms |
| [OTEL-026](./memory/otel-026.md) | `batch` timeout above 60s |
| [OTEL-030](./memory/otel-030.md) | `memory_limiter` `limit_percentage` outside safe range (20–90%) |
| [OTEL-032](./security/otel-032.md) | `insecure_skip_verify` enabled |
| [OTEL-033](./security/otel-033.md) | Receiver on non-localhost endpoint without TLS |
| [OTEL-036](./security/otel-036.md) | gRPC `max_recv_msg_size_mib` > 128 (decompression bomb risk) |
| [OTEL-037](./security/otel-037.md) | Inline `key_pem` detected (use `key_file` instead) |
| [OTEL-038](./pipeline/otel-038.md) | Filter processor after batch (filter early to reduce waste) |
| [OTEL-039](./pipeline/otel-039.md) | Transform/attributes processor after batch |
| [OTEL-041](./pipeline/otel-041.md) | Routing connector without `default_pipelines` |
| [OTEL-042](./pipeline/otel-042.md) | Duplicate processor in same pipeline |
| [OTEL-043](./pipeline/otel-043.md) | Batch before `tail_sampling`/`groupbytrace` |
| [OTEL-045](./exporter/otel-045.md) | OTLP gRPC endpoint missing port number |
| [OTEL-046](./exporter/otel-046.md) | OTLP HTTP endpoint missing URL scheme |
| [OTEL-047](./exporter/otel-047.md) | OTLP HTTP exporter using gRPC port 4317 (HTTP is 4318) |
| [OTEL-048](./exporter/otel-048.md) | `sending_queue` explicitly disabled |
| [OTEL-049](./exporter/otel-049.md) | `sending_queue.queue_size` below 10 |
| [OTEL-050](./exporter/otel-050.md) | `sending_queue.queue_size` above 50000 (OOM risk) |
| [OTEL-051](./exporter/otel-051.md) | `sending_queue.num_consumers` below 2 |
| [OTEL-052](./exporter/otel-052.md) | Compression disabled for network exporter |
| [OTEL-053](./exporter/otel-053.md) | Retry `max_elapsed_time` set to 0 (infinite retries) |
| [OTEL-054](./receiver/otel-054.md) | Prometheus `scrape_interval` below 10s |
| [OTEL-055](./receiver/otel-055.md) | `hostmetrics` `collection_interval` below 10s |
| [OTEL-056](./receiver/otel-056.md) | `filelog` `start_at:beginning` without storage |
| [OTEL-057](./receiver/otel-057.md) | `filelog` overly broad include pattern |
| [OTEL-059](./extension/otel-059.md) | `pprof` extension enabled in production |
| [OTEL-060](./extension/otel-060.md) | `zpages` endpoint bound to `0.0.0.0` |
| [OTEL-061](./extension/otel-061.md) | `memory_ballast` extension (deprecated, use `GOMEMLIMIT`) |
| [OTEL-062](./extension/otel-062.md) | Extension in `service.extensions` but not defined |
| [OTEL-063](./reliability/otel-063.md) | `tail_sampling` without `groupbytrace` |
| [OTEL-064](./reliability/otel-064.md) | Both `probabilistic_sampler` and `tail_sampling` in same pipeline |
| [OTEL-065](./reliability/otel-065.md) | `sending_queue` without persistent storage |
| [OTEL-067](./reliability/otel-067.md) | K8s environment without `k8sattributes` processor |
| [OTEL-068](./reliability/otel-068.md) | K8s environment without `resourcedetection` processor |
| [OTEL-069](./lifecycle/otel-069.md) | Telemetry metrics level set to `none` |
| [OTEL-070](./lifecycle/otel-070.md) | Telemetry metrics address bound to `0.0.0.0` |
| [OTEL-071](./lifecycle/otel-071.md) | `logging` exporter deprecated (renamed to `debug` in v0.111.0) |
| [OTEL-072](./lifecycle/otel-072.md) | OpenCensus receiver/exporter (deprecated, migrate to OTLP) |
| [OTEL-073](./lifecycle/otel-073.md) | `memory_limiter` `ballast_size_mib` (deprecated, use `GOMEMLIMIT`) |
| [OTEL-074](./lifecycle/otel-074.md) | `service.telemetry.metrics.address` (deprecated, use `readers` config) |
