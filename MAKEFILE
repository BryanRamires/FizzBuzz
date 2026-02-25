.PHONY: tidy test race build run exec clean tools lint vuln ci fmt

APP_NAME := fizzbuzz-api
BIN_DIR  := bin
PKG      := ./...

# ---- Tools installation (local, pinned) ----
GOLANGCI_LINT_VERSION := v2.10.1
GOVULNCHECK_VERSION := v1.1.4

tidy:
	go mod tidy

fmt:
	gofmt -w .

build:
	go build -o $(BIN_DIR)/$(APP_NAME) ./cmd/api

test:
	go test $(PKG) -count=1

race:
	CGO_ENABLED=1 go test $(PKG) -race -count=1

cover:
	go test $(PKG) -count=1 -coverprofile=$(BIN_DIR)/coverage.out
	go tool cover -func=$(BIN_DIR)/coverage.out

run:
	go run ./cmd/api

exec: build
	$(BIN_DIR)/$(APP_NAME)

# ---- Tools installation (local, reproducible) ----
tools:
	@mkdir -p $(BIN_DIR)

	@if [ ! -f ./$(BIN_DIR)/golangci-lint ]; then \
		echo "Installing golangci-lint $(GOLANGCI_LINT_VERSION) into ./$(BIN_DIR) ..."; \
		GOBIN=$$(pwd)/$(BIN_DIR) go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION); \
	fi

	@if [ ! -f ./$(BIN_DIR)/govulncheck ]; then \
		echo "Installing govulncheck $(GOVULNCHECK_VERSION) into ./$(BIN_DIR) ..."; \
		GOBIN=$$(pwd)/$(BIN_DIR) go install golang.org/x/vuln/cmd/govulncheck@$(GOVULNCHECK_VERSION); \
	fi

lint: tools
	./$(BIN_DIR)/golangci-lint run --timeout=5m

vuln: tools
	GOTOOLCHAIN=local ./$(BIN_DIR)/govulncheck ./...

# "CI parity" target: run what you expect in GitHub Actions
ci: tools tidy fmt test race lint vuln