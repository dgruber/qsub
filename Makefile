BINARY_NAME := qsub
BUILD_DIR := build
LDFLAGS := -ldflags="-s -w" # make it small

.PHONY: all clean build build-darwin build-linux

all: clean build

build: build-darwin build-linux

build-darwin:
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)_darwin

build-linux:
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)_linux

clean:
	rm -rf $(BUILD_DIR)
	mkdir $(BUILD_DIR)


