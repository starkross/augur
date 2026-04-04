package lib

import future.keywords.if
import future.keywords.in

pipeline_types contains t if {
	some t, _ in input.service.pipelines
}

pipeline_processors(ptype) := ps if {
	ps := input.service.pipelines[ptype].processors
} else := []

pipeline_receivers(ptype) := rs if {
	rs := input.service.pipelines[ptype].receivers
} else := []

pipeline_exporters(ptype) := es if {
	es := input.service.pipelines[ptype].exporters
} else := []

all_used_receivers contains r if {
	some t in pipeline_types
	some r in pipeline_receivers(t)
}

all_used_exporters contains e if {
	some t in pipeline_types
	some e in pipeline_exporters(t)
}

all_used_processors contains p if {
	some t in pipeline_types
	some p in pipeline_processors(t)
}

is_env_var(val) if {
	startswith(val, "${env:")
	endswith(val, "}")
}

is_env_var(val) if {
	startswith(val, "${ENV:")
	endswith(val, "}")
}

looks_like_secret(key) if {
	secret_patterns := {
		"api_key", "apikey", "token", "secret", "password", "credential",
		"private_key", "passphrase", "signing_key", "access_key",
		"auth_token", "bearer", "connection_string",
	}
	some pattern in secret_patterns
	contains(lower(key), pattern)
}

obj_contains_string(obj, substr) if {
	marshaled := json.marshal(obj)
	contains(marshaled, substr)
}

# Duration parsing: converts OTel duration strings to milliseconds.
parse_duration_ms(val) := ms if {
	is_string(val)
	endswith(val, "ms")
	num_str := substring(val, 0, count(val) - 2)
	ms := to_number(num_str)
}

parse_duration_ms(val) := ms if {
	is_string(val)
	endswith(val, "s")
	not endswith(val, "ms")
	num_str := substring(val, 0, count(val) - 1)
	ms := to_number(num_str) * 1000
}

parse_duration_ms(val) := ms if {
	is_string(val)
	endswith(val, "m")
	num_str := substring(val, 0, count(val) - 1)
	ms := to_number(num_str) * 60000
}

parse_duration_ms(val) := ms if {
	is_string(val)
	endswith(val, "h")
	num_str := substring(val, 0, count(val) - 1)
	ms := to_number(num_str) * 3600000
}

parse_duration_ms(val) := ms if {
	is_number(val)
	ms := val * 1000
}

# Check if a duration value represents zero or is unset.
is_zero_duration(val) if is_null(val)
is_zero_duration(0) := true
is_zero_duration("0") := true
is_zero_duration("0s") := true
is_zero_duration("0ms") := true
