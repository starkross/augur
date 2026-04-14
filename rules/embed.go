package rules

import "embed"

// Policies contains the built-in Rego policy files, embedded at compile time.
//
//go:embed policy
var Policies embed.FS

// PolicyDir is the root directory within Policies containing the .rego files.
const PolicyDir = "policy"
