# A Self-Documenting Makefile: http://marmelab.com/blog/2016/02/29/auto-documented-makefile.html

.PHONY: test prepare-release release build
.DEFAULT_GOAL := help

VERSION := $(shell git describe --tags --always --dirty)

build: ## Build the development binaries
	mkdir -p build
	GOOS=linux GOARCH=amd64 go build -o build/chagen-linux-amd64 -ldflags "-X main.version=$(VERSION)" chagen.go
	GOOS=darwin GOARCH=amd64 go build -o build/chagen-darwin-amd64 -ldflags "-X main.version=$(VERSION)" chagen.go
	GOOS=windows GOARCH=amd64 go build -o build/chagen-windows-amd64 -ldflags "-X main.version=$(VERSION)" chagen.go

prepare-env: ## Prepare the development/test environment
	go get -u github.com/alecthomas/gometalinter
	gometalinter --install

test: ## Run the tests
	go test $$(go list ./... | grep -v /vendor/)
	gometalinter --enable-all --line-length=100 -s vendor ./...

prepare-release: ## Prepare a new release
ifndef NEW_VERSION
	@echo "Usage: make release NEW_VERSION=0.1.2"
	exit 1
endif
	sed -i "s/var version =.*/var version = \"${NEW_VERSION}\"/" chagen.go

release: ## Build a new release
	mkdir -p releases/$(VERSION)
	GOOS=linux GOARCH=amd64 go build -o releases/$(VERSION)/chagen-linux-amd64 -ldflags "-X main.version=$(VERSION)" chagen.go
	GOOS=darwin GOARCH=amd64 go build -o releases/$(VERSION)/chagen-darwin-amd64 -ldflags "-X main.version=$(VERSION)" chagen.go
	GOOS=windows GOARCH=amd64 go build -o releases/$(VERSION)/chagen-windows-amd64 -ldflags "-X main.version=$(VERSION)" chagen.go

clean: ## Cleanup the builds
	rm -rf build releases

help:
	@grep -E '^[a-zA-Z0-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'