package main_test

import future.keywords.if
import future.keywords.in

import data.main

test_038_warn_filter_after_batch if {
	cfg := json.patch(valid_config, [
		{"op": "add", "path": "/processors/filter", "value": {}},
		{"op": "replace", "path": "/service/pipelines/traces/processors", "value": ["memory_limiter", "batch", "filter"]},
	])
	msgs := main.warn with input as cfg
	some msg in msgs
	contains(msg, "OTEL-038")
}

test_039_warn_transform_after_batch if {
	cfg := json.patch(valid_config, [
		{"op": "add", "path": "/processors/transform", "value": {}},
		{"op": "replace", "path": "/service/pipelines/traces/processors", "value": ["memory_limiter", "batch", "transform"]},
	])
	msgs := main.warn with input as cfg
	some msg in msgs
	contains(msg, "OTEL-039")
}

test_040_deny_circular_dependency if {
	cfg := {
		"receivers": {"otlp": {"protocols": {"grpc": {"endpoint": "localhost:4317"}}}},
		"processors": {
			"memory_limiter": {"check_interval": "5s", "limit_mib": 4000},
			"batch": {"send_batch_max_size": 16384},
		},
		"exporters": {},
		"connectors": {"forward/a": null, "forward/b": null},
		"extensions": {"health_check": {}},
		"service": {
			"extensions": ["health_check"],
			"pipelines": {
				"traces/a": {
					"receivers": ["forward/b"],
					"processors": ["memory_limiter", "batch"],
					"exporters": ["forward/a"],
				},
				"traces/b": {
					"receivers": ["forward/a"],
					"processors": ["memory_limiter", "batch"],
					"exporters": ["forward/b"],
				},
			},
		},
	}
	msgs := main.deny with input as cfg
	some msg in msgs
	contains(msg, "OTEL-040")
}

test_041_warn_routing_no_default if {
	val := {"routing": {"table": [{"statement": "route()", "pipelines": ["traces"]}]}}
	cfg := json.patch(valid_config, [{"op": "add", "path": "/connectors", "value": val}])
	msgs := main.warn with input as cfg
	some msg in msgs
	contains(msg, "OTEL-041")
}

test_042_warn_duplicate_processor if {
	val := ["memory_limiter", "batch", "batch"]
	cfg := json.patch(valid_config, [{"op": "replace", "path": "/service/pipelines/traces/processors", "value": val}])
	msgs := main.warn with input as cfg
	some msg in msgs
	contains(msg, "OTEL-042")
}

test_043_warn_batch_before_tail_sampling if {
	procs := ["memory_limiter", "batch", "tail_sampling"]
	cfg := json.patch(valid_config, [
		{"op": "add", "path": "/processors/tail_sampling", "value": {}},
		{"op": "replace", "path": "/service/pipelines/traces/processors", "value": procs},
	])
	msgs := main.warn with input as cfg
	some msg in msgs
	contains(msg, "OTEL-043")
}
