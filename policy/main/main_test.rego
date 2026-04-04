package main

import future.keywords.if
import future.keywords.in

valid_config := {
	"receivers": {"otlp": {"protocols": {"grpc": {"endpoint": "localhost:4317"}}}},
	"processors": {
		"batch": {},
		"memory_limiter": {"check_interval": "5s", "limit_mib": 4000},
	},
	"exporters": {"otlp/backend": {
		"endpoint": "${env:OTEL_EXPORTER_ENDPOINT}",
		"retry_on_failure": {"enabled": true},
		"sending_queue": {"enabled": true},
	}},
	"extensions": {"health_check": {}},
	"service": {
		"extensions": ["health_check"],
		"pipelines": {"traces": {
			"receivers": ["otlp"],
			"processors": ["memory_limiter", "batch"],
			"exporters": ["otlp/backend"],
		}},
	},
}

test_001_deny_when_no_memory_limiter if {
	cfg := json.patch(valid_config, [{"op": "remove", "path": "/processors/memory_limiter"}])
	msgs := deny with input as cfg
	some msg in msgs
	contains(msg, "OTEL-001")
}

test_001_pass_when_memory_limiter_exists if {
	msgs := deny with input as valid_config
	not_contains_rule(msgs, "OTEL-001")
}

test_002_deny_when_pipeline_missing_memory_limiter if {
	cfg := json.patch(valid_config, [
		{"op": "replace", "path": "/service/pipelines/traces/processors", "value": ["batch"]},
	])
	msgs := deny with input as cfg
	some msg in msgs
	contains(msg, "OTEL-002")
}

test_003_deny_when_no_batch if {
	cfg := json.patch(valid_config, [{"op": "remove", "path": "/processors/batch"}])
	msgs := deny with input as cfg
	some msg in msgs
	contains(msg, "OTEL-003")
}

test_004_deny_hardcoded_api_key if {
	cfg := json.patch(valid_config, [
		{"op": "add", "path": "/exporters/someexporter", "value": {"api_key": "sk-12345", "endpoint": "https://example.com"}},
		{"op": "add", "path": "/service/pipelines/traces/exporters/-", "value": "someexporter"},
	])
	msgs := deny with input as cfg
	some msg in msgs
	contains(msg, "OTEL-004")
}

test_004_pass_env_var_secret if {
	cfg := json.patch(valid_config, [
		{"op": "add", "path": "/exporters/someexporter", "value": {"api_key": "${env:MY_KEY}"}},
	])
	msgs := deny with input as cfg
	not_contains_rule(msgs, "OTEL-004")
}

test_010_warn_on_wildcard_bind if {
	cfg := json.patch(valid_config, [
		{"op": "replace", "path": "/receivers/otlp/protocols/grpc/endpoint", "value": "0.0.0.0:4317"},
	])
	msgs := warn with input as cfg
	some msg in msgs
	contains(msg, "OTEL-010")
}

test_010_pass_on_localhost_bind if {
	msgs := warn with input as valid_config
	not_contains_rule(msgs, "OTEL-010")
}

test_013_warn_batch_not_last if {
	cfg := json.patch(valid_config, [
		{"op": "replace", "path": "/service/pipelines/traces/processors", "value": ["batch", "memory_limiter"]},
	])
	msgs := warn with input as cfg
	some msg in msgs
	contains(msg, "OTEL-013")
}

test_014_warn_memory_limiter_not_first if {
	cfg := json.patch(valid_config, [
		{"op": "replace", "path": "/service/pipelines/traces/processors", "value": ["batch", "memory_limiter"]},
	])
	msgs := warn with input as cfg
	some msg in msgs
	contains(msg, "OTEL-014")
}

test_015_warn_debug_exporter if {
	cfg := json.patch(valid_config, [
		{"op": "add", "path": "/exporters/debug", "value": {}},
	])
	msgs := warn with input as cfg
	some msg in msgs
	contains(msg, "OTEL-015")
}

test_020_warn_unused_receiver if {
	cfg := json.patch(valid_config, [
		{"op": "add", "path": "/receivers/jaeger", "value": {"protocols": {"grpc": {}}}},
	])
	msgs := warn with input as cfg
	some msg in msgs
	contains(msg, "OTEL-020")
}

test_valid_config_no_denials if {
	msgs := deny with input as valid_config
	count(msgs) == 0
}

not_contains_rule(msgs, rule_id) if {
	every msg in msgs {
		not contains(msg, rule_id)
	}
}
