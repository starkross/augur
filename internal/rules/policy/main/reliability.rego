package main

import data.lib
import future.keywords.contains
import future.keywords.if
import future.keywords.in

# Reliability and sampling rules

# OTEL-063: tail_sampling without groupbytrace
warn contains msg if {
	some ptype in lib.pipeline_types
	procs := lib.pipeline_processors(ptype)
	_pipeline_has_processor_type(procs, "tail_sampling")
	not _pipeline_has_processor_type(procs, "groupbytrace")
	msg := sprintf(
		"OTEL-063: pipeline '%s' uses tail_sampling without groupbytrace. Incomplete traces cause inconsistent sampling.",
		[ptype],
	)
}

# OTEL-064: both probabilistic_sampler and tail_sampling
warn contains msg if {
	some ptype in lib.pipeline_types
	procs := lib.pipeline_processors(ptype)
	_pipeline_has_processor_type(procs, "probabilistic_sampler")
	_pipeline_has_processor_type(procs, "tail_sampling")
	msg := sprintf(
		"OTEL-064: pipeline '%s' has both probabilistic_sampler and tail_sampling. Use one sampler only.",
		[ptype],
	)
}

# OTEL-065: sending_queue without persistent storage
warn contains msg if {
	some name, exporter in input.exporters
	exporter.sending_queue.enabled != false
	not exporter.sending_queue.storage
	msg := sprintf(
		"OTEL-065: exporter '%s' sending_queue has no persistent storage. Data lost on restart.",
		[name],
	)
}

# OTEL-067: K8s environment without k8sattributes
warn contains msg if {
	not _has_processor_type("k8sattributes")
	some name in object.keys(input.receivers)
	split(name, "/")[0] in {"kubeletstats", "k8s_cluster", "k8s_events", "k8sobjects"}
	msg := "OTEL-067: K8s environment detected but k8sattributes processor is not configured."
}

# OTEL-068: K8s environment without resourcedetection
warn contains msg if {
	not _has_processor_type("resourcedetection")
	some name in object.keys(input.receivers)
	split(name, "/")[0] in {"kubeletstats", "k8s_cluster", "k8s_events", "k8sobjects"}
	msg := "OTEL-068: K8s environment detected but resourcedetection processor is not configured."
}

# ============================================================
# Private helpers
# ============================================================

_has_processor_type(ptype) if {
	some name in object.keys(input.processors)
	split(name, "/")[0] == ptype
}
