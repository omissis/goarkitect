_PROJECT_DIRECTORY = $(dir $(realpath $(firstword $(MAKEFILE_LIST))))
_GOLANG_IMAGE = golang:1.19.1
_PROJECTNAME = goarkitect
_GOARCH = "amd64"

NETRC_FILE ?= ~/.netrc

ifeq ("$(shell uname -m)", "arm64")
	_GOARCH = "arm64"
endif

#1: docker image
#2: make target
define run-docker
	@docker run --rm \
		-e CGO_ENABLED=0 \
		-e GOARCH=${_GOARCH} \
		-e GOOS=linux \
		-w /app \
		-v ${NETRC_FILE}:/root/.netrc \
		-v ${_PROJECT_DIRECTORY}:/app \
		$(1) $(2)
endef

.PHONY: env

env:
	@echo 'export CGO_ENABLED=0'
	@echo 'export GOARCH=${_GOARCH}'
	@grep -v '^#' .env | sed 's/^/export /'

.PHONY: mod-download mod-tidy mod-verify

mod-download:
	@go mod download

mod-tidy:
	@go mod tidy

mod-verify:
	@go mod verify

.PHONY: mod-check-upgrades mod-upgrade

mod-check-upgrades:
	@go list -mod=readonly -u -f "{{if (and (not (or .Main .Indirect)) .Update)}}{{.Path}}: {{.Version}} -> {{.Update.Version}}{{end}}" -m all

mod-upgrade:
	@go get -u ./... && go mod tidy

.PHONY: fmt fumpt imports gci

fmt:
	@find . -name "*.go" -type f -not -path '*/vendor/*' \
	| sed 's/^\.\///g' \
	| xargs -I {} sh -c 'echo "formatting {}.." && gofmt -w -s {}'

fumpt:
	@find . -name "*.go" -type f -not -path '*/vendor/*' \
	| sed 's/^\.\///g' \
	| xargs -I {} sh -c 'echo "formatting {}.." && gofumpt -w -extra {}'

imports:
	@goimports -v -w -e -local github.com/omissis main.go
	@goimports -v -w -e -local github.com/omissis cmd/
	@goimports -v -w -e -local github.com/omissis internal/

gci:
	@find . -name "*.go" -type f -not -path '*/vendor/*' \
	| sed 's/^\.\///g' \
	| xargs -I {} sh -c 'echo "formatting imports for {}.." && \
	gci write --skip-generated -s standard,default,"prefix(github.com/omissis)" {}'

.PHONY: lint lint-go

lint: lint-go

lint-go:
	@golangci-lint -v run --color=always --config=${_PROJECT_DIRECTORY}/.rules/.golangci.yml ./...

.PHONY: test

test:
	@go test -v -race -covermode=atomic -coverprofile=coverage.out ./...

.PHONY: quickbuild

quickbuild:
	@go build -o /dev/null main.go

.PHONY: fix qa

fix: mod-tidy fmt fumpt imports gci

qa: mod-verify lint test quickbuild

.PHONY: examples

examples:
	@echo "== VALIDATE | text | config file generating violations ===========================\n"
	@go run main.go validate examples/.goarkitect.violations.yaml
	@echo "\n================================================================================\n"

	@echo "== VALIDATE | text | config file with no violations ==============================\n"
	@go run main.go validate examples/.goarkitect.yaml
	@echo "\n================================================================================\n"

	@echo "== VALIDATE | json | config file generating violations ===========================\n"
	@go run main.go validate --output=json examples/.goarkitect.violations.yaml
	@echo "\n================================================================================\n"

	@echo "== VALIDATE | json | config file with no violations ==============================\n"
	@go run main.go validate --output=json examples/.goarkitect.yaml
	@echo "\n================================================================================\n"

	@echo "== VERIFY | text | invalid config file =========================================\n"
	@go run main.go verify examples/.goarkitect.invalid.yaml
	@echo "\n================================================================================\n"

	@echo "== VERIFY | text | config file generating violations ===========================\n"
	@go run main.go verify examples/.goarkitect.violations.yaml
	@echo "\n================================================================================\n"

	@echo "== VERIFY | text | config file with no violations ==============================\n"
	@go run main.go verify examples/.goarkitect.yaml
	@echo "\n================================================================================\n"

	@echo "== VERIFY | json | invalid config file =========================================\n"
	@go run main.go verify --output=json examples/.goarkitect.invalid.yaml
	@echo "\n================================================================================\n"

	@echo "== VERIFY | json | config file generating violations ===========================\n"
	@go run main.go verify --output=json examples/.goarkitect.violations.yaml
	@echo "\n================================================================================\n"

	@echo "== VERIFY | json | config file with no violations ==============================\n"
	@go run main.go verify --output=json examples/.goarkitect.yaml
	@echo "\n================================================================================\n"

.PHONY: build release

build:
	@export GO_VERSION=$$(go version | cut -d ' ' -f 3) && \
	goreleaser check && \
	goreleaser release --debug --snapshot --rm-dist

release:
	@export GO_VERSION=$$(go version | cut -d ' ' -f 3) && \
	goreleaser check && \
	goreleaser --debug release --rm-dist

# Helpers

%-docker:
	$(call run-docker,${_GOLANG_IMAGE},make $*)

check-variable-%: # detection of undefined variables.
	@[[ "${${*}}" ]] || (echo '*** Please define variable `${*}` ***' && exit 1)
