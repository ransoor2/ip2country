export

LOCAL_BIN:=$(CURDIR)/bin
PATH:=$(LOCAL_BIN):$(PATH)
BIN_DIR := .tools/bin
GOLANGCI_LINT_VERSION := 1.55.2
GOLANGCI_LINT := $(BIN_DIR)/golangci-lint_$(GOLANGCI_LINT_VERSION)

# HELP =================================================================================================================
# This will output the help for each task
# thanks to https://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
.PHONY: help

help: ## Display this help screen
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

install:
	go install github.com/swaggo/swag/cmd/swag@v1.8.4

swag-v1: install ### swag init
	$(GOPATH)/bin/swag init -g internal/controller/http/v1/router.go
.PHONY: swag-v1

run: swag-v1 ### swag run
	go mod tidy && go mod download && \
	DISABLE_SWAGGER_HTTP_HANDLER='' GIN_MODE=debug CGO_ENABLED=0 go run -tags ip2country ./cmd/app
.PHONY: run

$(GOLANGCI_LINT):
	curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(BIN_DIR) v$(GOLANGCI_LINT_VERSION)
	mv $(BIN_DIR)/golangci-lint $(GOLANGCI_LINT)

lint: $(GOLANGCI_LINT)
	$(GOLANGCI_LINT) run --fast

test: ### run test
	go test -v -cover -race ./internal/... ./pkg/... ./test/...
.PHONY: test

docker-build: ### build docker image
	@docker build -t ip2country:latest .
.PHONY: docker-build

kind-install: docker-build ### create kind cluster and install everything
	./scripts/cluster_install.sh
.PHONY: kind-install

kind-delete:  ### delete kind cluster
	./scripts/cluster_delete.sh
.PHONY: kind-install