package main_test

import future.keywords.if
import future.keywords.in

import data.main

test_069_warn_metrics_level_none if {
	val := {"metrics": {"level": "none"}}
	cfg := json.patch(valid_config, [{"op": "add", "path": "/service/telemetry", "value": val}])
	msgs := main.warn with input as cfg
	some msg in msgs
	contains(msg, "OTEL-069")
}

test_070_warn_metrics_address_wildcard if {
	val := {"metrics": {"address": "0.0.0.0:8888"}}
	cfg := json.patch(valid_config, [{"op": "add", "path": "/service/telemetry", "value": val}])
	msgs := main.warn with input as cfg
	some msg in msgs
	contains(msg, "OTEL-070")
}

test_071_warn_logging_exporter if {
	cfg := json.patch(valid_config, [{"op": "add", "path": "/exporters/logging", "value": {"loglevel": "debug"}}])
	msgs := main.warn with input as cfg
	some msg in msgs
	contains(msg, "OTEL-071")
}

test_072_warn_opencensus_receiver if {
	val := {"endpoint": "0.0.0.0:55678"}
	cfg := json.patch(valid_config, [{"op": "add", "path": "/receivers/opencensus", "value": val}])
	msgs := main.warn with input as cfg
	some msg in msgs
	contains(msg, "OTEL-072")
}

test_073_warn_ballast_size_mib if {
	val := {"check_interval": "1s", "limit_mib": 512, "ballast_size_mib": 256}
	cfg := json.patch(valid_config, [{"op": "replace", "path": "/processors/memory_limiter", "value": val}])
	msgs := main.warn with input as cfg
	some msg in msgs
	contains(msg, "OTEL-073")
}

test_074_warn_deprecated_metrics_address if {
	val := {"metrics": {"address": "127.0.0.1:8888"}}
	cfg := json.patch(valid_config, [{"op": "add", "path": "/service/telemetry", "value": val}])
	msgs := main.warn with input as cfg
	some msg in msgs
	contains(msg, "OTEL-074")
}
