package main

import data.lib
import future.keywords.contains
import future.keywords.if
import future.keywords.in

# Security rules

# OTEL-031: TLS min_version below 1.2 (exporters)
deny contains msg if {
	some name, exporter in input.exporters
	_tls_version_below_1_2(exporter.tls.min_version)
	msg := sprintf("OTEL-031: exporter '%s' TLS min_version '%s' is below 1.2.", [name, exporter.tls.min_version])
}

# OTEL-031: TLS min_version below 1.2 (receivers)
deny contains msg if {
	some name, receiver in input.receivers
	some proto, proto_cfg in receiver.protocols
	_tls_version_below_1_2(proto_cfg.tls.min_version)
	msg := sprintf(
		"OTEL-031: receiver '%s/%s' TLS min_version '%s' is below 1.2.",
		[name, proto, proto_cfg.tls.min_version],
	)
}

# OTEL-034: CORS allowed_origins wildcard
deny contains msg if {
	some name, receiver in input.receivers
	some proto, proto_cfg in receiver.protocols
	some origin in proto_cfg.cors.allowed_origins
	origin == "*"
	msg := sprintf("OTEL-034: receiver '%s/%s' CORS allowed_origins contains wildcard '*'.", [name, proto])
}

# OTEL-035: hardcoded secrets in extensions (top-level fields)
deny contains msg if {
	some name, ext in input.extensions
	is_object(ext)
	some k, v in ext
	is_string(v)
	lib.looks_like_secret(k)
	not lib.is_env_var(v)
	msg := sprintf("OTEL-035: extension '%s' has hardcoded '%s'. Use ${env:VAR_NAME} instead.", [name, k])
}

# OTEL-035: hardcoded secrets in extensions (nested one level, e.g. client_auth)
deny contains msg if {
	some name, ext in input.extensions
	is_object(ext)
	some _, sub_obj in ext
	is_object(sub_obj)
	some k, v in sub_obj
	is_string(v)
	lib.looks_like_secret(k)
	not lib.is_env_var(v)
	msg := sprintf("OTEL-035: extension '%s' has hardcoded '%s'. Use ${env:VAR_NAME} instead.", [name, k])
}

# OTEL-032: insecure_skip_verify in exporters
warn contains msg if {
	some name, exporter in input.exporters
	exporter.tls.insecure_skip_verify == true
	msg := sprintf("OTEL-032: exporter '%s' has insecure_skip_verify enabled. TLS verification is bypassed.", [name])
}

# OTEL-032: insecure_skip_verify in receivers
warn contains msg if {
	some name, receiver in input.receivers
	some proto, proto_cfg in receiver.protocols
	proto_cfg.tls.insecure_skip_verify == true
	msg := sprintf("OTEL-032: receiver '%s/%s' has insecure_skip_verify enabled.", [name, proto])
}

# OTEL-033: receiver on non-localhost without TLS
warn contains msg if {
	some name, receiver in input.receivers
	some proto, proto_cfg in receiver.protocols
	endpoint := object.get(proto_cfg, "endpoint", "")
	endpoint != ""
	not contains(endpoint, "localhost")
	not contains(endpoint, "127.0.0.1")
	object.get(proto_cfg, "transport", "") != "unix"
	not proto_cfg.tls
	msg := sprintf("OTEL-033: receiver '%s/%s' on non-localhost endpoint '%s' without TLS.", [name, proto, endpoint])
}

# OTEL-036: gRPC max_recv_msg_size_mib > 128
warn contains msg if {
	some name, receiver in input.receivers
	some proto, proto_cfg in receiver.protocols
	proto_cfg.max_recv_msg_size_mib > 128
	msg := sprintf(
		"OTEL-036: receiver '%s/%s' max_recv_msg_size_mib is %d (>128). Risk of decompression bomb.",
		[name, proto, proto_cfg.max_recv_msg_size_mib],
	)
}

# OTEL-037: inline key_pem in exporters
warn contains msg if {
	some name, exporter in input.exporters
	exporter.tls.key_pem
	msg := sprintf("OTEL-037: exporter '%s' has inline key_pem. Use key_file instead.", [name])
}

# OTEL-037: inline key_pem in receivers
warn contains msg if {
	some name, receiver in input.receivers
	some proto, proto_cfg in receiver.protocols
	proto_cfg.tls.key_pem
	msg := sprintf("OTEL-037: receiver '%s/%s' has inline key_pem. Use key_file instead.", [name, proto])
}

# ============================================================
# Private helpers
# ============================================================

_tls_version_below_1_2("1.0") := true
_tls_version_below_1_2("1.1") := true
