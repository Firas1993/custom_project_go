MAKEFLAGS    += --always-make --warn-undefined-variables
SHELL        := /usr/bin/env bash
.SHELLFLAGS  := -e -o pipefail -c
.NOTPARALLEL :

PROJECT_ROOT ?= $(shell git rev-parse --show-toplevel)
PROJECT_NAME ?= irp-app-from-template

export GOPRIVATE ?= github.com/mime-rona,github.com/irp

format:
	pre-commit run --all-files golangci-lint-fmt

lint:
	pre-commit run --all-files

checkauth:
	if [[ -z "$${CI:-}" ]] && ! gcloud auth print-access-token >/dev/null; then gcloud auth login --update-adc; fi

test: checkauth
	@set -o pipefail; \
		go run github.com/dave/courtney -t "-json" -v -o .coverage "$${package:-./...}" 2>&1 \
		| grep -v "no packages being tested depend on matches for pattern" \
		| go run gotest.tools/gotestsum --raw-command --ignore-non-json-output-lines --format testname -- cat
	@sed -i.bak -e '/\.gen\.go:/d' .coverage
	@rm -f .coverage.bak
	@if grep -E '0$$' .coverage; then printf "\\e[31m[ERROR]\\e[0m   Missing test coverage\\n"; exit 1; fi;

mock.gen:
	@find . -type d -name mocks -print0 | xargs -0 --no-run-if-empty rm -rf
	@go run github.com/vektra/mockery/v3@v3.2.5

cov.render:
	go tool cover -html=.coverage

mod:
	go mod tidy

update:
	go get -u -t toolchain@1.26.0 ./...
	go mod tidy

upgrade: _cruft.update update

clean:
	go clean

# Private targets
_cruft.update:
	cruft update

# Imports
ifneq ("$(wildcard ${PROJECT_ROOT}/Dockerfile)", "")
  include ${PROJECT_ROOT}/.makefiles/Makefile.docker
endif
ifneq ("$(wildcard ${PROJECT_ROOT}/atlas.hcl)", "")
  include ${PROJECT_ROOT}/.makefiles/Makefile.migration
endif
ifneq ("$(wildcard ${PROJECT_ROOT}/gqlgen.yaml)", "")
  include ${PROJECT_ROOT}/.makefiles/Makefile.gql
endif
ifneq ("$(wildcard ${HOME}/.local/share/mime/Makefile)", "")
  include ${HOME}/.local/share/mime/Makefile
endif
