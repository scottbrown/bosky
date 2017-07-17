.PHONY: help deps clean build dist

PKG = github.com/scottbrown/bosky
ARTIFACT = bosky

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

deps: ## ensure all dependencies are downloaded
	go get -u gopkg.in/urfave/cli.v1
	go get -u github.com/aws/aws-sdk-go

clean: ## removes all derived files
	rm -f bosky
	rm -rf $(GOPATH)/dist/bosky

build: ## builds a local version of the application
	go build -o $(ARTIFACT) $(PKG)

dist: ## builds distributable versions of the app in all OS/ARCH combos
	mkdir -p $(GOPATH)/dist/$(ARTIFACT)
	GOOS=linux GOARCH=amd64 go build -o $(GOPATH)/dist/$(ARTIFACT)/$(ARTIFACT).linux-amd64 $(PKG)
	GOOS=linux GOARCH=arm go build -o $(GOPATH)/dist/$(ARTIFACT)/$(ARTIFACT).linux-arm $(PKG)
	GOOS=linux GOARCH=386 go build -o $(GOPATH)/dist/$(ARTIFACT)/$(ARTIFACT).linux-386 $(PKG)
	GOOS=darwin GOARCH=amd64 go build -o $(GOPATH)/dist/$(ARTIFACT)/$(ARTIFACT).darwin-amd64 $(PKG)
	GOOS=windows GOARCH=amd64 go build -o $(GOPATH)/dist/bosky/$(ARTIFACT).windows-amd64.exe $(PKG)
	GOOS=windows GOARCH=386 go build -o $(GOPATH)/dist/bosky/$(ARTIFACT).windows-386.exe github.com/scottbrown/bosky

