.PHONY: clean
clean: ## Removes builds and compiled templates
	@rm -rf tmp/*.hashcode
	@rm -rf out
	@find ./views -type f -name '*.html.go' -exec rm {} +

.PHONY: dev
dev: ## Start the project, reloading on changes
	@bin/dev.sh

.PHONY: templates
templates:
	@bin/templates.sh

.PHONY: build
build: export GOEXPERIMENT=jsonv2
build: templates ## Build all binaries
	@go build -gcflags "all=-N -l" -o build/debug/projectforge .

.PHONY: build-verbose
build-verbose: export GOEXPERIMENT=jsonv2
build-verbose: templates ## Build all binaries
	@go build -v -x -gcflags "all=-N -l" -o build/debug/projectforge .

.PHONY: build-release
build-release: export GOEXPERIMENT=jsonv2
build-release: templates ## Build all binaries without debug information, clean up after
	@go build -ldflags '-s -w' -trimpath -o build/release/projectforge .

.PHONY: lint
lint: ## Run linter
	@bin/check.sh

.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.DEFAULT_GOAL := help
