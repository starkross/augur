BINARY     := augur
VERSION    := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS    := -s -w -X main.version=$(VERSION)
POLICY_SRC := rules/policy

.PHONY: build test test-rego lint-rego snapshot install clean demo \
	rust-build rust-test rust-wasm difftest

build:
	CGO_ENABLED=0 go build -ldflags "$(LDFLAGS)" -o $(BINARY) ./cmd/augur

test:
	go test ./...

test-rego:
	@command -v conftest >/dev/null 2>&1 || { echo "conftest needed for rego tests"; exit 1; }
	conftest verify --policy $(POLICY_SRC)/

lint-rego:
	@command -v regal >/dev/null 2>&1 || { echo "regal needed: brew install styrainc/packages/regal"; exit 1; }
	regal lint $(POLICY_SRC)/

snapshot:
	goreleaser build --snapshot --clean

rust-build:
	cd rust && cargo build --release
	@echo "✓ rust/target/release/augur ($$(du -h rust/target/release/augur | cut -f1))"

rust-test:
	cd rust && cargo test

# The Rust linter as a WASI module: same rules, ~40x smaller than the Go build
# because it carries regorus instead of the OPA interpreter.
rust-wasm:
	cd rust && cargo build --release --target wasm32-wasip1
	@echo "✓ rust/target/wasm32-wasip1/release/augur.wasm ($$(du -h rust/target/wasm32-wasip1/release/augur.wasm | cut -f1))"

# Both implementations must agree byte-for-byte across the corpus.
difftest: build rust-build
	./rust/difftest.sh

install: build
	cp $(BINARY) /usr/local/bin/$(BINARY)
	@echo "✓ $(BINARY) $(VERSION) installed"

demo: build
	@echo "=== Good config ===" && ./$(BINARY) examples/good.yaml || true
	@echo "" && echo "=== Bad config ===" && ./$(BINARY) examples/bad.yaml || true

clean:
	rm -f $(BINARY)
	rm -rf dist/
