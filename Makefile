BINARY     := augur
VERSION    := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS    := -s -w -X main.version=$(VERSION)
POLICY_SRC := policy
POLICY_DST := internal/rules/policy

.PHONY: build test test-rego lint-rego snapshot install clean sync-policies demo

sync-policies:
	@rm -rf $(POLICY_DST)
	@mkdir -p $(POLICY_DST)
	@cp -R $(POLICY_SRC)/main $(POLICY_SRC)/lib $(POLICY_DST)/

build: sync-policies
	CGO_ENABLED=0 go build -ldflags "$(LDFLAGS)" -o $(BINARY) ./cmd/augur

test: sync-policies
	go test ./...

test-rego:
	@command -v conftest >/dev/null 2>&1 || { echo "conftest needed for rego tests"; exit 1; }
	conftest verify --policy $(POLICY_SRC)/

lint-rego:
	@command -v regal >/dev/null 2>&1 || { echo "regal needed: brew install styrainc/packages/regal"; exit 1; }
	regal lint $(POLICY_SRC)/

snapshot: sync-policies
	goreleaser build --snapshot --clean

install: build
	cp $(BINARY) /usr/local/bin/$(BINARY)
	@echo "✓ $(BINARY) $(VERSION) installed"

demo: build
	@echo "=== Good config ===" && ./$(BINARY) examples/good.yaml || true
	@echo "" && echo "=== Bad config ===" && ./$(BINARY) examples/bad.yaml || true

clean:
	rm -f $(BINARY)
	rm -rf dist/ $(POLICY_DST)
