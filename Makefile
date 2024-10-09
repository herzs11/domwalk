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


# Install the CLI tool (including Go installation and build)
install: install-go install-cli
	@if [ -z "$$GOOGLE_APPLICATION_CREDENTIALS" ]; then \
		echo "Error: \$GOOGLE_APPLICATION_CREDENTIALS environment variable is not set."; \
		echo "Please set it to the path of your Google Cloud service account key file."; \
		exit 1; \
	fi
	@read -p "Enter the path to your SQLite database (default: $$HOME/.domwalk.db): " DB_PATH; \
	if [ -z "$$DB_PATH" ]; then \
		DB_PATH="$$HOME/.domwalk.db"; \
	fi; \
	if [ -n "$ZSH_VERSION" ]; then \
		echo "export DOMWALK_SQLITE_NAME=\"$$DB_PATH\"" >> ~/.zshrc; \
		(zsh -c "source ~/.zshrc"); \
	elif [ -n "$BASH_VERSION" ]; then \
		echo "export DOMWALK_SQLITE_NAME=\"$$DB_PATH\"" >> ~/.bash_profile; \
		(bash -c "source ~/.bash_profile"); \
	else \
		echo "Warning: Unsupported shell. Please manually set the DOMWALK_SQLITE_NAME environment variable."; \
	fi
	@echo "Installation complete!";
	@echo "The DOMWALK_SQLITE_NAME environment variable has been set (or instructions provided).";
