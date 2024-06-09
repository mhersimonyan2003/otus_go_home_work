# Variables
APP_NAME = calendar
CONFIG_FILE = configs/config.toml

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

lint:
	@echo "Running golangci-lint..."
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@golangci-lint run

lint-fix:
	@echo "Running golangci-lint --fix..."
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@golangci-lint run --fix

clean:
	@echo "Cleaning up..."
	@rm -rf $(BIN_DIR)
	@echo "Cleanup completed."
