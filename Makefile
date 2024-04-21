include config.mk

.PHONY: docker-compose-up docker-compose-down setup-go-work setup-redis-conf

GO111MODULE ?= on
export GO111MODULE

MODULES ?= $(shell find $(PWD) -name "go.mod" | grep -v ".bingo" | xargs dirname)

COMMON_DIR=common

# Go versions may look like goM, goM.N, or goM.N.P. Only 1.22.* is supported.
supported_go_minor_version = go1.22
system_go_version := $(shell go version | sed 's/.*\(go[[:digit:]][[:digit:].]*\).*/\1/')
.PHONY: go-version-check
go-version-check:
	@: $(if $(findstring $(supported_go_minor_version), $(system_go_version)), \
				, \
				$(error go version $(system_go_version) not supported. Must use $(supported_go_minor_version).x))

setup-go-work: ## Sets up your go.work file
ifneq ($(IGNORE_GO_WORK_IF_EXISTS),true)
	@echo "Creating a go.work file"
	rm -f go.work
	$(GO) work init
	$(GO) work use -r .
endif

setup-redis-conf: # run this for loop to create each redis container configuration file
	@echo "Creating redis conf for all nodes: $(REDIES_NODES)"
	@rm -rf $(PWD)/infra/redis
	./scripts/gen-redis-conf.sh

lint: golangci-lint ## Runs golangci-lint with predefined configuration
	@echo "linting all of the Go files"
	golangci-lint version
	for dir in $(MODULES) ; do \
		cd $${dir} && golangci-lint run -v; \
	done

dep: wire
	@echo "generating dependency injection"
	cd purchase-svc && wire ./di
	cd product-svc && wire ./di

proto:
	cd $(COMMON_DIR) && $(MAKE) proto

docker-compose-up: setup-redis-conf
	docker compose up -d --remove-orphans

docker-compose-down:
	docker compose down -v
