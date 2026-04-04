package main

import data.lib
import future.keywords.contains
import future.keywords.if
import future.keywords.in

# Receiver configuration rules

# OTEL-058: multiple receivers bound to the same endpoint
deny contains msg if {
	some ep1 in _receiver_endpoints
	some ep2 in _receiver_endpoints
	ep1.endpoint == ep2.endpoint
	ep1.receiver != ep2.receiver
	ep1.receiver < ep2.receiver
	msg := sprintf("OTEL-058: receivers '%s' and '%s' both bind to '%s'.", [ep1.receiver, ep2.receiver, ep1.endpoint])
}

# OTEL-054: prometheus scrape_interval below 10s
warn contains msg if {
	some name, receiver in input.receivers
	split(name, "/")[0] == "prometheus"
	some _, sc in receiver.config.scrape_configs
	ms := lib.parse_duration_ms(sc.scrape_interval)
	ms < 10000
	job := object.get(sc, "job_name", "unknown")
	msg := sprintf(
		"OTEL-054: receiver '%s' job '%s' scrape_interval %v is below 10s.",
		[name, job, sc.scrape_interval],
	)
}

# OTEL-055: hostmetrics collection_interval below 10s
warn contains msg if {
	some name, receiver in input.receivers
	split(name, "/")[0] == "hostmetrics"
	receiver.collection_interval
	ms := lib.parse_duration_ms(receiver.collection_interval)
	ms < 10000
	msg := sprintf(
		"OTEL-055: receiver '%s' collection_interval %v is below 10s. High CPU overhead.",
		[name, receiver.collection_interval],
	)
}

# OTEL-056: filelog start_at:beginning without storage
warn contains msg if {
	some name, receiver in input.receivers
	split(name, "/")[0] == "filelog"
	receiver.start_at == "beginning"
	not receiver.storage
	msg := sprintf(
		"OTEL-056: receiver '%s' has start_at:beginning without storage. Causes duplicate ingestion on restart.",
		[name],
	)
}

# OTEL-057: filelog overly broad include pattern
warn contains msg if {
	some name, receiver in input.receivers
	split(name, "/")[0] == "filelog"
	some pattern in receiver.include
	_is_broad_glob(pattern)
	msg := sprintf(
		"OTEL-057: receiver '%s' has overly broad include pattern '%s'. Scope to specific files.",
		[name, pattern],
	)
}

# ============================================================
# Private helpers
# ============================================================

_is_broad_glob(pattern) if endswith(pattern, "/*")
_is_broad_glob(pattern) if endswith(pattern, "/**")
_is_broad_glob(pattern) if endswith(pattern, "/**/*")

# Receiver endpoint collection for conflict detection (OTEL-058)
_receiver_endpoints contains {"receiver": name, "endpoint": endpoint} if {
	some name, receiver in input.receivers
	some _, proto_cfg in receiver.protocols
	endpoint := proto_cfg.endpoint
	is_string(endpoint)
}

_receiver_endpoints contains {"receiver": name, "endpoint": endpoint} if {
	some name, receiver in input.receivers
	not receiver.protocols
	endpoint := receiver.endpoint
	is_string(endpoint)
}
