SHELL:=/bin/bash
PORJ=adc-genius
ORG_PATH=code.htres.cn/casicloud
REPO_PATH=$(ORG_PATH)/$(PORJ)
export PATH := $(PWD)/bin:$(PATH)
export GOBIN=$(PWD)/bin
# build number
BN=$(shell ./scripts/gen-bn.sh)
VERSION=$(shell cat VERSION).$(BN)
LD_FLAGS="-w -X $(REPO_PATH)/version.Version=$(VERSION)"
SRCS := $(shell find . -name '*.go'| grep -v vendor)

clean:
	@rm -rf bin/

.PHONY: version
version:
	@echo $(VERSION)

test:
	@go test -v ./...

testrace:
	@go test -v --race ./...

.PHONY: lint
lint: 
	@for file in $(SRCS); do \
		golint $${file}; \
		if [ -n "$$(golint $${file})" ]; then \
			exit 1; \
		fi; \
	done