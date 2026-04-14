# Rules

## Deny (blocking)

| ID | Description |
|----|-------------|
| OTEL-001 | `memory_limiter` processor must be configured |
| OTEL-002 | `memory_limiter` must be included in every pipeline |
| OTEL-003 | `batch` processor must be configured |
| OTEL-004 | No hardcoded secrets in exporters |
| OTEL-005 | No hardcoded secrets in receivers |
| OTEL-006 | `service.pipelines` must be defined |
| OTEL-007 | Every pipeline must have receivers and exporters |
| OTEL-024 | `batch` `send_batch_max_size` < `send_batch_size` |
| OTEL-027 | `memory_limiter` `check_interval` is 0 or unset |
| OTEL-028 | `spike_limit_mib` >= `limit_mib` (soft limit zero or negative) |
| OTEL-029 | Neither `limit_mib` nor `limit_percentage` set on `memory_limiter` |
| OTEL-031 | TLS `min_version` below 1.2 |
| OTEL-034 | CORS `allowed_origins` contains wildcard `*` |
| OTEL-035 | Hardcoded secrets in extensions |
| OTEL-040 | Circular pipeline dependency via connectors |
| OTEL-044 | OTLP gRPC exporter endpoint has `http(s)://` scheme (use bare `host:port`) |
| OTEL-058 | Multiple receivers bound to the same endpoint |
| OTEL-066 | `sending_queue.storage` references undefined extension |

## Warn (advisory)

| ID | Description |
|----|-------------|
| OTEL-010 | Receivers should not bind to `0.0.0.0` |
| OTEL-011 | `health_check` extension recommended |
| OTEL-012 | `health_check` configured but not listed in `service.extensions` |
| OTEL-013 | `batch` processor should be last in pipeline |
| OTEL-014 | `memory_limiter` should be first processor in pipeline |
| OTEL-015 | `debug`/`logging` exporter detected |
| OTEL-016 | Telemetry log level set to `debug` |
| OTEL-017 | Exporter missing `retry_on_failure`/`sending_queue`/`max_retries` |
| OTEL-018 | OTLP exporter without TLS on non-local endpoint |
| OTEL-020 | Unused receiver |
| OTEL-021 | Unused exporter |
| OTEL-022 | Unused processor |
| OTEL-023 | `batch` `send_batch_max_size` unset (unlimited) |
| OTEL-025 | `batch` timeout below 100ms |
| OTEL-026 | `batch` timeout above 60s |
| OTEL-030 | `memory_limiter` `limit_percentage` outside safe range (20–90%) |
| OTEL-032 | `insecure_skip_verify` enabled |
| OTEL-033 | Receiver on non-localhost endpoint without TLS |
| OTEL-036 | gRPC `max_recv_msg_size_mib` > 128 (decompression bomb risk) |
| OTEL-037 | Inline `key_pem` detected (use `key_file` instead) |
| OTEL-038 | Filter processor after batch (filter early to reduce waste) |
| OTEL-039 | Transform/attributes processor after batch |
| OTEL-041 | Routing connector without `default_pipelines` |
| OTEL-042 | Duplicate processor in same pipeline |
| OTEL-043 | Batch before `tail_sampling`/`groupbytrace` |
| OTEL-045 | OTLP gRPC endpoint missing port number |
| OTEL-046 | OTLP HTTP endpoint missing URL scheme |
| OTEL-047 | OTLP HTTP exporter using gRPC port 4317 (HTTP is 4318) |
| OTEL-048 | `sending_queue` explicitly disabled |
| OTEL-049 | `sending_queue.queue_size` below 10 |
| OTEL-050 | `sending_queue.queue_size` above 50000 (OOM risk) |
| OTEL-051 | `sending_queue.num_consumers` below 2 |
| OTEL-052 | Compression disabled for network exporter |
| OTEL-053 | Retry `max_elapsed_time` set to 0 (infinite retries) |
| OTEL-054 | Prometheus `scrape_interval` below 10s |
| OTEL-055 | `hostmetrics` `collection_interval` below 10s |
| OTEL-056 | `filelog` `start_at:beginning` without storage |
| OTEL-057 | `filelog` overly broad include pattern |
| OTEL-059 | `pprof` extension enabled in production |
| OTEL-060 | `zpages` endpoint bound to `0.0.0.0` |
| OTEL-061 | `memory_ballast` extension (deprecated, use `GOMEMLIMIT`) |
| OTEL-062 | Extension in `service.extensions` but not defined |
| OTEL-063 | `tail_sampling` without `groupbytrace` |
| OTEL-064 | Both `probabilistic_sampler` and `tail_sampling` in same pipeline |
| OTEL-065 | `sending_queue` without persistent storage |
| OTEL-067 | K8s environment without `k8sattributes` processor |
| OTEL-068 | K8s environment without `resourcedetection` processor |
| OTEL-069 | Telemetry metrics level set to `none` |
| OTEL-070 | Telemetry metrics address bound to `0.0.0.0` |
| OTEL-071 | `logging` exporter deprecated (renamed to `debug` in v0.111.0) |
| OTEL-072 | OpenCensus receiver/exporter (deprecated, migrate to OTLP) |
| OTEL-073 | `memory_limiter` `ballast_size_mib` (deprecated, use `GOMEMLIMIT`) |
| OTEL-074 | `service.telemetry.metrics.address` (deprecated, use `readers` config) |