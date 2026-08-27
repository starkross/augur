BINARY     := augur
VERSION    := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS    := -s -w -X main.version=$(VERSION)
POLICY_SRC := rules/policy
WASM_DIR   := dist/wasm
WASM_ENTRY := -e main/deny -e main/warn

.PHONY: build test test-rego lint-rego snapshot install clean demo wasm wasm-policy

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

# Compile the Rego rules to a standalone wasm module via OPA's Rego->wasm
# compiler. Contains the policies only -- no interpreter, no Go runtime -- so
# any OPA-ABI host (JS, Rust, Python, Go, ...) can evaluate the rules. This is
# the artifact to ship to a browser.
wasm-policy:
	@command -v opa >/dev/null 2>&1 || { echo "opa needed: brew install opa"; exit 1; }
	@rm -rf $(WASM_DIR)/bundle && mkdir -p $(WASM_DIR)/bundle
	opa build -t wasm $(WASM_ENTRY) --ignore '*_test.rego' -o $(WASM_DIR)/bundle/bundle.tar.gz $(POLICY_SRC)
	@tar -xzf $(WASM_DIR)/bundle/bundle.tar.gz -C $(WASM_DIR)/bundle 2>/dev/null
	@mv $(WASM_DIR)/bundle/policy.wasm $(WASM_DIR)/policy.wasm
	@rm -rf $(WASM_DIR)/bundle
	@echo "✓ $(WASM_DIR)/policy.wasm ($$(du -h $(WASM_DIR)/policy.wasm | cut -f1))"

# Build the full linter -- CLI, YAML loading, env expansion and the embedded
# OPA interpreter -- as a WASI preview 1 module, runnable under wasmtime,
# wasmer or node. Self-contained but large: the bundled interpreter dominates.
wasm:
	@mkdir -p $(WASM_DIR)
	GOOS=wasip1 GOARCH=wasm go build -ldflags "$(LDFLAGS)" -o $(WASM_DIR)/$(BINARY).wasm ./cmd/augur
	@echo "✓ $(WASM_DIR)/$(BINARY).wasm ($$(du -h $(WASM_DIR)/$(BINARY).wasm | cut -f1))"

install: build
	cp $(BINARY) /usr/local/bin/$(BINARY)
	@echo "✓ $(BINARY) $(VERSION) installed"

demo: build
	@echo "=== Good config ===" && ./$(BINARY) examples/good.yaml || true
	@echo "" && echo "=== Bad config ===" && ./$(BINARY) examples/bad.yaml || true

clean:
	rm -f $(BINARY)
	rm -rf dist/
