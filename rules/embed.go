package rules

import "embed"

//go:embed policy
var Policies embed.FS

const PolicyDir = "policy"
