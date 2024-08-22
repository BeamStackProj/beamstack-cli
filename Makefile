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

.PHONY: uninstall
uninstall:
	@echo "🧹 Uninstalling !"
	@sudo rm -f /usr/local/bin/$(BINARY_NAME)
	@sudo rm -rf ~/.$(BINARY_NAME)
	@echo "✅ $(BINARY_NAME) Uninstalled..!"


.PHONY: dryinstall
dryinstall: go-build
	@sudo mv $(SOURCE_DIR)/$(BINARY_NAME) /usr/local/bin/$(BINARY_NAME)
	@echo "✅ $(BINARY_NAME) installed..!"
	@echo
	make clean

.PHONY: install
install: uninstall go-build
	@sudo mkdir -p ~/.$(BINARY_NAME)
	@sudo mkdir -p ~/.$(BINARY_NAME)/config
	@sudo mkdir -p ~/.$(BINARY_NAME)/profiles
	@sudo cp ./tests/config.json ~/.$(BINARY_NAME)/config/config.json
	@sudo mv $(SOURCE_DIR)/$(BINARY_NAME) /usr/local/bin/$(BINARY_NAME)
	@sudo chmod -R 777 ~/.$(BINARY_NAME)
	@echo "✅ $(BINARY_NAME) installed..!"
	@echo
	@$(MAKE) clean
