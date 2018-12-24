NAME     := lsas
VERSION  := $(shell git describe --abbrev=0)
REVISION := $(shell git rev-parse --short HEAD)
.DEFAULT_GOAL := help

DIST_DIRS := find * -type d -exec

.PHONY: pre-dep
pre-dep:  ## Prepare tools for resolve dependencies
ifeq ($(shell command -v dep 2> /dev/null),)
	go get -u github.com/golang/dep/cmd/dep
endif

.PHONY: dep
dep: pre-dep ## Resolve dependencies
	dep ensure

.PHONY: build
build: dep ## Build binary
	go build -o ./bin/lsas ./cmd

.PHONY: clean
clean: ## Cleanup destination directories
	@rm -rf bin/*
	@rm -rf dist/*

.PHONY: test
test:  ## Execute tests
	go test ./...

.PHONY: pre-dist
pre-dist:  ## Prepare tools for release
ifeq ($(shell command -v goxc 2> /dev/null),)
	go get -v -u github.com/laher/goxc
endif
ifeq ($(shell command -v ghr 2> /dev/null),)
	go get -v -u github.com/tcnksm/ghr
endif

.PHONY: dist
dist: dep pre-dist ## Build release objects
	goxc
	openssl dgst -sha256 dist/snapshot/lsas_linux_386.zip
	openssl dgst -sha256 dist/snapshot/lsas_darwin_386.zip
	openssl dgst -sha256 dist/snapshot/lsas_darwin_amd64.zip
	openssl dgst -sha256 dist/snapshot/lsas_linux_amd64.zip

.PHONY: releases
releases: dist  ## Upload release object to GitHub
	git push origin --tags
	ghr ${VERSION} dist/snapshot

.PHONY: help
help: ## Show options
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'
