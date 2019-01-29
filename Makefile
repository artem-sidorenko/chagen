# A Self-Documenting Makefile: http://marmelab.com/blog/2016/02/29/auto-documented-makefile.html

.PHONY: test prepare-release release build
.DEFAULT_GOAL := help

VERSION := $(shell git describe --tags --always --dirty)

build: ## Build the development binaries
	mkdir -p build
	GOOS=linux GOARCH=amd64 go build -o build/chagen-linux-amd64 -ldflags "-X github.com/artem-sidorenko/chagen/internal/info.version=$(VERSION)" chagen.go
	GOOS=darwin GOARCH=amd64 go build -o build/chagen-darwin-amd64 -ldflags "-X github.com/artem-sidorenko/chagen/internal/info.version=$(VERSION)" chagen.go
	GOOS=windows GOARCH=amd64 go build -o build/chagen-windows-amd64 -ldflags "-X github.com/artem-sidorenko/chagen/internal/info.version=$(VERSION)" chagen.go

prepare-env: ## Prepare the development/test environment
	go get -u github.com/alecthomas/gometalinter
	gometalinter --install
	go get -u github.com/tcnksm/ghr

test: ## Run the tests
	go test -timeout 10s -race -cpu 4 -count 1 $$(go list ./... | grep -v /vendor/)
	gometalinter --enable-all --disable=dupl --deadline=300s --line-length=100 -s vendor ./...
# calculate the coverage only if everything was ok (see go test, coverage can have side effects)
	go test -timeout 10s -cover  $$(go list ./... | grep -v /vendor/)

prepare-release: ## Prepare a new release
ifndef NEW_VERSION
	@echo "Usage: make prepare-release NEW_VERSION=0.1.2"
	@exit 1
endif
	@if [ "$$(git symbolic-ref HEAD)" != "refs/heads/master" ]; then \
		echo "make prepare-release should be executed on the master branch" ;\
		exit 1 ;\
	fi
	chagen generate --github-owner artem-sidorenko --github-repo chagen -r v${NEW_VERSION} --github-release-url
	git add -u CHANGELOG.md chagen.go
	git commit -m "Release ${NEW_VERSION}"
	git tag -u 8B4B87B9 v${NEW_VERSION} -m "Release v${NEW_VERSION}"
	git push
	git push origin refs/tags/v${NEW_VERSION}

release: ## Build a new release
	rm -rf release/$(VERSION)
	mkdir -p release/$(VERSION)
	GOOS=linux GOARCH=amd64 go build -o release/$(VERSION)/chagen -ldflags "-X github.com/artem-sidorenko/chagen/internal/info.version=$(VERSION)" chagen.go
	tar cfzC release/$(VERSION)/chagen_$(VERSION)_Linux-64bit.tgz release/$(VERSION) chagen
	GOOS=darwin GOARCH=amd64 go build -o release/$(VERSION)/chagen -ldflags "-X github.com/artem-sidorenko/chagen/internal/info.version=$(VERSION)" chagen.go
	tar cfzC release/$(VERSION)/chagen_$(VERSION)_MacOS-64bit.tgz release/$(VERSION) chagen
	GOOS=windows GOARCH=amd64 go build -o release/$(VERSION)/chagen -ldflags "-X github.com/artem-sidorenko/chagen/internal/info.version=$(VERSION)" chagen.go
	zip -FS -j release/$(VERSION)/chagen_$(VERSION)_Windows-64bit.zip release/$(VERSION)/chagen
	rm release/$(VERSION)/chagen
	cd release/$(VERSION); sha256sum * > chagen_$(VERSION)_checksums.sha256

sign-release: ## Sign the checksums of release
	gpg -a -u 8B4B87B9 --detach-sig --output release/$(VERSION)/chagen_$(VERSION)_checksums.sha256.sig release/$(VERSION)/chagen_$(VERSION)_checksums.sha256

upload-release: ## Upload the new release builds to GitHub releases
	@tag=$$(git describe --tags --exact-match);\
	ghr -u artem-sidorenko -r chagen $$tag release/$$tag

clean: ## Cleanup the builds
	rm -rf build release

help:
	@grep -E '^[a-zA-Z0-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
