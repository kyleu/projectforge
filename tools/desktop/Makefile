.PHONY: build
build: ## Build all binaries
	go build -gcflags "all=-N -l" -o ../../build/debug/projectforge-desktop .

.PHONY: build-release
build-release: ## Build all binaries without debug information, clean up after
	go build -ldflags '-s -w' -trimpath -o ../../build/release/projectforge-desktop .

.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.DEFAULT_GOAL := help
