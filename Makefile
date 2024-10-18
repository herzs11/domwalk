.DEFAULT_GOAL := install

EXECUTABLE := domwalk
ARCH := $(shell uname -m)

install-go:
	@if ! command -v go &> /dev/null; then \
		echo "Go is not installed. Installing..."; \
		if [ "$(ARCH)" == "x86_64" ]; then \
			wget "https://go.dev/dl/go1.23.1.darwin-amd64.pkg"; \
			sudo installer -pkg go1.23.1.darwin-amd64.pkg -target /; \
			rm go1.23.1.darwin-amd64.pkg; \
		elif [ "$(ARCH)" == "arm64" ]; then \
			wget "https://go.dev/dl/go1.23.1.darwin-arm64.pkg"; \
			sudo installer -pkg go1.23.1.darwin-arm64.pkg -target /; \
			rm go1.23.1.darwin-arm64.pkg; \
		else \
			echo "Unsupported architecture: $(ARCH)"; \
			exit 1; \
		fi; \
		echo "Go installed successfully."; \
	fi

get-deps:
	go mod download

install-cli:
	GOARCH=$(ARCH) go install .

test-cf:
	go run ./cloud_functions/cloud_functions_test/main.go