SHELL := /bin/bash
CWD:=$(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
SOURCE_DIR=./src
BINARY_NAME=beamflow

.PHONY: go-build
go-build:
	@echo "running go build ..."
	@cd $(SOURCE_DIR) && $(GOBUILD) -o $(BINARY_NAME) -v
	@sudo mv $(SOURCE_DIR)/$(BINARY_NAME) /usr/local/bin/$(BINARY_NAME)
	@echo

.PHONY: clean
clean:
	@echo "🧹 Cleaning up binary file"
	@rm -f $(SOURCE_DIR)/$(BINARY_NAME)
	@echo

.PHONY: install
install: go-build clean
	@echo "✅ $(BINARY_NAME) installed..!"

