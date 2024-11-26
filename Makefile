APP_NAME?=order-processor
SHELL := env APP_NAME=$(APP_NAME) $(SHELL)

BIN_DIR?=$(CURDIR)/bin

GOVERSION:=1.23

format-code: fmt goimports
.PHONY: format-code

fmt:
	@echo "Formatting code..."
	./scripts/style/fmt.sh
.PHONY: fmt

goimports:
	@echo "Formatting code..."
	./scripts/style/fix-imports.sh
.PHONY: goimports

vet:
	@echo "Vetting code..."
	@go vet ./...
	@echo "Done"
.PHONY: vet

test:
	@echo "Running tests..."
	@go test -v ./...
	@echo "Done"
.PHONY: test

build:
	@echo "Building..."
	@./scripts/build/app.sh
	@echo "Done"
.PHONY: build

run:
	@echo "Running..."
	@${BIN_DIR}/$(APP_NAME)
	@echo "Done"
.PHONY: run

vendor:
	@echo "Vendoring..."
	@go mod tidy && go mod vendor
	@echo "Done"
.PHONY: vendor

## Issue new release.
new-version: vet test build
	./scripts/release/new-version.sh
.PHONY: new-release

## Bump go version
bump-go-version:
	./scripts/bump-go.sh $(GOVERSION)
.PHONY: bump-go-version

