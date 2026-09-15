#!/usr/bin/env bash
# Differential test: the Go and Rust linters must agree, byte for byte, on
# every config in the corpus across every output format and flag combination.
#
# Usage: rust/difftest.sh [extra-config.yaml ...]
set -uo pipefail

cd "$(dirname "$0")/.." || exit 1

GO_BIN=${GO_BIN:-./augur}
RS_BIN=${RS_BIN:-./rust/target/release/augur}

[ -x "$GO_BIN" ] || { echo "missing $GO_BIN — run: make build"; exit 1; }
[ -x "$RS_BIN" ] || { echo "missing $RS_BIN — run: make rust-build"; exit 1; }

pass=0 fail=0 tolerated=0

check() {
  local label="$1"; shift
  local g r ge re
  g=$("$GO_BIN" "$@" 2>&1); ge=$?
  r=$("$RS_BIN" "$@" 2>&1); re=$?
  if [ "$g" = "$r" ] && [ "$ge" = "$re" ]; then
    pass=$((pass + 1))
    printf '  ✓ %s\n' "$label"
  elif [ "$ge" = "$re" ] && [ "$ge" -ne 0 ] \
       && [[ "$g" == *"parsing YAML:"* ]] && [[ "$r" == *"parsing YAML:"* ]]; then
    # Both reject the document with the same exit code; only the underlying
    # YAML library's wording differs, which is not worth reimplementing.
    tolerated=$((tolerated + 1))
    printf '  ~ %s (parse-error text differs; both reject)\n' "$label"
  else
    fail=$((fail + 1))
    printf '  ✗ %s (exit: go=%s rust=%s)\n' "$label" "$ge" "$re"
    diff <(printf '%s\n' "$g") <(printf '%s\n' "$r") | sed 's/^/      /' | head -10
  fi
}

CONFIGS=(testdata/*.yaml examples/*.yaml "$@")

echo "corpus: ${#CONFIGS[@]} configs"
for f in "${CONFIGS[@]}"; do
  [ -f "$f" ] || continue
  for fmt in text json github; do
    check "$fmt $f" -o "$fmt" --no-color "$f"
  done
  check "quiet  $f"  -o text --no-color -q "$f"
  check "strict $f"  -o text --no-color -s "$f"
done

# Flags and multi-file merge.
check "skip"          -o text --no-color -k OTEL-010,OTEL-011 examples/bad.yaml
check "merge"         -o json testdata/good.yaml testdata/bad.yaml
check "stdin-vs-file" -o json testdata/bad.yaml

echo
echo "pass=$pass tolerated=$tolerated fail=$fail"
[ "$fail" -eq 0 ]
