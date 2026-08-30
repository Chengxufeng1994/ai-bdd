.DEFAULT_GOAL := help

## help: List the available targets
.PHONY: help
help:
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'

# Code generation
#
# //go:generate directives live in the package that consumes their output, so a
# deleted package takes its directive with it and never leaves a path pointing
# at something that no longer exists. This Makefile is the entry point on top of
# them: `go generate -run` selects one generator at a time, which is what makes
# per-tool targets possible without moving the directives to one central file.
#
# A future generator fits the same shape. buf and protoc cannot share a single
# root generate.go — they carry their own config and plugin flags — but they can
# be selected the same way:
#
#	gen-proto:
#		go generate -run buf ./internal/...

# Directory names written by generators. verify-generated checksums whichever of
# these currently exist, so adding a generator means adding its output directory
# name here — otherwise its output silently stops being checked for staleness.
#
# Matched by name rather than by path because some of them do not exist yet:
# mocks/ directories appear next to the first interface mockery is pointed at.
GENERATED_DIR_NAMES := apigen mocks

## gen: Run every generator
#
# Go generators are declared in ./generate.go and run in the order they appear
# there. This target exists so that non-Go generators — a TypeScript client, an
# asset pipeline — have somewhere to join, and so there is one command to learn.
.PHONY: gen
gen:
	go generate ./...

## gen-api: Regenerate only the HTTP server from api/openapi.yaml
.PHONY: gen-api
gen-api:
	go generate -run oapi-codegen ./...

## gen-mock: Regenerate only the test doubles for the application layer's ports
.PHONY: gen-mock
gen-mock:
	go generate -run mockery ./...

# -race everywhere is not optional here. Scenarios run concurrently
# (Concurrency: 4), so the claim that per-scenario state is isolated is a claim
# the race detector is the only thing actually checking.

## test: Unit and acceptance tests, no external dependencies
.PHONY: test
test:
	go test -race ./...

## test-cover: Tests with a coverage summary
.PHONY: test-cover
test-cover:
	go test -race -cover ./...

## test-integration: Add the tests that need real adapters
.PHONY: test-integration
test-integration:
	go test -race -tags=integration ./...

## test-acceptance: Run the .feature files only
.PHONY: test-acceptance
test-acceptance:
	go test -race -count=1 ./test/acceptance/

## test-wip: Run only scenarios tagged @wip
.PHONY: test-wip
test-wip:
	go test -race -count=1 ./test/acceptance/ -godog.tags=@wip

## steps: List every registered step definition
.PHONY: steps
steps:
	go test ./test/acceptance/ -godog.definitions

# Tooling
#
# Tools are pinned in the `tool` block of go.mod, which replaces the pre-Go-1.24
# tools.go workaround entirely: `go tool` lists them, `go tool <name>` runs one
# without installing it, `go get -tool <pkg>@<version>` adds, `@none` removes,
# and `go get tool` upgrades every one of them at once.

## tools: Compile every pinned tool into bin/, so CI can cache them
.PHONY: tools
tools:
	go build -o bin/ tool

## tools-upgrade: Upgrade every pinned tool to its latest version
.PHONY: tools-upgrade
tools-upgrade:
	go get tool
	go mod tidy

## lint: Run golangci-lint
.PHONY: lint
lint:
	go tool golangci-lint run ./...

## lint-fix: Run golangci-lint and apply the fixes it can make itself
.PHONY: lint-fix
lint-fix:
	go tool golangci-lint run --fix ./...

## build: Build the server with the version stamped in
.PHONY: build
build:
	go build -ldflags "-X skeleton/pkg/version.value=$$(git describe --tags --always --dirty 2>/dev/null || echo dev)" -o bin/server ./cmd/server

## verify: Everything CI should run
.PHONY: verify
verify: verify-fmt verify-generated verify-stamp lint
	go vet ./...
	go test -race ./...

.PHONY: verify-fmt
verify-fmt:
	@test -z "$$(gofmt -l .)" || { gofmt -l .; echo "run: gofmt -w ."; exit 1; }

# Checksum around a regeneration rather than asking git whether anything
# changed: the generated files may be untracked, or the repository may have no
# commits yet, and in both cases `git diff` reports success no matter how stale
# the committed code actually is.
.PHONY: verify-generated
verify-generated:
	@sum() { find . -type d \( $(foreach n,$(GENERATED_DIR_NAMES),-name $(n) -o) -false \) \
		-exec find {} -type f \; 2>/dev/null | sort | xargs -r sha256sum | sha256sum; }; \
	before=$$(sum); \
	$(MAKE) --no-print-directory gen >/dev/null; \
	after=$$(sum); \
	if [ "$$before" != "$$after" ]; then \
		echo "generated code is stale — run 'make gen' and commit the result"; \
		exit 1; \
	fi

# `-X` aimed at a symbol that does not exist is not an error: the linker writes
# nothing and the binary reports its default. Renaming the package or the
# variable would therefore ship a build reporting "dev" with every test still
# green, because the acceptance suite injects its own version and never touches
# this path. This is the only gate that would notice.
#
# Like verify-fmt and verify-generated, this has no `## ` help line: `make help`
# lists what you would run by hand, and these three are reached through `verify`.
.PHONY: verify-stamp
verify-stamp:
	@command -v strings >/dev/null || { echo "verify-stamp needs 'strings' (binutils)"; exit 1; }
	@go build -ldflags "-X skeleton/pkg/version.value=STAMP-PROBE" -o bin/.stamp-probe ./cmd/server
	@strings bin/.stamp-probe | grep -qx STAMP-PROBE || { rm -f bin/.stamp-probe; echo "ldflags symbol path no longer resolves: skeleton/pkg/version.value"; exit 1; }
	@rm -f bin/.stamp-probe

## clean: Remove build output
.PHONY: clean
clean:
	rm -rf bin
