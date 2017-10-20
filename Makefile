.PHONY: help deps clean build dist
.DEFAULT_GOAL := help

GO_BIN = go

PKG = github.com/scottbrown/bosky
ARTIFACT = bosky

ARTIFACT_DIR := $(GOPATH)/dist/$(ARTIFACT)

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

deps: ## ensure all dependencies are downloaded
	@hash $(GO_BIN) > /dev/null 2>&1 || (echo "Install Go to continue"; exit 1)
	@mkdir -p $(ARTIFACT_DIR)

clean: ## removes all derived files
	rm -f bosky
	rm -rf $(GOPATH)/dist/bosky

build: deps ## builds a local version of the application
	$(GO_BIN) fmt
	$(GO_BIN) build -o $(ARTIFACT) $(PKG)

dist: deps ## builds distributable versions of the app in all OS/ARCH combos
	mkdir -p $(ARTIFACT_DIR)/linux-amd64
	GOOS=linux GOARCH=amd64 $(GO_BIN) build -o $(ARTIFACT_DIR)/linux-amd64/$(ARTIFACT) $(PKG)
	mkdir -p $(ARTIFACT_DIR)/linux-arm
	GOOS=linux GOARCH=arm $(GO_BIN) build -o $(ARTIFACT_DIR)/linux-arm/$(ARTIFACT) $(PKG)
	mkdir -p $(ARTIFACT_DIR)/linux-i386
	GOOS=linux GOARCH=386 $(GO_BIN) build -o $(ARTIFACT_DIR)/linux-i386/$(ARTIFACT) $(PKG)
	mkdir -p $(ARTIFACT_DIR)/darwin-amd64
	GOOS=darwin GOARCH=amd64 $(GO_BIN) build -o $(ARTIFACT_DIR)/darwin-amd64/$(ARTIFACT) $(PKG)
	mkdir -p $(ARTIFACT_DIR)/windows-amd64
	GOOS=windows GOARCH=amd64 $(GO_BIN) build -o $(ARTIFACT_DIR)/windows-amd64/$(ARTIFACT).exe $(PKG)
	mkdir -p $(ARTIFACT_DIR)/windows-i386
	GOOS=windows GOARCH=386 $(GO_BIN) build -o $(ARTIFACT_DIR)/windows-i386/$(ARTIFACT).exe github.com/scottbrown/bosky

