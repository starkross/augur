package main

import data.lib
import future.keywords.contains
import future.keywords.if
import future.keywords.in

# Lifecycle and observability rules

# OTEL-069: telemetry metrics level none
warn contains msg if {
	input.service.telemetry.metrics.level == "none"
	msg := "OTEL-069: service telemetry metrics level is 'none'. No visibility into Collector health."
}

# OTEL-070: telemetry metrics address on 0.0.0.0
warn contains msg if {
	addr := input.service.telemetry.metrics.address
	contains(addr, "0.0.0.0")
	msg := "OTEL-070: telemetry metrics address bound to 0.0.0.0. Use localhost for security."
}

# OTEL-071: logging exporter (renamed to debug)
warn contains msg if {
	some name, _ in input.exporters
	split(name, "/")[0] == "logging"
	msg := sprintf("OTEL-071: exporter '%s' uses deprecated 'logging' exporter. Renamed to 'debug' in v0.111.0.", [name])
}

# OTEL-072: opencensus receiver (deprecated)
warn contains msg if {
	some name, _ in input.receivers
	split(name, "/")[0] == "opencensus"
	msg := sprintf("OTEL-072: receiver '%s' uses deprecated OpenCensus protocol. Migrate to OTLP.", [name])
}

# OTEL-072: opencensus exporter (deprecated)
warn contains msg if {
	some name, _ in input.exporters
	split(name, "/")[0] == "opencensus"
	msg := sprintf("OTEL-072: exporter '%s' uses deprecated OpenCensus protocol. Migrate to OTLP.", [name])
}

# OTEL-073: memory_limiter ballast_size_mib (deprecated)
warn contains msg if {
	some name, proc in input.processors
	split(name, "/")[0] == "memory_limiter"
	proc.ballast_size_mib
	msg := sprintf("OTEL-073: processor '%s' uses deprecated ballast_size_mib. Use GOMEMLIMIT env var.", [name])
}

# OTEL-074: service.telemetry.metrics.address (deprecated)
warn contains msg if {
	input.service.telemetry.metrics.address
	msg := "OTEL-074: service.telemetry.metrics.address is deprecated. Use 'readers' configuration instead."
}
