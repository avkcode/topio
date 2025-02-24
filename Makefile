# Variables
BIN_NAME := proc-monitor
SRC_FILES := $(wildcard *.go)
BUILD_DIR := build
TEST_FLAGS := -v -coverprofile=coverage.out

# Default target
all: build

# Build the binary
build: $(BUILD_DIR)/$(BIN_NAME)

# Run the binary
run: build
	@echo "Running $(BIN_NAME)..."
	./$(BUILD_DIR)/$(BIN_NAME) -interval=1 -html=true -port=8080

# Test the code
test:
	@echo "Running tests..."
	go test $(TEST_FLAGS) ./...

# Generate coverage report
coverage: test
	@echo "Generating coverage report..."
	go tool cover -html=coverage.out -o coverage.html

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -rf $(BUILD_DIR) coverage.out coverage.html

# Install the binary to /usr/local/bin
install: build
	@echo "Installing $(BIN_NAME) to /usr/local/bin..."
	sudo cp $(BUILD_DIR)/$(BIN_NAME) /usr/local/bin/

# Uninstall the binary from /usr/local/bin
uninstall:
	@echo "Uninstalling $(BIN_NAME)..."
	sudo rm -f /usr/local/bin/$(BIN_NAME)

# Create the build directory
$(BUILD_DIR):
	mkdir -p $(BUILD_DIR)

# Compile the binary
$(BUILD_DIR)/$(BIN_NAME): $(SRC_FILES) | $(BUILD_DIR)
	@echo "Building $(BIN_NAME)..."
	go build -o $(BUILD_DIR)/$(BIN_NAME) .

.PHONY: all build run test coverage clean install uninstall
