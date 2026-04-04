package main

import data.lib
import future.keywords.contains
import future.keywords.if
import future.keywords.in

# Exporter configuration rules

# OTEL-044: OTLP gRPC exporter endpoint has http(s):// prefix
deny contains msg if {
	some name, exporter in input.exporters
	split(name, "/")[0] == "otlp"
	endpoint := object.get(exporter, "endpoint", "")
	_has_http_scheme(endpoint)
	msg := sprintf("OTEL-044: OTLP gRPC exporter '%s' endpoint '%s' has URL scheme. Use bare host:port.", [name, endpoint])
}

# OTEL-045: OTLP gRPC endpoint missing port
warn contains msg if {
	some name, exporter in input.exporters
	split(name, "/")[0] == "otlp"
	endpoint := object.get(exporter, "endpoint", "")
	endpoint != ""
	not contains(endpoint, ":")
	msg := sprintf("OTEL-045: OTLP gRPC exporter '%s' endpoint '%s' is missing port number.", [name, endpoint])
}

# OTEL-046: OTLP HTTP endpoint missing URL scheme
warn contains msg if {
	some name, exporter in input.exporters
	split(name, "/")[0] == "otlphttp"
	endpoint := object.get(exporter, "endpoint", "")
	endpoint != ""
	not startswith(endpoint, "http://")
	not startswith(endpoint, "https://")
	not lib.is_env_var(endpoint)
	msg := sprintf("OTEL-046: OTLP HTTP exporter '%s' endpoint '%s' is missing URL scheme.", [name, endpoint])
}

# OTEL-047: OTLP HTTP exporter using gRPC port 4317
warn contains msg if {
	some name, exporter in input.exporters
	split(name, "/")[0] == "otlphttp"
	endpoint := object.get(exporter, "endpoint", "")
	endswith(endpoint, ":4317")
	msg := sprintf("OTEL-047: OTLP HTTP exporter '%s' uses gRPC port 4317. HTTP standard port is 4318.", [name])
}

warn contains msg if {
	some name, exporter in input.exporters
	split(name, "/")[0] == "otlphttp"
	endpoint := object.get(exporter, "endpoint", "")
	contains(endpoint, ":4317/")
	msg := sprintf("OTEL-047: OTLP HTTP exporter '%s' uses gRPC port 4317. HTTP standard port is 4318.", [name])
}

# OTEL-048: sending_queue explicitly disabled
warn contains msg if {
	pull_based := {"debug", "logging", "prometheus", "prometheusremotewrite", "file"}
	some name, exporter in input.exporters
	not split(name, "/")[0] in pull_based
	exporter.sending_queue.enabled == false
	msg := sprintf("OTEL-048: exporter '%s' has sending_queue disabled. Risk of data loss on export failures.", [name])
}

# OTEL-049: queue_size below 10
warn contains msg if {
	some name, exporter in input.exporters
	exporter.sending_queue.queue_size < 10
	msg := sprintf(
		"OTEL-049: exporter '%s' sending_queue.queue_size is %d (<10). Queue fills instantly.",
		[name, exporter.sending_queue.queue_size],
	)
}

# OTEL-050: queue_size above 50000
warn contains msg if {
	some name, exporter in input.exporters
	exporter.sending_queue.queue_size > 50000
	msg := sprintf(
		"OTEL-050: exporter '%s' sending_queue.queue_size is %d (>50000). Risk of OOM outside memory_limiter.",
		[name, exporter.sending_queue.queue_size],
	)
}

# OTEL-051: num_consumers below 2
warn contains msg if {
	some name, exporter in input.exporters
	exporter.sending_queue.num_consumers < 2
	msg := sprintf(
		"OTEL-051: exporter '%s' sending_queue.num_consumers is %d. Single consumer limits throughput.",
		[name, exporter.sending_queue.num_consumers],
	)
}

# OTEL-052: compression disabled for network exporter
warn contains msg if {
	pull_based := {"debug", "logging", "prometheus", "prometheusremotewrite", "file"}
	some name, exporter in input.exporters
	not split(name, "/")[0] in pull_based
	exporter.compression == "none"
	msg := sprintf("OTEL-052: exporter '%s' has compression disabled. gzip reduces bandwidth by 70-90%%.", [name])
}

# OTEL-053: retry max_elapsed_time set to 0 (infinite)
warn contains msg if {
	some name, exporter in input.exporters
	lib.is_zero_duration(exporter.retry_on_failure.max_elapsed_time)
	msg := sprintf(
		"OTEL-053: exporter '%s' retry max_elapsed_time is 0 (infinite retries). Risk of unbounded queue growth.",
		[name],
	)
}

# ============================================================
# Private helpers
# ============================================================

_has_http_scheme(endpoint) if startswith(endpoint, "http://")
_has_http_scheme(endpoint) if startswith(endpoint, "https://")
