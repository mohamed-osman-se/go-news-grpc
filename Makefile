
# Local bin for tools
GO_BIN?=$(shell pwd)/.bin

# golangci-lint version
GOCI_LINT_VERSION?=v1.64.5

# Use local bin in PATH
SHELL:=env PATH=$(GO_BIN):$(PATH) $(SHELL)



# ----------------------------------------
# Generate Go code from proto files
generate-proto::
	@echo ">> Generating Go code from proto files..."
	go tool buf generate --template ./proto/news/v1/buf.gen.yml




# ----------------------------------------
# Install necessary tools locally
install-tools::
	@echo ">> Installing tools..."
	mkdir -p ${GO_BIN}
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
setup-run:: install-tools tidy generate-proto 
	@echo ">>All Done You Can Run The Server. "
	


