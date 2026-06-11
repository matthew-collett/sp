BINARY=sp

VERSION := $(shell git describe --tags --abbrev=0 2>/dev/null || echo "dev")
DATE := $(shell date +%Y-%m-%d)
LDFLAGS=-ldflags "-s -w -X github.com/matthew-collett/sp/internal/build.v=$(VERSION) -X github.com/matthew-collett/sp/internal/build.d=$(DATE)"

INSTALL_PATH := /usr/local/bin

.PHONY: build test clean run coverage install uninstall help

help:
	@echo "Targets:"
	@echo "  build      Build the binary to build/$(BINARY)"
	@echo "  test       Run all tests with verbose output"
	@echo "  clean      Remove build directory and coverage output"
	@echo "  coverage   Run tests and open HTML coverage report"
	@echo "  install    Build and install $(BINARY) to $(INSTALL_PATH)"
	@echo "  uninstall  Remove $(BINARY) from $(INSTALL_PATH)"
	@echo "  help       Show this help message"

build:
	go build $(LDFLAGS) -o build/$(BINARY) ./cmd/$(BINARY)
	go build $(LDFLAGS) -o build/$(BINARY)-mcp ./cmd/$(BINARY)-mcp

test:
	go test -v ./...

clean:
	rm -rf build
	rm -f coverage.out

coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

install: build
	@mkdir -p $(INSTALL_PATH)
	cp build/$(BINARY) $(INSTALL_PATH)/$(BINARY)
	chmod +x $(INSTALL_PATH)/$(BINARY)
	@echo $(BINARY) installed to $(INSTALL_PATH)

uninstall:
	rm -f $(INSTALL_PATH)/$(BINARY)
	@echo $(BINARY) uninstalled from $(INSTALL_PATH)

.DEFAULT_GOAL := help
