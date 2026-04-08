---
id: custom-policies
title: Custom policies
sidebar_position: 6
---

# Custom policies

All built-in rules live in [`policy/`](https://github.com/starkross/augur/tree/main/policy) as standard [Rego](https://www.openpolicyagent.org/docs/latest/policy-language/) files. To add your own:

## 1. Write a rule

```rego
# my-policies/main/custom.rego
package main

import future.keywords.contains
import future.keywords.if

deny contains msg if {
    not input.processors.filter
    msg := "CUSTOM-001: filter processor is required by our platform team."
}
```

## 2. Run augur with `--policy`

```sh
augur --policy ./my-policies config.yaml
```

Custom policies are **merged** with the built-in rules — your rules run alongside every default check.

## Rule conventions

- `deny contains msg` — blocking rule, fails the run
- `warn contains msg` — advisory rule, reported but non-blocking
- Prefix message IDs with your own namespace (e.g. `ACME-001`) to avoid colliding with augur's `OTEL-*` IDs
- Keep messages actionable: state what's wrong AND what to do about it

## Testing your rules

Rego ships with a built-in test runner. Put tests next to your rules:

```rego
# my-policies/main/custom_test.rego
package main

test_custom_001_denies_missing_filter if {
    result := deny with input as {"processors": {}}
    count(result) == 1
}
```

Run them with:

```sh
opa test my-policies/
```
