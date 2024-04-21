GOPATH ?= $(shell go env GOPATH)
GOBIN  ?= $(firstword $(subst :, ,${GOPATH}))/bin
GO     ?= $(shell which go)
GOTEST  = $(GO) test
GOBUILD = $(GO) build
GOCLEAN = $(GO) clean
GOINSTALL = $(GO) install

GOLANGCI := $(shell command -v golangci-lint || echo "$(GOBIN)/golangci-lint")
golangci-lint: VERSION := 1.57.1
golangci-lint:
	@echo "(re)installing golangci-lint-v$(VERSION)"
	$(GOINSTALL) github.com/golangci/golangci-lint/cmd/golangci-lint@v$(VERSION)
	@echo "done."

WIRE := $(shell command -v wire || echo "$(GOBIN)/wire")
wire: VERSION := 0.6.0
wire:
	@echo "(re)installing wire-v$(VERSION)"
	$(GOINSTALL) github.com/google/wire/cmd/wire@v$(VERSION)
	@echo "done."