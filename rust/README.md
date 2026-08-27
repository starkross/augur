# augur — Rust implementation

A second implementation of the linter, sharing the Rego rules with the Go one.
`rules/policy/` remains the single source of truth: `augur-core/build.rs` walks
that directory at compile time, so a new `.rego` file is picked up by both
without touching Rust.

The motivation is wasm. The Go binary embeds the OPA interpreter, which is
~96% of a `wasip1` build; swapping it for [regorus] and trimming unused
features brings the module from 46.6 MB to 1.2 MB.

| build | raw | brotli |
|---|---:|---:|
| Go + OPA (`make wasm`) | 46,593,356 | 6.04 MB |
| Rust + regorus (`make rust-wasm`) | **1,184,533** | **330 KB** |

## Layout

- `augur-core` — parsing, env substitution, merge, evaluation. No I/O: callers
  pass strings, so it works unchanged in a CLI, a server, or a browser.
- `augur-cli` — the `augur` binary. Flags and output are compatible with the Go
  CLI.

## Building

```sh
make rust-build      # native binary at rust/target/release/augur
make rust-test       # cargo test
make rust-wasm       # wasm32-wasip1 module
make difftest        # Go vs Rust, byte-for-byte over the corpus
```

## Fidelity

`difftest.sh` runs both binaries over every config in `testdata/` and
`examples/` across all three output formats and the flag combinations, and
requires byte-identical stdout and matching exit codes. Current status:
**28 pass, 5 tolerated, 0 fail**.

### Deliberate matches

Two behaviours had to be reproduced rather than "fixed", because diverging
changes lint results:

**YAML 1.1 scalars.** `sigs.k8s.io/yaml` wraps go-yaml v2, which resolves
`yes`/`no`/`on`/`off`/`y`/`n` as booleans. Rust YAML crates follow the YAML 1.2
core schema and leave them as strings. Without coercion a config with
`insecure_skip_verify: yes` silently stops matching `== true`, so OTEL-032 —
"TLS verification is bypassed" — goes unreported. The coercion applies to
mapping *keys* too, because go-yaml does: a bare `y:` key becomes `true`.

**Builtin error handling.** OPA leaves an expression undefined when a builtin
hits a runtime type error and continues evaluating; regorus aborts by default.
An empty config section (`grpc:` with no body) parses as null, reaches
`object.get`, and under the strict default takes the whole run down — 14
findings become 0. `set_strict_builtin_errors(false)` restores OPA's behaviour.

`"findings": null` for a clean config is also preserved: Go marshals a nil
slice that way, and emitting `[]` would change the shape consumers parse.

### Known divergences

Tolerated by `difftest.sh`:

- **Parse-error wording.** Both implementations reject malformed YAML with the
  same exit code, but the message comes from the underlying library and differs.

Not currently reconciled, and not covered by the corpus — each needs a decision
before the Rust CLI could replace the Go one:

| case | Go | Rust |
|---|---|---|
| `0644`, `010` | `420`, `8` (YAML 1.1 octal) | `"0644"`, `"010"` |
| duplicate keys | last wins | parse error |
| `<<:` merge key | merged | literal `<<` key |
| ints > 2^53 | precision lost via float64 | exact |

The last one is arguably a Go bug worth fixing rather than copying. The corpus
is also small — six configs — so this list should be treated as what has been
looked for, not as proof of equivalence.

[regorus]: https://github.com/microsoft/regorus
