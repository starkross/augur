package main

import data.lib
import future.keywords.contains
import future.keywords.if
import future.keywords.in

# ============================================================
# DENY — blocking rules (must fix before deploy)
# ============================================================

# METADATA
# entrypoint: true
deny contains msg if {
	not input.processors.memory_limiter
	msg := "OTEL-001: memory_limiter processor is not configured. Required to prevent OOM in production."
}

deny contains msg if {
	input.processors.memory_limiter
	some ptype in lib.pipeline_types
	procs := lib.pipeline_processors(ptype)
	not "memory_limiter" in procs
	msg := sprintf("OTEL-002: pipeline '%s' does not include memory_limiter processor.", [ptype])
}

deny contains msg if {
	not input.processors.batch
	msg := "OTEL-003: batch processor is not configured. Required for efficient data export."
}

deny contains msg if {
	some name, exporter in input.exporters
	some k, v in exporter
	is_string(v)
	lib.looks_like_secret(k)
	not lib.is_env_var(v)
	msg := sprintf("OTEL-004: exporter '%s' has hardcoded '%s'. Use ${env:VAR_NAME} instead.", [name, k])
}

deny contains msg if {
	some name, receiver in input.receivers
	some k, v in receiver
	is_string(v)
	lib.looks_like_secret(k)
	not lib.is_env_var(v)
	msg := sprintf("OTEL-005: receiver '%s' has hardcoded '%s'. Use ${env:VAR_NAME} instead.", [name, k])
}

deny contains msg if {
	not input.service.pipelines
	msg := "OTEL-006: service.pipelines is not defined. At least one pipeline is required."
}

deny contains msg if {
	some ptype in lib.pipeline_types
	rs := lib.pipeline_receivers(ptype)
	count(rs) == 0
	msg := sprintf("OTEL-007: pipeline '%s' has no receivers.", [ptype])
}

deny contains msg if {
	some ptype in lib.pipeline_types
	es := lib.pipeline_exporters(ptype)
	count(es) == 0
	msg := sprintf("OTEL-007: pipeline '%s' has no exporters.", [ptype])
}

# ============================================================
# WARN — advisory rules (best practices)
# ============================================================

# METADATA
# entrypoint: true
warn contains msg if {
	some name, receiver in input.receivers
	lib.obj_contains_string(receiver, "0.0.0.0")
	msg := sprintf("OTEL-010: receiver '%s' binds to 0.0.0.0. Use localhost or a specific interface for security.", [name])
}

warn contains msg if {
	not input.extensions.health_check
	msg := "OTEL-011: health_check extension is not configured. Recommended for k8s liveness/readiness probes."
}

warn contains msg if {
	input.extensions.health_check
	extensions := object.get(input.service, "extensions", [])
	not "health_check" in extensions
	msg := "OTEL-012: health_check is configured but not listed in service.extensions."
}

warn contains msg if {
	some ptype in lib.pipeline_types
	procs := lib.pipeline_processors(ptype)
	count(procs) > 1
	"batch" in procs
	procs[count(procs) - 1] != "batch"
	msg := sprintf("OTEL-013: pipeline '%s' — batch processor should be last for optimal performance.", [ptype])
}

warn contains msg if {
	some ptype in lib.pipeline_types
	procs := lib.pipeline_processors(ptype)
	count(procs) > 1
	"memory_limiter" in procs
	procs[0] != "memory_limiter"
	msg := sprintf("OTEL-014: pipeline '%s' — memory_limiter should be first processor.", [ptype])
}

warn contains msg if {
	some name, _ in input.exporters
	name in {"debug", "logging"}
	msg := sprintf("OTEL-015: '%s' exporter detected. Remove or disable in production.", [name])
}

warn contains msg if {
	input.service.telemetry.logs.level == "debug"
	msg := "OTEL-016: service telemetry log level is 'debug'. Use 'info' or 'warn' in production."
}

warn contains msg if {
	pull_based := {"debug", "logging", "prometheus", "prometheusremotewrite"}
	some name, exporter in input.exporters
	base_type := split(name, "/")[0]
	not base_type in pull_based
	not exporter.retry_on_failure
	not exporter.sending_queue
	not _exporter_has_alt_retry(base_type, exporter)
	msg := sprintf("OTEL-017: exporter '%s' has no retry_on_failure, sending_queue, or exporter-native retry (e.g. max_retries). Risk of data loss.", [name])
}

# AWS exporters that implement their own retry via max_retries. These exporters
# manage retry internally and do not require sending_queue; omitting it is
# intentional. Matches the allowlist pattern used for pull_based exporters above.
_aws_alt_retry_exporters := {"awsemf", "awscloudwatchlogs", "awsxray", "awss3"}

# Helper: recognises exporter-specific alternative retry mechanisms.
_exporter_has_alt_retry(base_type, exporter) if {
	base_type in _aws_alt_retry_exporters
	exporter.max_retries
}

warn contains msg if {
	some name, exporter in input.exporters
	startswith(name, "otlp")
	not exporter.tls
	endpoint := object.get(exporter, "endpoint", "")
	not startswith(endpoint, "https://")
	not contains(endpoint, "localhost")
	not contains(endpoint, "127.0.0.1")
	msg := sprintf("OTEL-018: exporter '%s' has no TLS configured for non-local endpoint.", [name])
}

warn contains msg if {
	some name, _ in input.receivers
	not name in lib.all_used_receivers
	msg := sprintf("OTEL-020: receiver '%s' is configured but not used in any pipeline.", [name])
}

warn contains msg if {
	some name, _ in input.exporters
	not name in lib.all_used_exporters
	msg := sprintf("OTEL-021: exporter '%s' is configured but not used in any pipeline.", [name])
}

warn contains msg if {
	some name, _ in input.processors
	not name in lib.all_used_processors
	msg := sprintf("OTEL-022: processor '%s' is configured but not used in any pipeline.", [name])
}
