SHELL := /bin/bash

BINARY_NAME := beamstack
VERSION ?= 1.0.0
MAINTAINER := kennedy
ARCHITECTURE := amd64

GOCMD := go
GOBUILD := $(GOCMD) build
SOURCE_DIR := ../src

.PHONY: go-build
go-build:
	@echo "running go build ..."
	@pushd $(SOURCE_DIR); $(GOBUILD) -o ../package/$(BINARY_NAME) -v; popd

.PHONY: clean
clean:
	@echo "🧹 Cleaning up package files"
	@rm -rf $(BINARY_NAME)-$(VERSION)-$(ARCHITECTURE)
	@rm -f $(BINARY_NAME)
	@rm -f $(BINARY_NAME)-$(VERSION)-$(ARCHITECTURE).deb
	@echo

.PHONY: package
package: go-build
	@mkdir -p $(BINARY_NAME)-$(VERSION)-$(ARCHITECTURE)/DEBIAN
	@mkdir -p $(BINARY_NAME)-$(VERSION)-$(ARCHITECTURE)/usr/local/bin

	@echo "#!/bin/bash" > $(BINARY_NAME)-$(VERSION)-$(ARCHITECTURE)/DEBIAN/postinst
	@echo "set -e" >> $(BINARY_NAME)-$(VERSION)-$(ARCHITECTURE)/DEBIAN/postinst
	@echo 'USER=$$SUDO_USER' >> $(BINARY_NAME)-$(VERSION)-$(ARCHITECTURE)/DEBIAN/postinst
	
	@echo 'echo "Running post installation script for user: $$USER ..." | tee /tmp/postinst.log' >> $(BINARY_NAME)-$(VERSION)-$(ARCHITECTURE)/DEBIAN/postinst
	@echo 'echo "Creating directories ..." | tee -a /tmp/postinst.log' >> $(BINARY_NAME)-$(VERSION)-$(ARCHITECTURE)/DEBIAN/postinst
	
	@echo "mkdir -p /home/\$$USER/.$(BINARY_NAME)/config" >> $(BINARY_NAME)-$(VERSION)-$(ARCHITECTURE)/DEBIAN/postinst
	@echo "mkdir -p /home/\$$USER/.$(BINARY_NAME)/profiles" >> $(BINARY_NAME)-$(VERSION)-$(ARCHITECTURE)/DEBIAN/postinst

	@echo 'echo "Directories created successfully" | tee -a /tmp/postinst.log' >> $(BINARY_NAME)-$(VERSION)-$(ARCHITECTURE)/DEBIAN/postinst
	@echo 'echo "Granting necessary permissions to config files" | tee -a /tmp/postinst.log' >> $(BINARY_NAME)-$(VERSION)-$(ARCHITECTURE)/DEBIAN/postinst
	
	@echo "chmod 755 /home/\$$USER/.$(BINARY_NAME)" >> $(BINARY_NAME)-$(VERSION)-$(ARCHITECTURE)/DEBIAN/postinst
	@echo "chmod 755 /home/\$$USER/.$(BINARY_NAME)/config" >> $(BINARY_NAME)-$(VERSION)-$(ARCHITECTURE)/DEBIAN/postinst
	@echo "chmod 755 /home/\$$USER/.$(BINARY_NAME)/profiles" >> $(BINARY_NAME)-$(VERSION)-$(ARCHITECTURE)/DEBIAN/postinst
	
	@echo 'echo "$(BINARY_NAME) ($(ARCHITECTURE)-$(VERSION)) installation completed successfully" | tee -a /tmp/postinst.log' >> $(BINARY_NAME)-$(VERSION)-$(ARCHITECTURE)/DEBIAN/postinst
	@chmod 755 $(BINARY_NAME)-$(VERSION)-$(ARCHITECTURE)/DEBIAN/postinst

	@echo "#!/bin/bash" > $(BINARY_NAME)-$(VERSION)-$(ARCHITECTURE)/DEBIAN/postrm
	@echo "set -e" >> $(BINARY_NAME)-$(VERSION)-$(ARCHITECTURE)/DEBIAN/postrm
	@echo 'USER=$$SUDO_USER' >> $(BINARY_NAME)-$(VERSION)-$(ARCHITECTURE)/DEBIAN/postrm

	@echo 'echo "Uninstalling $(BINARY_NAME) ($(ARCHITECTURE)-$(VERSION)) for user: $$USER ..." | tee /tmp/postrm.log' >> $(BINARY_NAME)-$(VERSION)-$(ARCHITECTURE)/DEBIAN/postrm
	@echo "rm -rf /home/\$$USER/.$(BINARY_NAME) | tee -a /tmp/postrm.log" >> $(BINARY_NAME)-$(VERSION)-$(ARCHITECTURE)/DEBIAN/postrm
	
	@echo 'echo "$(BINARY_NAME) ($(ARCHITECTURE)-$(VERSION)) successfully uninstalled for user: $$USER" | tee -a /tmp/postrm.log' >> $(BINARY_NAME)-$(VERSION)-$(ARCHITECTURE)/DEBIAN/postrm
	@chmod 755 $(BINARY_NAME)-$(VERSION)-$(ARCHITECTURE)/DEBIAN/postrm

	@cp $(BINARY_NAME) $(BINARY_NAME)-$(VERSION)-$(ARCHITECTURE)/usr/local/bin/

	@echo "Package: $(BINARY_NAME)" > $(BINARY_NAME)-$(VERSION)-$(ARCHITECTURE)/DEBIAN/control
	@echo "Version: $(VERSION)" >> $(BINARY_NAME)-$(VERSION)-$(ARCHITECTURE)/DEBIAN/control
	@echo "Section: base" >> $(BINARY_NAME)-$(VERSION)-$(ARCHITECTURE)/DEBIAN/control
	@echo "Priority: optional" >> $(BINARY_NAME)-$(VERSION)-$(ARCHITECTURE)/DEBIAN/control
	@echo "Architecture: $(ARCHITECTURE)" >> $(BINARY_NAME)-$(VERSION)-$(ARCHITECTURE)/DEBIAN/control
	@echo "Maintainer: $(MAINTAINER) <$(MAINTAINER)@mavencode.com>" >> $(BINARY_NAME)-$(VERSION)-$(ARCHITECTURE)/DEBIAN/control
	@echo "Description: $(BINARY_NAME) CLI Tool" >> $(BINARY_NAME)-$(VERSION)-$(ARCHITECTURE)/DEBIAN/control

	@dpkg-deb --build $(BINARY_NAME)-$(VERSION)-$(ARCHITECTURE)

	@mkdir -p ../releases/$(VERSION)
	@cp $(BINARY_NAME)-$(VERSION)-$(ARCHITECTURE).deb ../releases/$(VERSION)/

	@$(MAKE) clean