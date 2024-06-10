# Variables
APP_NAME = calendar
CONFIG_FILE = configs/config.toml
GOLANGCI_LINT = $(shell which golangci-lint)

# Directories
CMD_DIR = ./cmd/calendar
BIN_DIR = ./bin
BUILD_DIR = $(BIN_DIR)/$(APP_NAME)

# Commands
.PHONY: all build run test lint clean

all: build

build:
	@echo "Building binary..."
	@mkdir -p $(BIN_DIR)
	@go build -o $(BUILD_DIR) $(CMD_DIR)
	@echo "Build completed. Binary is located at $(BUILD_DIR)"

run: build
	@echo "Running the service..."
	@$(BUILD_DIR) --config=$(CONFIG_FILE)

test:
	@echo "Running tests..."
	@go test ./... -race -v

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.50.1

lint: install-lint-deps
	@echo "Running golangci-lint..."
	@golangci-lint run


lint-fix: install-lint-deps
	@echo "Running golangci-lint --fix..."
	@golangci-lint run --fix

clean:
	@echo "Cleaning up..."
	@rm -rf $(BIN_DIR)
	@echo "Cleanup completed."
