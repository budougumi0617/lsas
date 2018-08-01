NAME     := lsas
VERSION  := $(shell git describe --abbrev=0)
REVISION := $(shell git rev-parse --short HEAD)

DIST_DIRS := find * -type d -exec

.PHONY: pre-dep
pre-dep:
ifeq ($(shell command -v dep 2> /dev/null),)
	go get -u github.com/golang/dep/cmd/dep
endif

.PHONY: dep
dep: pre-dep
	dep ensure

.PHONY: clean
clean:
	rm -rf bin/*
	rm -rf dist/*

.PHONY: test
test:
	go test ./...

.PHONY: pre-dist
pre-dist:
ifeq ($(shell command -v goxc 2> /dev/null),)
	go get -v -u github.com/laher/goxc
endif
ifeq ($(shell command -v ghr 2> /dev/null),)
	go get -v -u github.com/tcnksm/ghr
endif

.PHONY: dist
dist: dep pre-dist
	goxc
	openssl dgst -sha256 dist/snapshot/lsas_linux_386.zip
	openssl dgst -sha256 dist/snapshot/lsas_darwin_386.zip
	openssl dgst -sha256 dist/snapshot/lsas_darwin_amd64.zip
	openssl dgst -sha256 dist/snapshot/lsas_linux_amd64.zip

.PHONY: releases
releases: dist
	git push origin --tags
	ghr ${VERSION} dist/snapshot
