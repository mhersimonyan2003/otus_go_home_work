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

lint: golangci-lint
	@echo "Running golangci-lint..."
	@$(GOLANGCI_LINT) run

golangci-lint:
ifeq (, $(GOLANGCI_LINT))
	@echo "Installing golangci-lint..."
	@curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(HOME)/bin v1.50.1
	@export PATH=$(HOME)/bin:$PATH
else
	@echo "golangci-lint is already installed."
endif

lint-fix: golangci-lint
	@echo "Running golangci-lint --fix..."
	@$(GOLANGCI_LINT) run --fix

clean:
	@echo "Cleaning up..."
	@rm -rf $(BIN_DIR)
	@echo "Cleanup completed."
