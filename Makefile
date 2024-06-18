SHELL := /bin/bash
CWD:=$(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
SOURCE_DIR=./src
BINARY_NAME=beamstack

.PHONY: go-build
go-build:
	@echo "running go build ..."
	@cd $(SOURCE_DIR) && $(GOBUILD) -o $(BINARY_NAME) -v

.PHONY: clean
clean:
	@echo "🧹 Cleaning up binary file"
	@rm -f $(SOURCE_DIR)/$(BINARY_NAME)
	@echo

.PHONY: install
install: go-build
	@sudo mkdir ~/.$(BINARY_NAME) -p
	@sudo mkdir ~/.$(BINARY_NAME)/config -p
	@sudo mkdir ~/.$(BINARY_NAME)/profiles -p
	@sudo cp ./tests/config.json ~/.$(BINARY_NAME)/config/config.json
	@sudo mv $(SOURCE_DIR)/$(BINARY_NAME) /usr/local/bin/$(BINARY_NAME)
	@chmod -R 755 ~/.$(BINARY_NAME)
	@echo "✅ $(BINARY_NAME) installed..!"
	@echo
	make clean

.PHONY: uninstall
uninstall:
	@echo "🧹 Uninstalling !"
	@sudo rm -f /usr/local/bin/$(BINARY_NAME)
	@sudo rm -rf ~/.$(BINARY_NAME)
	@echo "✅ $(BINARY_NAME) Uninstalled..!"


