package main_test

import future.keywords.if
import future.keywords.in

import data.main

test_063_warn_tail_sampling_no_groupbytrace if {
	procs := ["memory_limiter", "tail_sampling", "batch"]
	cfg := json.patch(valid_config, [
		{"op": "add", "path": "/processors/tail_sampling", "value": {}},
		{"op": "replace", "path": "/service/pipelines/traces/processors", "value": procs},
	])
	msgs := main.warn with input as cfg
	some msg in msgs
	contains(msg, "OTEL-063")
}

test_064_warn_both_samplers if {
	procs := ["memory_limiter", "probabilistic_sampler", "tail_sampling", "batch"]
	cfg := json.patch(valid_config, [
		{"op": "add", "path": "/processors/probabilistic_sampler", "value": {}},
		{"op": "add", "path": "/processors/tail_sampling", "value": {}},
		{"op": "replace", "path": "/service/pipelines/traces/processors", "value": procs},
	])
	msgs := main.warn with input as cfg
	some msg in msgs
	contains(msg, "OTEL-064")
}

test_065_warn_queue_no_storage if {
	msgs := main.warn with input as valid_config
	some msg in msgs
	contains(msg, "OTEL-065")
}

test_067_warn_k8s_no_k8sattributes if {
	val := {"auth_type": "serviceAccount"}
	cfg := json.patch(valid_config, [{"op": "add", "path": "/receivers/kubeletstats", "value": val}])
	msgs := main.warn with input as cfg
	some msg in msgs
	contains(msg, "OTEL-067")
}

test_068_warn_k8s_no_resourcedetection if {
	val := {"auth_type": "serviceAccount"}
	cfg := json.patch(valid_config, [{"op": "add", "path": "/receivers/kubeletstats", "value": val}])
	msgs := main.warn with input as cfg
	some msg in msgs
	contains(msg, "OTEL-068")
}
