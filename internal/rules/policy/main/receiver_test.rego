package main_test

import future.keywords.if
import future.keywords.in

import data.main

test_054_warn_fast_scrape_interval if {
	val := {"config": {"scrape_configs": [{"job_name": "test", "scrape_interval": "2s"}]}}
	cfg := json.patch(valid_config, [{"op": "add", "path": "/receivers/prometheus", "value": val}])
	msgs := main.warn with input as cfg
	some msg in msgs
	contains(msg, "OTEL-054")
}

test_055_warn_fast_collection_interval if {
	val := {"collection_interval": "1s", "scrapers": {"cpu": {}}}
	cfg := json.patch(valid_config, [{"op": "add", "path": "/receivers/hostmetrics", "value": val}])
	msgs := main.warn with input as cfg
	some msg in msgs
	contains(msg, "OTEL-055")
}

test_056_warn_filelog_start_beginning_no_storage if {
	val := {"include": ["/var/log/app/*.log"], "start_at": "beginning"}
	cfg := json.patch(valid_config, [{"op": "add", "path": "/receivers/filelog", "value": val}])
	msgs := main.warn with input as cfg
	some msg in msgs
	contains(msg, "OTEL-056")
}

test_057_warn_filelog_broad_include if {
	cfg := json.patch(valid_config, [{"op": "add", "path": "/receivers/filelog", "value": {"include": ["/var/log/*"]}}])
	msgs := main.warn with input as cfg
	some msg in msgs
	contains(msg, "OTEL-057")
}

test_058_deny_duplicate_endpoint if {
	cfg := {
		"receivers": {
			"otlp": {"protocols": {"grpc": {"endpoint": "0.0.0.0:4317"}}},
			"jaeger": {"protocols": {"grpc": {"endpoint": "0.0.0.0:4317"}}},
		},
		"processors": {
			"memory_limiter": {"check_interval": "5s", "limit_mib": 4000},
			"batch": {"send_batch_max_size": 16384},
		},
		"exporters": {"debug": {}},
		"extensions": {"health_check": {}},
		"service": {
			"extensions": ["health_check"],
			"pipelines": {"traces": {
				"receivers": ["otlp", "jaeger"],
				"processors": ["memory_limiter", "batch"],
				"exporters": ["debug"],
			}},
		},
	}
	msgs := main.deny with input as cfg
	some msg in msgs
	contains(msg, "OTEL-058")
}
