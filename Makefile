# Makefile

# Project-specific settings
BINARY_NAME=opper
MAIN_FILE=cmd/opper/main.go

# Go related variables.
GOBASE=$(shell pwd)
GOBIN=$(GOBASE)/bin
GOPATH=$(shell go env GOPATH)

# Go commands
GOCMD=go
GOBUILD=$(GOCMD) build
GOINSTALL=$(GOCMD) install
GOCLEAN=$(GOCMD) clean

# Build the project
build:
	$(GOBUILD) -o $(GOBIN)/$(BINARY_NAME) $(MAIN_FILE)

# Install dependencies and the project
install:
	$(GOINSTALL) $(MAIN_FILE)

# Clean build files
clean:
	$(GOCLEAN)
	rm -f $(GOBIN)/$(BINARY_NAME)

.PHONY: build install clean