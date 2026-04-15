package main_test

import future.keywords.if
import future.keywords.in

import data.main

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
	msgs := main.deny with input as cfg
	some msg in msgs
	contains(msg, "OTEL-001")
}

test_001_pass_when_memory_limiter_exists if {
	msgs := main.deny with input as valid_config
	not_contains_rule(msgs, "OTEL-001")
}

test_002_deny_when_pipeline_missing_memory_limiter if {
	val := ["batch"]
	cfg := json.patch(valid_config, [{"op": "replace", "path": "/service/pipelines/traces/processors", "value": val}])
	msgs := main.deny with input as cfg
	some msg in msgs
	contains(msg, "OTEL-002")
}

test_003_deny_when_no_batch if {
	cfg := json.patch(valid_config, [{"op": "remove", "path": "/processors/batch"}])
	msgs := main.deny with input as cfg
	some msg in msgs
	contains(msg, "OTEL-003")
}

test_004_deny_hardcoded_api_key if {
	cfg := json.patch(valid_config, [
		{"op": "add", "path": "/exporters/someexporter", "value": {"api_key": "sk-12345", "endpoint": "https://example.com"}},
		{"op": "add", "path": "/service/pipelines/traces/exporters/-", "value": "someexporter"},
	])
	msgs := main.deny with input as cfg
	some msg in msgs
	contains(msg, "OTEL-004")
}

test_004_pass_env_var_secret if {
	val := {"api_key": "${env:MY_KEY}"}
	cfg := json.patch(valid_config, [{"op": "add", "path": "/exporters/someexporter", "value": val}])
	msgs := main.deny with input as cfg
	not_contains_rule(msgs, "OTEL-004")
}

test_010_warn_on_wildcard_bind if {
	p := "/receivers/otlp/protocols/grpc/endpoint"
	cfg := json.patch(valid_config, [{"op": "replace", "path": p, "value": "0.0.0.0:4317"}])
	msgs := main.warn with input as cfg
	some msg in msgs
	contains(msg, "OTEL-010")
}

test_010_pass_on_localhost_bind if {
	msgs := main.warn with input as valid_config
	not_contains_rule(msgs, "OTEL-010")
}

test_018_pass_on_unix_scheme if {
	val := {"endpoint": "unix:///var/run/otel.sock"}
	cfg := json.patch(valid_config, [{"op": "replace", "path": "/exporters/otlp~1backend", "value": val}])
	msgs := main.warn with input as cfg
	not_contains_rule(msgs, "OTEL-018")
}

test_013_warn_batch_not_last if {
	val := ["batch", "memory_limiter"]
	cfg := json.patch(valid_config, [{"op": "replace", "path": "/service/pipelines/traces/processors", "value": val}])
	msgs := main.warn with input as cfg
	some msg in msgs
	contains(msg, "OTEL-013")
}

test_014_warn_memory_limiter_not_first if {
	val := ["batch", "memory_limiter"]
	cfg := json.patch(valid_config, [{"op": "replace", "path": "/service/pipelines/traces/processors", "value": val}])
	msgs := main.warn with input as cfg
	some msg in msgs
	contains(msg, "OTEL-014")
}

test_015_warn_debug_exporter if {
	cfg := json.patch(valid_config, [{"op": "add", "path": "/exporters/debug", "value": {}}])
	msgs := main.warn with input as cfg
	some msg in msgs
	contains(msg, "OTEL-015")
}

test_017_warn_otlp_no_retry if {
	cfg := json.patch(valid_config, [
		{"op": "add", "path": "/exporters/otlp~1noretry", "value": {"endpoint": "backend:4317"}},
		{"op": "add", "path": "/service/pipelines/traces/exporters/-", "value": "otlp/noretry"},
	])
	msgs := main.warn with input as cfg
	some msg in msgs
	contains(msg, "OTEL-017")
	contains(msg, "otlp/noretry")
}

test_017_pass_awsemf_with_max_retries if {
	cfg := json.patch(valid_config, [
		{"op": "add", "path": "/exporters/awsemf~1foo", "value": {"region": "us-east-1", "max_retries": 5}},
		{"op": "add", "path": "/service/pipelines/traces/exporters/-", "value": "awsemf/foo"},
	])
	msgs := main.warn with input as cfg
	not_contains_rule(msgs, "OTEL-017")
}

test_017_warn_awsemf_without_max_retries if {
	cfg := json.patch(valid_config, [
		{"op": "add", "path": "/exporters/awsemf~1bar", "value": {"region": "us-east-1"}},
		{"op": "add", "path": "/service/pipelines/traces/exporters/-", "value": "awsemf/bar"},
	])
	msgs := main.warn with input as cfg
	some msg in msgs
	contains(msg, "OTEL-017")
	contains(msg, "awsemf/bar")
}

test_017_pass_awsxray_with_max_retries if {
	cfg := json.patch(valid_config, [
		{"op": "add", "path": "/exporters/awsxray", "value": {"region": "us-east-1", "max_retries": 3}},
		{"op": "add", "path": "/service/pipelines/traces/exporters/-", "value": "awsxray"},
	])
	msgs := main.warn with input as cfg
	not_contains_rule(msgs, "OTEL-017")
}

test_017_warn_awsxray_without_max_retries if {
	cfg := json.patch(valid_config, [
		{"op": "add", "path": "/exporters/awsxray", "value": {"region": "us-east-1"}},
		{"op": "add", "path": "/service/pipelines/traces/exporters/-", "value": "awsxray"},
	])
	msgs := main.warn with input as cfg
	some msg in msgs
	contains(msg, "OTEL-017")
	contains(msg, "awsxray")
}

test_017_pass_awscloudwatchlogs_with_max_retries if {
	cfg := json.patch(valid_config, [
		{"op": "add", "path": "/exporters/awscloudwatchlogs", "value": {"region": "us-east-1", "max_retries": 3}},
		{"op": "add", "path": "/service/pipelines/traces/exporters/-", "value": "awscloudwatchlogs"},
	])
	msgs := main.warn with input as cfg
	not_contains_rule(msgs, "OTEL-017")
}

test_017_pass_awss3_with_retry_max_attempts if {
	cfg := json.patch(valid_config, [
		{
			"op": "add", "path": "/exporters/awss3",
			"value": {"s3uploader": {"region": "us-east-1", "s3_bucket": "b", "retry_max_attempts": 5}},
		},
		{"op": "add", "path": "/service/pipelines/traces/exporters/-", "value": "awss3"},
	])
	msgs := main.warn with input as cfg
	not_contains_rule(msgs, "OTEL-017")
}

test_017_pass_awss3_with_retry_mode_standard if {
	cfg := json.patch(valid_config, [
		{
			"op": "add", "path": "/exporters/awss3",
			"value": {"s3uploader": {"region": "us-east-1", "s3_bucket": "b", "retry_mode": "standard"}},
		},
		{"op": "add", "path": "/service/pipelines/traces/exporters/-", "value": "awss3"},
	])
	msgs := main.warn with input as cfg
	not_contains_rule(msgs, "OTEL-017")
}

test_017_warn_awss3_with_retry_mode_nop if {
	cfg := json.patch(valid_config, [
		{
			"op": "add", "path": "/exporters/awss3",
			"value": {"s3uploader": {"region": "us-east-1", "s3_bucket": "b", "retry_mode": "nop"}},
		},
		{"op": "add", "path": "/service/pipelines/traces/exporters/-", "value": "awss3"},
	])
	msgs := main.warn with input as cfg
	some msg in msgs
	contains(msg, "OTEL-017")
	contains(msg, "awss3")
}

test_017_warn_awss3_without_retry if {
	cfg := json.patch(valid_config, [
		{
			"op": "add", "path": "/exporters/awss3",
			"value": {"s3uploader": {"region": "us-east-1", "s3_bucket": "b"}},
		},
		{"op": "add", "path": "/service/pipelines/traces/exporters/-", "value": "awss3"},
	])
	msgs := main.warn with input as cfg
	some msg in msgs
	contains(msg, "OTEL-017")
	contains(msg, "awss3")
}

test_020_warn_unused_receiver if {
	cfg := json.patch(valid_config, [{"op": "add", "path": "/receivers/jaeger", "value": {"protocols": {"grpc": {}}}}])
	msgs := main.warn with input as cfg
	some msg in msgs
	contains(msg, "OTEL-020")
}

test_valid_config_no_denials if {
	msgs := main.deny with input as valid_config
	count(msgs) == 0
}

test_valid_config_still_no_denials if {
	msgs := main.deny with input as valid_config
	count(msgs) == 0
}

not_contains_rule(msgs, rule_id) if {
	every msg in msgs {
		not contains(msg, rule_id)
	}
}
