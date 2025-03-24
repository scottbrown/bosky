.PHONY: help deps clean build dist
.DEFAULT_GOAL := help

GO_BIN = go

PKG = github.com/scottbrown/bosky
ARTIFACT = bosky

ARTIFACT_DIR := .build

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

deps: ## ensure all dependencies are downloaded
	@hash $(GO_BIN) > /dev/null 2>&1 || (echo "Install Go to continue"; exit 1)
	@mkdir -p $(ARTIFACT_DIR)

clean: ## removes all derived files
	rm -rf $(ARTIFACT_DIR)

build: deps ## builds a local version of the application
	$(GO_BIN) fmt
	$(eval GIT_HASH := $(shell git rev-parse --short HEAD))
	$(GO_BIN) build -ldflags "-X $(PKG).VERSION=$(GIT_HASH)" -o $(ARTIFACT_DIR)/$(ARTIFACT) $(PKG)/cmd

test:
	go test ./...

dist: deps ## builds distributable versions of the app in all OS/ARCH combos
ifdef VERSION
	$(eval BUILD_VERSION := $(VERSION))
else
	$(eval BUILD_VERSION := $(shell git rev-parse --short HEAD))
endif
	mkdir -p $(ARTIFACT_DIR)/linux-amd64
	GOOS=linux GOARCH=amd64 $(GO_BIN) build -ldflags "-X $(PKG).VERSION=$(BUILD_VERSION)" -o $(ARTIFACT_DIR)/linux-amd64/$(ARTIFACT) $(PKG)/cmd
	mkdir -p $(ARTIFACT_DIR)/linux-arm
	GOOS=linux GOARCH=arm $(GO_BIN) build -ldflags "-X $(PKG).VERSION=$(BUILD_VERSION)" -o $(ARTIFACT_DIR)/linux-arm/$(ARTIFACT) $(PKG)/cmd
	mkdir -p $(ARTIFACT_DIR)/linux-i386
	GOOS=linux GOARCH=386 $(GO_BIN) build -ldflags "-X $(PKG).VERSION=$(BUILD_VERSION)" -o $(ARTIFACT_DIR)/linux-i386/$(ARTIFACT) $(PKG)/cmd
	mkdir -p $(ARTIFACT_DIR)/darwin-amd64
	GOOS=darwin GOARCH=amd64 $(GO_BIN) build -ldflags "-X $(PKG).VERSION=$(BUILD_VERSION)" -o $(ARTIFACT_DIR)/darwin-amd64/$(ARTIFACT) $(PKG)/cmd
	mkdir -p $(ARTIFACT_DIR)/windows-amd64
	GOOS=windows GOARCH=amd64 $(GO_BIN) build -ldflags "-X $(PKG).VERSION=$(BUILD_VERSION)" -o $(ARTIFACT_DIR)/windows-amd64/$(ARTIFACT).exe $(PKG)/cmd
	mkdir -p $(ARTIFACT_DIR)/windows-i386
	GOOS=windows GOARCH=386 $(GO_BIN) build -ldflags "-X $(PKG).VERSION=$(BUILD_VERSION)" -o $(ARTIFACT_DIR)/windows-i386/$(ARTIFACT).exe $(PKG)/cmd

release: deps ## Creates releaseable artifacts ready for public download
ifndef VERSION
	$(error "Specify a VERSION to continue.")
endif
	$(MAKE) dist VERSION=$(VERSION)
	tar cfz $(ARTIFACT_DIR)/$(ARTIFACT)_$(VERSION)_linux_amd64.tar.gz -C $(ARTIFACT_DIR)/linux-amd64 bosky
	tar cfz $(ARTIFACT_DIR)/$(ARTIFACT)_$(VERSION)_linux_arm.tar.gz -C $(ARTIFACT_DIR)/linux-arm bosky
	tar cfz $(ARTIFACT_DIR)/$(ARTIFACT)_$(VERSION)_linux_i386.tar.gz -C $(ARTIFACT_DIR)/linux-i386 bosky
	tar cfz $(ARTIFACT_DIR)/$(ARTIFACT)_$(VERSION)_darwin_amd64.tar.gz -C $(ARTIFACT_DIR)/darwin-amd64 bosky
	tar cfz $(ARTIFACT_DIR)/$(ARTIFACT)_$(VERSION)_windows_amd64.tar.gz -C $(ARTIFACT_DIR)/windows-amd64 bosky.exe
	tar cfz $(ARTIFACT_DIR)/$(ARTIFACT)_$(VERSION)_windows_i386.tar.gz -C $(ARTIFACT_DIR)/windows-i386 bosky.exe
