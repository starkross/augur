package rules

import "embed"

// Policies embeds all .rego files from the policy/ directory at compile time.
// Before building, run: cp policy/*.rego internal/rules/policy/
//
//go:embed policy
var Policies embed.FS

// PolicyDir is the root directory within the embedded filesystem containing policy files.
const PolicyDir = "policy"
