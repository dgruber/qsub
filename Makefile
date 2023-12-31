BINARY_NAME := qsub
BINARY_NAME_QSTAT := qstat
BUILD_DIR := build
LDFLAGS := -ldflags="-s -w" # make it small

.PHONY: all clean build build-darwin build-linux

all: clean build

build: build-darwin build-linux

build-darwin:
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)_darwin
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/${BINARY_NAME_QSTAT}_darwin ./cmd/qstat/

build-linux:
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)_linux
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/${BINARY_NAME_QSTAT}_linux ./cmd/qstat/

clean:
	rm -rf $(BUILD_DIR)
	mkdir $(BUILD_DIR)


