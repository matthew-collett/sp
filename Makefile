BINARY=sp

VERSION := $(shell git describe --tags --abbrev=0 2>/dev/null || echo "dev")
DATE := $(shell date +%Y-%m-%d)
LDFLAGS=-ldflags "-s -w -X github.com/matthew-collett/sp/internal/build.v=$(VERSION) -X github.com/matthew-collett/sp/internal/build.d=$(DATE)"

INSTALL_PATH := /usr/local/bin
TARGET ?= all

.PHONY: build test clean coverage install uninstall help

help:
	@echo "Targets:"
	@echo "  build      Build binaries to build/ (TARGET=sp|sp-mcp to build one)"
	@echo "  test       Run all tests with verbose output"
	@echo "  clean      Remove build directory and coverage output"
	@echo "  coverage   Run tests and open HTML coverage report"
	@echo "  install    Build and install to $(INSTALL_PATH)"
	@echo "  uninstall  Remove from $(INSTALL_PATH)"
	@echo "  help       Show this help message"

build:
	@mkdir -p build
ifeq ($(TARGET),sp-mcp)
	go build $(LDFLAGS) -o build/$(BINARY)-mcp ./cmd/$(BINARY)-mcp
else ifeq ($(TARGET),sp)
	go build $(LDFLAGS) -o build/$(BINARY) ./cmd/$(BINARY)
else
	go build $(LDFLAGS) -o build/$(BINARY) ./cmd/$(BINARY)
	go build $(LDFLAGS) -o build/$(BINARY)-mcp ./cmd/$(BINARY)-mcp
endif

install: build
ifeq ($(TARGET),sp-mcp)
	cp build/$(BINARY)-mcp $(INSTALL_PATH)/$(BINARY)-mcp
	chmod +x $(INSTALL_PATH)/$(BINARY)-mcp
	@echo $(BINARY)-mcp installed to $(INSTALL_PATH)
else ifeq ($(TARGET),sp)
	cp build/$(BINARY) $(INSTALL_PATH)/$(BINARY)
	chmod +x $(INSTALL_PATH)/$(BINARY)
	@echo $(BINARY) installed to $(INSTALL_PATH)
else
	cp build/$(BINARY) $(INSTALL_PATH)/$(BINARY)
	chmod +x $(INSTALL_PATH)/$(BINARY)
	cp build/$(BINARY)-mcp $(INSTALL_PATH)/$(BINARY)-mcp
	chmod +x $(INSTALL_PATH)/$(BINARY)-mcp
	@echo $(BINARY) and $(BINARY)-mcp installed to $(INSTALL_PATH)
endif

uninstall:
ifeq ($(TARGET),sp-mcp)
	rm -f $(INSTALL_PATH)/$(BINARY)-mcp
	@echo $(BINARY)-mcp uninstalled from $(INSTALL_PATH)
else ifeq ($(TARGET),sp)
	rm -f $(INSTALL_PATH)/$(BINARY)
	@echo $(BINARY) uninstalled from $(INSTALL_PATH)
else
	rm -f $(INSTALL_PATH)/$(BINARY) $(INSTALL_PATH)/$(BINARY)-mcp
	@echo $(BINARY) and $(BINARY)-mcp uninstalled from $(INSTALL_PATH)
endif

test:
	go test -v ./...

clean:
	rm -rf build
	rm -f coverage.out

coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

.DEFAULT_GOAL := help
