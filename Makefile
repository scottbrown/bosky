.PHONY: help deps clean build dist
.DEFAULT_GOAL := help

GO_BIN = go

PKG = github.com/scottbrown/bosky
ARTIFACT = bosky

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

deps: ## ensure all dependencies are downloaded
	@hash $(GO_BIN) > /dev/null 2>&1 || (echo "Install Go to continue"; exit 1)
	@mkdir -p $(GOPATH)/dist/$(ARTIFACT)

clean: ## removes all derived files
	rm -f bosky
	rm -rf $(GOPATH)/dist/bosky

build: deps ## builds a local version of the application
	$(GO_BIN) fmt
	$(GO_BIN) build -o $(ARTIFACT) $(PKG)

dist: deps ## builds distributable versions of the app in all OS/ARCH combos
	GOOS=linux GOARCH=amd64 $(GO_BIN) build -o $(GOPATH)/dist/$(ARTIFACT)/$(ARTIFACT).linux-amd64 $(PKG)
	GOOS=linux GOARCH=arm $(GO_BIN) build -o $(GOPATH)/dist/$(ARTIFACT)/$(ARTIFACT).linux-arm $(PKG)
	GOOS=linux GOARCH=386 $(GO_BIN) build -o $(GOPATH)/dist/$(ARTIFACT)/$(ARTIFACT).linux-386 $(PKG)
	GOOS=darwin GOARCH=amd64 $(GO_BIN) build -o $(GOPATH)/dist/$(ARTIFACT)/$(ARTIFACT).darwin-amd64 $(PKG)
	GOOS=windows GOARCH=amd64 $(GO_BIN) build -o $(GOPATH)/dist/bosky/$(ARTIFACT).windows-amd64.exe $(PKG)
	GOOS=windows GOARCH=386 $(GO_BIN) build -o $(GOPATH)/dist/bosky/$(ARTIFACT).windows-386.exe github.com/scottbrown/bosky

