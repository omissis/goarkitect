_PROJECT_DIRECTORY = $(dir $(realpath $(firstword $(MAKEFILE_LIST))))
_GOLANG_IMAGE = golang:1.19.0
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

.PHONY: fmt fumpt imports

fmt:
	@find . -name "*.go" -type f -not -path '*/vendor/*' \
	| sed 's/^\.\///g' \
	| xargs -I {} sh -c 'echo "formatting {}.." && gofmt -w -s {}'

fumpt:
	@find . -name "*.go" -type f -not -path '*/vendor/*' \
	| sed 's/^\.\///g' \
	| xargs -I {} sh -c 'echo "formatting {}.." && gofumpt -w -extra {}'

imports:
	@goimports -w -e -local github.com/omissis main.go
	@goimports -w -e -local github.com/omissis cmd/
	@goimports -w -e -local github.com/omissis internal/

.PHONY: lint lint-go

lint: lint-go

lint-go:
	@golangci-lint -v run --color=always --config=${_PROJECT_DIRECTORY}/.rules/.golangci.yml ./...

.PHONY: test

test:
	@go test ./...

.PHONY: release-local release

release-local:
	@GO_VERSION=$$(go version | cut -d ' ' -f 3) OS_ARCH=$$(go version | cut -d ' ' -f 4) goreleaser release --debug --snapshot --rm-dist

release:
	@GO_VERSION=$$(go version | cut -d ' ' -f 3) OS_ARCH=$$(go version | cut -d ' ' -f 4) goreleaser --debug release --rm-dist

# Helpers

%-docker:
	$(call run-docker,${_GOLANG_IMAGE},make $*)

check-variable-%: # detection of undefined variables.
	@[[ "${${*}}" ]] || (echo '*** Please define variable `${*}` ***' && exit 1)
