package main_test

import future.keywords.if
import future.keywords.in

import data.main

test_023_warn_batch_max_size_unset if {
	msgs := main.warn with input as valid_config
	some msg in msgs
	contains(msg, "OTEL-023")
}

test_024_deny_batch_max_less_than_size if {
	val := {"send_batch_size": 10000, "send_batch_max_size": 5000}
	cfg := json.patch(valid_config, [{"op": "replace", "path": "/processors/batch", "value": val}])
	msgs := main.deny with input as cfg
	some msg in msgs
	contains(msg, "OTEL-024")
}

test_024_pass_when_max_gte_size if {
	val := {"send_batch_size": 8192, "send_batch_max_size": 16384}
	cfg := json.patch(valid_config, [{"op": "replace", "path": "/processors/batch", "value": val}])
	msgs := main.deny with input as cfg
	not_contains_rule(msgs, "OTEL-024")
}

test_025_warn_batch_timeout_too_low if {
	cfg := json.patch(valid_config, [{"op": "replace", "path": "/processors/batch", "value": {"timeout": "10ms"}}])
	msgs := main.warn with input as cfg
	some msg in msgs
	contains(msg, "OTEL-025")
}

test_026_warn_batch_timeout_too_high if {
	cfg := json.patch(valid_config, [{"op": "replace", "path": "/processors/batch", "value": {"timeout": "120s"}}])
	msgs := main.warn with input as cfg
	some msg in msgs
	contains(msg, "OTEL-026")
}

test_027_deny_check_interval_missing if {
	cfg := json.patch(valid_config, [{"op": "replace", "path": "/processors/memory_limiter", "value": {"limit_mib": 512}}])
	msgs := main.deny with input as cfg
	some msg in msgs
	contains(msg, "OTEL-027")
}

test_027_deny_check_interval_zero if {
	val := {"check_interval": "0s", "limit_mib": 512}
	cfg := json.patch(valid_config, [{"op": "replace", "path": "/processors/memory_limiter", "value": val}])
	msgs := main.deny with input as cfg
	some msg in msgs
	contains(msg, "OTEL-027")
}

test_027_pass_valid_check_interval if {
	msgs := main.deny with input as valid_config
	not_contains_rule(msgs, "OTEL-027")
}

test_028_deny_spike_exceeds_limit if {
	val := {"check_interval": "1s", "limit_mib": 512, "spike_limit_mib": 512}
	cfg := json.patch(valid_config, [{"op": "replace", "path": "/processors/memory_limiter", "value": val}])
	msgs := main.deny with input as cfg
	some msg in msgs
	contains(msg, "OTEL-028")
}

test_029_deny_no_limit if {
	val := {"check_interval": "1s"}
	cfg := json.patch(valid_config, [{"op": "replace", "path": "/processors/memory_limiter", "value": val}])
	msgs := main.deny with input as cfg
	some msg in msgs
	contains(msg, "OTEL-029")
}

test_029_pass_with_limit_mib if {
	msgs := main.deny with input as valid_config
	not_contains_rule(msgs, "OTEL-029")
}

test_030_warn_limit_percentage_too_high if {
	val := {"check_interval": "1s", "limit_percentage": 95}
	cfg := json.patch(valid_config, [{"op": "replace", "path": "/processors/memory_limiter", "value": val}])
	msgs := main.warn with input as cfg
	some msg in msgs
	contains(msg, "OTEL-030")
}

test_030_warn_limit_percentage_too_low if {
	val := {"check_interval": "1s", "limit_percentage": 10}
	cfg := json.patch(valid_config, [{"op": "replace", "path": "/processors/memory_limiter", "value": val}])
	msgs := main.warn with input as cfg
	some msg in msgs
	contains(msg, "OTEL-030")
}
