.PHONY: clean
clean: ## Removes builds and compiled templates
	@rm -rf tmp/*.hashcode
	@rm -rf out{{{ if .HasModule "database" }}}
	@find ./queries -type f -name '*.sql.go' -exec rm {} +{{{ end }}}
	@find ./views -type f -name '*.html.go' -exec rm {} +

.PHONY: dev
dev: ## Start the project, reloading on changes
	@bin/dev.sh

.PHONY: templates
templates:
	@bin/templates.sh

.PHONY: build
build: templates ## Build all binaries{{{ if .HasModule "wasmclient" }}}
	@GOOS=js GOARCH=wasm {{{ .GoBinarySafe }}} build -o ./assets/wasm/{{{ .Exec }}}.wasm ./app/wasm/...{{{ end }}}
	@{{{ .GoBinarySafe }}} build -gcflags "all=-N -l" -o build/debug/{{{ .Exec }}} .

.PHONY: build-release
build-release: templates ## Build all binaries without debug information, clean up after{{{ if .HasModule "wasmclient" }}}
	@GOOS=js GOARCH=wasm {{{ .GoBinarySafe }}} build -o ./assets/wasm/{{{ .Exec }}}.wasm ./app/wasm/...{{{ end }}}
	@{{{ .GoBinarySafe }}} build -ldflags '-s -w' -trimpath -o build/release/{{{ .Exec }}} .

.PHONY: lint
lint: ## Run linter
	@bin/check.sh

.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.DEFAULT_GOAL := help
