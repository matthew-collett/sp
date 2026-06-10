BINARY=sp

VERSION := $(shell git describe --tags --abbrev=0 2>/dev/null || echo "dev")
DATE := $(shell date +%Y-%m-%d)
LDFLAGS=-ldflags "-s -w -X github.com/matthew-collett/sp/internal/build.v=$(VERSION) -X github.com/matthew-collett/sp/internal/build.d=$(DATE)"

# Installation paths
ifeq ($(OS),Windows_NT)
	RM_DIR := rmdir /s /q
	RM_FILE := del /f /q
	RUN_BINARY := build\$(BINARY)
	INSTALL_PATH := $(USERPROFILE)\bin
	BINARY_EXT := .exe
	MKDIR := mkdir
else
	RM_DIR := rm -rf
	RM_FILE := rm -f
	RUN_BINARY := ./build/$(BINARY)
	INSTALL_PATH := /usr/local/bin
	BINARY_EXT :=
	MKDIR := mkdir -p
endif

.PHONY: build test clean run coverage install uninstall help

help:
	@echo "Targets:"
	@echo "  build      Build the binary to build/$(BINARY)$(BINARY_EXT) (default)"
	@echo "  test       Run all tests with verbose output"
	@echo "  clean      Remove build directory and coverage output"
	@echo "  coverage   Run tests and open HTML coverage report"
	@echo "  install    Build and install $(BINARY) to $(INSTALL_PATH)"
	@echo "  uninstall  Remove $(BINARY) from $(INSTALL_PATH)"
	@echo "  help       Show this help message"

build:
	go build $(LDFLAGS) -o build/$(BINARY)$(BINARY_EXT) ./cmd/$(BINARY)

test:
	go test -v ./...

clean:
	$(RM_DIR) build
	$(RM_FILE) coverage.out

coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

install: build
ifeq ($(OS),Windows_NT)
	@if not exist "$(INSTALL_PATH)" $(MKDIR) "$(INSTALL_PATH)"
	copy "build\$(BINARY)$(BINARY_EXT)" "$(INSTALL_PATH)\"
	@echo $(BINARY) installed to $(INSTALL_PATH)
	@echo Make sure $(INSTALL_PATH) is in your PATH environment variable
else
	@$(MKDIR) $(INSTALL_PATH)
	cp build/$(BINARY) $(INSTALL_PATH)/$(BINARY)
	chmod +x $(INSTALL_PATH)/$(BINARY)
	@echo $(BINARY) installed to $(INSTALL_PATH)
endif

uninstall:
ifeq ($(OS),Windows_NT)
	$(RM_FILE) "$(INSTALL_PATH)\$(BINARY)$(BINARY_EXT)"
else
	rm -f $(INSTALL_PATH)/$(BINARY)
endif
	@echo $(BINARY) uninstalled from $(INSTALL_PATH)

.DEFAULT_GOAL := help