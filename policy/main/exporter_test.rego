package main_test

import future.keywords.if
import future.keywords.in

import data.main

test_044_deny_grpc_endpoint_with_scheme if {
	val := {"endpoint": "http://backend:4317"}
	cfg := json.patch(valid_config, [{"op": "add", "path": "/exporters/otlp", "value": val}])
	msgs := main.deny with input as cfg
	some msg in msgs
	contains(msg, "OTEL-044")
}

test_045_warn_grpc_endpoint_no_port if {
	val := {"endpoint": "backend.example.com"}
	cfg := json.patch(valid_config, [{"op": "add", "path": "/exporters/otlp", "value": val}])
	msgs := main.warn with input as cfg
	some msg in msgs
	contains(msg, "OTEL-045")
}

test_046_warn_http_endpoint_no_scheme if {
	val := {"endpoint": "backend.example.com:4318"}
	cfg := json.patch(valid_config, [{"op": "add", "path": "/exporters/otlphttp", "value": val}])
	msgs := main.warn with input as cfg
	some msg in msgs
	contains(msg, "OTEL-046")
}

test_047_warn_http_using_grpc_port if {
	val := {"endpoint": "https://backend.example.com:4317"}
	cfg := json.patch(valid_config, [{"op": "add", "path": "/exporters/otlphttp", "value": val}])
	msgs := main.warn with input as cfg
	some msg in msgs
	contains(msg, "OTEL-047")
}

test_048_warn_sending_queue_disabled if {
	val := {"endpoint": "backend:4317", "sending_queue": {"enabled": false}}
	cfg := json.patch(valid_config, [{"op": "add", "path": "/exporters/otlp", "value": val}])
	msgs := main.warn with input as cfg
	some msg in msgs
	contains(msg, "OTEL-048")
}

test_049_warn_queue_size_too_small if {
	val := {"endpoint": "backend:4317", "sending_queue": {"queue_size": 1}}
	cfg := json.patch(valid_config, [{"op": "add", "path": "/exporters/otlp", "value": val}])
	msgs := main.warn with input as cfg
	some msg in msgs
	contains(msg, "OTEL-049")
}

test_050_warn_queue_size_too_large if {
	val := {"endpoint": "backend:4317", "sending_queue": {"queue_size": 100000}}
	cfg := json.patch(valid_config, [{"op": "add", "path": "/exporters/otlp", "value": val}])
	msgs := main.warn with input as cfg
	some msg in msgs
	contains(msg, "OTEL-050")
}

test_051_warn_num_consumers_low if {
	val := {"endpoint": "backend:4317", "sending_queue": {"num_consumers": 1}}
	cfg := json.patch(valid_config, [{"op": "add", "path": "/exporters/otlp", "value": val}])
	msgs := main.warn with input as cfg
	some msg in msgs
	contains(msg, "OTEL-051")
}

test_052_warn_compression_disabled if {
	val := {"endpoint": "backend:4317", "compression": "none"}
	cfg := json.patch(valid_config, [{"op": "add", "path": "/exporters/otlp", "value": val}])
	msgs := main.warn with input as cfg
	some msg in msgs
	contains(msg, "OTEL-052")
}

test_053_warn_infinite_retries if {
	val := {"endpoint": "backend:4317", "retry_on_failure": {"enabled": true, "max_elapsed_time": "0s"}}
	cfg := json.patch(valid_config, [{"op": "add", "path": "/exporters/otlp", "value": val}])
	msgs := main.warn with input as cfg
	some msg in msgs
	contains(msg, "OTEL-053")
}
