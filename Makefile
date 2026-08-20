APP_NAME := sandbox-daemon
BIN_DIR := bin

.PHONY: all build run clean

all: build

build:
	@mkdir -p $(BIN_DIR)
	@go build -o $(BIN_DIR)/$(APP_NAME) ./cmd

run:
	@go run ./cmd

clean:
	@rm -rf $(BIN_DIR)
