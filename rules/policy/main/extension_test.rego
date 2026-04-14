package main_test

import future.keywords.if
import future.keywords.in

import data.main

test_059_warn_pprof_enabled if {
	cfg := json.patch(valid_config, [{"op": "add", "path": "/extensions/pprof", "value": {"endpoint": "localhost:1777"}}])
	msgs := main.warn with input as cfg
	some msg in msgs
	contains(msg, "OTEL-059")
}

test_060_warn_zpages_wildcard_bind if {
	cfg := json.patch(valid_config, [{"op": "add", "path": "/extensions/zpages", "value": {"endpoint": "0.0.0.0:55679"}}])
	msgs := main.warn with input as cfg
	some msg in msgs
	contains(msg, "OTEL-060")
}

test_061_warn_memory_ballast if {
	cfg := json.patch(valid_config, [{"op": "add", "path": "/extensions/memory_ballast", "value": {"size_mib": 512}}])
	msgs := main.warn with input as cfg
	some msg in msgs
	contains(msg, "OTEL-061")
}

test_062_warn_undefined_extension if {
	cfg := json.patch(valid_config, [{"op": "replace", "path": "/service/extensions", "value": ["health_check", "pprof"]}])
	msgs := main.warn with input as cfg
	some msg in msgs
	contains(msg, "OTEL-062")
}

test_066_deny_undefined_storage_ref if {
	val := {"endpoint": "backend:4317", "sending_queue": {"storage": "file_storage/missing"}}
	cfg := json.patch(valid_config, [{"op": "add", "path": "/exporters/otlp", "value": val}])
	msgs := main.deny with input as cfg
	some msg in msgs
	contains(msg, "OTEL-066")
}
