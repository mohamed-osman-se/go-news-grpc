
# Local bin for tools
GO_BIN?=$(shell pwd)/.bin

# golangci-lint version
GOCI_LINT_VERSION?=v1.64.5

# Use local bin in PATH
SHELL:=env PATH=$(GO_BIN):$(PATH) $(SHELL)

# ----------------------------------------
# Format the Go code
format::
	@echo ">> Formatting Go code..."
	golangci-lint run --fix -v ./...

# ----------------------------------------
# Generate Go code from proto files
generate-proto::
	@echo ">> Generating Go code from proto files..."
	go tool buf generate --template ./proto/news/v1/buf.gen.yml

# ----------------------------------------
# Detect breaking changes in proto
lint-breaking::
	@echo ">> Checking for breaking changes in proto..."
	go tool buf breaking --against ''

# ----------------------------------------
# Lint Go code
lint-go::
	@echo ">> Linting Go code..."
	golangci-lint run -v ./...

# ----------------------------------------
# Lint proto files
lint-proto::
	@echo ">> Linting proto files..."
	go tool buf lint --config ./buf.yaml

# ----------------------------------------
# Run all linters
lint:: lint-breaking lint-go lint-proto

# ----------------------------------------
# Install necessary tools locally
install-tools::
	@echo ">> Installing tools..."
	mkdir -p ${GO_BIN}
	# Install golangci-lint
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b ${GO_BIN} ${GOCI_LINT_VERSION}
	# Install buf if not installed
	curl -sSL https://github.com/bufbuild/buf/releases/download/v1.25.0/buf-Linux-x86_64 -o ${GO_BIN}/buf
	chmod +x ${GO_BIN}/buf
	@echo ">> Tools installed in ${GO_BIN}"

# ----------------------------------------
# Clean and tidy Go modules
tidy::
	@echo ">> Running go mod tidy..."
	go mod tidy -v

# ----------------------------------------
# Setup environment, generate code, lint, and run server & client
setup-run:: install-tools tidy generate-proto lint
	@echo ">> Starting gRPC server in background..."
	# Run server in background
	cd cmd/server && nohup go run main.go > ../../server.log 2>&1 &
	@echo ">> Server is running in background (logs: server.log)"
	# Give server a moment to start


