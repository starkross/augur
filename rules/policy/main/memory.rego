package main

import data.lib
import future.keywords.contains
import future.keywords.if
import future.keywords.in

# Memory and resource limit rules

# OTEL-024: send_batch_max_size < send_batch_size
deny contains msg if {
	some name, proc in input.processors
	split(name, "/")[0] == "batch"
	proc.send_batch_max_size < proc.send_batch_size
	msg := sprintf(
		"OTEL-024: processor '%s' send_batch_max_size (%d) < send_batch_size (%d).",
		[name, proc.send_batch_max_size, proc.send_batch_size],
	)
}

# OTEL-027: memory_limiter check_interval is 0 or unset
deny contains msg if {
	some name, proc in input.processors
	split(name, "/")[0] == "memory_limiter"
	not _has_valid_check_interval(proc)
	msg := sprintf(
		"OTEL-027: processor '%s' check_interval is 0 or unset. Memory limiter is effectively disabled.",
		[name],
	)
}

# OTEL-028: spike_limit_mib >= limit_mib
deny contains msg if {
	some name, proc in input.processors
	split(name, "/")[0] == "memory_limiter"
	proc.spike_limit_mib >= proc.limit_mib
	msg := sprintf(
		"OTEL-028: processor '%s' spike_limit_mib (%d) >= limit_mib (%d). Soft limit is zero or negative.",
		[name, proc.spike_limit_mib, proc.limit_mib],
	)
}

# OTEL-029: neither limit_mib nor limit_percentage set
deny contains msg if {
	some name, proc in input.processors
	split(name, "/")[0] == "memory_limiter"
	not _has_memory_limit(proc)
	msg := sprintf("OTEL-029: processor '%s' has neither limit_mib nor limit_percentage set.", [name])
}

# OTEL-023: batch send_batch_max_size unset (unlimited)
warn contains msg if {
	some name, proc in input.processors
	split(name, "/")[0] == "batch"
	not _has_batch_max_size(proc)
	msg := sprintf(
		"OTEL-023: processor '%s' has send_batch_max_size unset (unlimited). Set to ~2x send_batch_size.",
		[name],
	)
}

# OTEL-025: batch timeout below 100ms
warn contains msg if {
	some name, proc in input.processors
	split(name, "/")[0] == "batch"
	proc.timeout
	ms := lib.parse_duration_ms(proc.timeout)
	ms < 100
	msg := sprintf("OTEL-025: processor '%s' timeout %v is below 100ms. Defeats batching purpose.", [name, proc.timeout])
}

# OTEL-026: batch timeout above 60s
warn contains msg if {
	some name, proc in input.processors
	split(name, "/")[0] == "batch"
	proc.timeout
	ms := lib.parse_duration_ms(proc.timeout)
	ms > 60000
	msg := sprintf("OTEL-026: processor '%s' timeout %v exceeds 60s. Data sits in memory too long.", [name, proc.timeout])
}

# OTEL-030: limit_percentage outside safe range
warn contains msg if {
	some name, proc in input.processors
	split(name, "/")[0] == "memory_limiter"
	proc.limit_percentage > 90
	msg := sprintf(
		"OTEL-030: processor '%s' limit_percentage %v exceeds 90%%. Risk of OOM before limiter activates.",
		[name, proc.limit_percentage],
	)
}

warn contains msg if {
	some name, proc in input.processors
	split(name, "/")[0] == "memory_limiter"
	proc.limit_percentage < 20
	msg := sprintf(
		"OTEL-030: processor '%s' limit_percentage %v is below 20%%. Wastes majority of allocated memory.",
		[name, proc.limit_percentage],
	)
}

# ============================================================
# Private helpers
# ============================================================

_has_valid_check_interval(proc) if {
	is_object(proc)
	not lib.is_zero_duration(proc.check_interval)
}

_has_memory_limit(proc) if proc.limit_mib
_has_memory_limit(proc) if proc.limit_percentage

_has_batch_max_size(proc) if {
	is_object(proc)
	proc.send_batch_max_size > 0
}
