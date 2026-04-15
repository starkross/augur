package main_test

import future.keywords.if
import future.keywords.in

import data.main

test_031_deny_tls_min_version_low if {
	val := {"endpoint": "backend:4317", "tls": {"min_version": "1.0"}}
	cfg := json.patch(valid_config, [{"op": "add", "path": "/exporters/otlp_test", "value": val}])
	msgs := main.deny with input as cfg
	some msg in msgs
	contains(msg, "OTEL-031")
}

test_031_deny_tls_min_version_receiver if {
	rcv := {"protocols": {"grpc": {"endpoint": "localhost:4317", "tls": {"min_version": "1.1"}}}}
	cfg := json.patch(valid_config, [{"op": "replace", "path": "/receivers/otlp", "value": rcv}])
	msgs := main.deny with input as cfg
	some msg in msgs
	contains(msg, "OTEL-031")
}

test_032_warn_insecure_skip_verify if {
	val := {"endpoint": "backend:4317", "tls": {"insecure_skip_verify": true}}
	cfg := json.patch(valid_config, [{"op": "add", "path": "/exporters/otlp_test", "value": val}])
	msgs := main.warn with input as cfg
	some msg in msgs
	contains(msg, "OTEL-032")
}

test_033_warn_receiver_no_tls if {
	val := {"protocols": {"grpc": {"endpoint": "10.0.0.5:4317"}}}
	cfg := json.patch(valid_config, [{"op": "replace", "path": "/receivers/otlp", "value": val}])
	msgs := main.warn with input as cfg
	some msg in msgs
	contains(msg, "OTEL-033")
}

test_033_pass_on_localhost if {
	msgs := main.warn with input as valid_config
	not_contains_rule(msgs, "OTEL-033")
}

test_033_pass_on_unix_transport if {
	val := {"protocols": {"grpc": {"endpoint": "/run/xxxx-metrics.sock", "transport": "unix"}}}
	cfg := json.patch(valid_config, [{"op": "replace", "path": "/receivers/otlp", "value": val}])
	msgs := main.warn with input as cfg
	not_contains_rule(msgs, "OTEL-033")
}

test_034_deny_cors_wildcard if {
	val := {"protocols": {"http": {"cors": {"allowed_origins": ["*"]}, "endpoint": "localhost:4318"}}}
	cfg := json.patch(valid_config, [{"op": "replace", "path": "/receivers/otlp", "value": val}])
	msgs := main.deny with input as cfg
	some msg in msgs
	contains(msg, "OTEL-034")
}

test_035_deny_hardcoded_extension_secret if {
	val := {"token": "my-secret-token"}
	cfg := json.patch(valid_config, [{"op": "add", "path": "/extensions/bearertokenauth", "value": val}])
	msgs := main.deny with input as cfg
	some msg in msgs
	contains(msg, "OTEL-035")
}

test_035_deny_nested_extension_secret if {
	val := {"client_auth": {"username": "admin", "password": "plaintext"}}
	cfg := json.patch(valid_config, [{"op": "add", "path": "/extensions/basicauth", "value": val}])
	msgs := main.deny with input as cfg
	some msg in msgs
	contains(msg, "OTEL-035")
}

test_036_warn_large_recv_msg_size if {
	val := {"protocols": {"grpc": {"endpoint": "localhost:4317", "max_recv_msg_size_mib": 1024}}}
	cfg := json.patch(valid_config, [{"op": "replace", "path": "/receivers/otlp", "value": val}])
	msgs := main.warn with input as cfg
	some msg in msgs
	contains(msg, "OTEL-036")
}

test_037_warn_inline_key_pem if {
	val := {"endpoint": "backend:4317", "tls": {"key_pem": "-----BEGIN PRIVATE KEY-----"}}
	cfg := json.patch(valid_config, [{"op": "add", "path": "/exporters/otlp_test", "value": val}])
	msgs := main.warn with input as cfg
	some msg in msgs
	contains(msg, "OTEL-037")
}
