package main

import data.lib
import future.keywords.contains
import future.keywords.if
import future.keywords.in

# Extension rules

# OTEL-066: sending_queue.storage references undefined extension
deny contains msg if {
	some name, exporter in input.exporters
	storage_ref := exporter.sending_queue.storage
	is_string(storage_ref)
	not input.extensions[storage_ref]
	msg := sprintf(
		"OTEL-066: exporter '%s' sending_queue.storage references undefined extension '%s'.",
		[name, storage_ref],
	)
}

# OTEL-059: pprof extension enabled
warn contains msg if {
	some name, _ in input.extensions
	split(name, "/")[0] == "pprof"
	msg := sprintf("OTEL-059: pprof extension '%s' is enabled. Exposes profiling endpoints in production.", [name])
}

# OTEL-060: zpages endpoint on 0.0.0.0
warn contains msg if {
	some name, ext in input.extensions
	split(name, "/")[0] == "zpages"
	is_object(ext)
	endpoint := object.get(ext, "endpoint", "")
	contains(endpoint, "0.0.0.0")
	msg := sprintf("OTEL-060: zpages extension '%s' endpoint bound to 0.0.0.0. Use localhost.", [name])
}

# OTEL-061: memory_ballast extension (deprecated)
warn contains msg if {
	some name, _ in input.extensions
	split(name, "/")[0] == "memory_ballast"
	msg := sprintf("OTEL-061: extension '%s' is deprecated. Use GOMEMLIMIT env var instead.", [name])
}

# OTEL-062: extension in service.extensions but not defined
warn contains msg if {
	svc_extensions := object.get(input.service, "extensions", [])
	some ext_name in svc_extensions
	not input.extensions[ext_name]
	msg := sprintf("OTEL-062: extension '%s' is in service.extensions but not defined in extensions section.", [ext_name])
}
