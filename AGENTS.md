# AGENTS.md

This file provides guidance to autonomus agents for working with code in this repository.

## Build Commands
- Build: `make build`
- Run: `make build`, then `./build/debug/projectforge`
- Run and watch: `make dev` or `./bin/dev.sh`
- Lint: `make lint` or `./bin/check.sh`
- Format: `./bin/format.sh`
- Test all: `./bin/test.sh` or `gotestsum -- ./app/...`
- Test single: `go test ./path/to/package -run TestName`
- Test with clean cache: `./bin/test.sh -c`
- Test and watch: `./bin/test.sh -w`

## Code Style Guidelines
- Use `gofumpt` for formatting (enforced by linters)
- Import order: standard libs, third-party libs, project imports
- Max line length: 160 characters
- Error handling: check all errors, use appropriate context
- Functions: max 100 lines recommended
- Cyclomatic complexity: max 30 recommended
- Naming: PascalCase for exported, camelCase for internal
- Comments: All exported items must be documented
- Use quicktemplate for HTML templates
- Follow existing patterns when adding new code