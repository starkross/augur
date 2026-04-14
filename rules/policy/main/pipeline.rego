package main

import data.lib
import future.keywords.contains
import future.keywords.if
import future.keywords.in

# Pipeline architecture rules

# OTEL-040: circular pipeline dependency via connectors
deny contains msg if {
	input.connectors
	some pipeline in lib.pipeline_types
	targets := _connector_targets[pipeline]
	count(targets) > 0
	reachable := graph.reachable(_connector_targets, targets)
	pipeline in reachable
	msg := sprintf("OTEL-040: circular pipeline dependency detected involving pipeline '%s'.", [pipeline])
}

# OTEL-038: filter processor after batch
warn contains msg if {
	some ptype in lib.pipeline_types
	procs := lib.pipeline_processors(ptype)
	some i, p in procs
	split(p, "/")[0] == "batch"
	some j, q in procs
	split(q, "/")[0] == "filter"
	j > i
	msg := sprintf("OTEL-038: pipeline '%s' has filter processor after batch. Filter early to reduce waste.", [ptype])
}

# OTEL-039: transform/attributes processor after batch
warn contains msg if {
	transform_types := {"transform", "attributes"}
	some ptype in lib.pipeline_types
	procs := lib.pipeline_processors(ptype)
	some i, p in procs
	split(p, "/")[0] == "batch"
	some j, q in procs
	split(q, "/")[0] in transform_types
	j > i
	msg := sprintf("OTEL-039: pipeline '%s' has transform/attributes processor after batch.", [ptype])
}

# OTEL-041: routing connector without default_pipelines
warn contains msg if {
	some name, _ in input.connectors
	split(name, "/")[0] == "routing"
	not input.connectors[name].default_pipelines
	msg := sprintf(
		"OTEL-041: routing connector '%s' has no default_pipelines. Unmatched data is silently dropped.",
		[name],
	)
}

# OTEL-042: duplicate processor in same pipeline
warn contains msg if {
	some ptype in lib.pipeline_types
	procs := lib.pipeline_processors(ptype)
	some i, p in procs
	some j, q in procs
	p == q
	i < j
	msg := sprintf("OTEL-042: pipeline '%s' has duplicate processor '%s'.", [ptype, p])
}

# OTEL-043: batch before tail_sampling or groupbytrace
warn contains msg if {
	sampling_types := {"tail_sampling", "groupbytrace"}
	some ptype in lib.pipeline_types
	procs := lib.pipeline_processors(ptype)
	some i, p in procs
	split(p, "/")[0] == "batch"
	some j, q in procs
	split(q, "/")[0] in sampling_types
	i < j
	msg := sprintf("OTEL-043: pipeline '%s' has batch before %s. Batch splits traces across batches.", [ptype, q])
}

# ============================================================
# Private helpers
# ============================================================

# Connector pipeline graph for cycle detection (OTEL-040)
_connector_targets[p] := ts if {
	some p in lib.pipeline_types
	ts := {t |
		some exp in lib.pipeline_exporters(p)
		input.connectors[exp]
		some t in lib.pipeline_types
		some rcv in lib.pipeline_receivers(t)
		rcv == exp
	}
}

_pipeline_has_processor_type(procs, ptype) if {
	some p in procs
	split(p, "/")[0] == ptype
}
