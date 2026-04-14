# Contributing

## Development

```sh
# Run Go tests (syncs policies automatically)
make test

# Run Rego unit tests (requires conftest)
make test-rego

# Lint Rego policies (requires regal)
make lint-rego

# Build binary
make build

# Run against example configs
make demo
```

## Adding a new rule

1. Create or edit a `.rego` file in `rules/policy/main/` (rules are grouped by category)
2. Add the rule as `deny contains msg if { ... }` (blocking) or `warn contains msg if { ... }` (advisory)
3. Format the message as `"OTEL-NNN: description."` — the engine extracts the rule ID from this prefix
4. Add a test in the matching `*_test.rego` file
5. Run `make test && make test-rego && make lint-rego` to verify
