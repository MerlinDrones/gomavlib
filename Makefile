SHELL := /bin/bash
# Project related variables
CWD= $(shell git rev-parse --show-toplevel)
APP_NAME = $(shell basename $(CWD))
PKG = $(shell echo $(CWD) | awk -F"/" '{ print $$(NF-2)"/"$$(NF-1)"/"$$NF }')

# Directories
WD := $(subst $(BSLASH),$(FSLASH),$(shell pwd))
MD := $(subst $(BSLASH),$(FSLASH),$(shell dirname "$(realpath $(lastword $(MAKEFILE_LIST)))"))
BUILD_DIR = $(WD)/build
PKG_DIR = $(MD)/pkg
CMD_DIR = $(PKG_DIR)/cmd
DIST_DIR = $(WD)/dist
LOG_DIR = $(WD)/log
CONFIG_DIR = $(WD)/config
REPORT_DIR = $(WD)/reports
IGNORE_DIRS = "/scratch"

M = $(shell printf "\033[34;1m\033[0m")
N = $(shell printf '\u2705')
DONE="$(N) Done: "
VERSION :=$(file < VERSION)
ifndef VERSION
	VERSION := dev
endif
GIT_TAG := $(shell git describe --exact-match --tags 2>git_describe_error.tmp; rm -f git_describe_error.tmp)
GIT_BRANCH := $(shell git branch --show-current)
GIT_COMMIT := $(shell git rev-parse HEAD)

MAKEFLAGS += --no-print-directory
GOPATH := $(shell go env GOPATH)
GOBIN := $(GOPATH)/bin

ARCHES ?= amd64 arm64
OSES ?= linux darwin windows
OUTTPL = $(DIST_DIR)/$(APP_NAME)-$(VERSION)-{{.OS}}_{{.Arch}}
GZCMD = tar -czf
ZIPCMD = zip
SHACMD = sha256sum
VET_RPT=vet.out
COVERAGE_RPT=coverage.out

LDFLAGS = -X '$(APP_NAME)/pkg/version.APP_NAME=$(APP_NAME)' \
	-X '$(APP_NAME)/pkg/version.Commit=$(GIT_COMMIT)' \
	-X '$(APP_NAME)/pkg/version.Branch=$(GIT_BRANCH)' \
	-X '$(APP_NAME)/pkg/version.Version=$(VERSION)' \
	-X '$(APP_NAME)/pkg/version.BuildTime=$(shell date -Iseconds)'

default: build

## deps: Download and Install any missing dependecies
.PHONY: deps
deps:
	go mod download -x
	@echo $(DONE) "Deps"

## tidy: Verifies and downloads all required dependencies
.PHONY: tidy
tidy:
	@echo "$(M) 🏃 go mod tidy..."
	@mkdir -p $(REPORT_DIR)
	#go mod verify
	#go mod tidy -e -v
	@if ! git diff --quiet; then \
		echo "WARNING:  'go mod tidy' resulted in changes or working tree is dirty. See diff.out for details"; \
		git --no-pager diff > $(REPORT_DIR)/diff.out; \
	fi
	@echo $(DONE) "Tidy\n"

## fmt: Runs gofmt on all source files
.PHONY: fmt
fmt:
	@echo "$(M) 🏃 gofmt..."
	@ret=0 && for d in $$(go list -f '{{.Dir}}' ./...); do \
		gofmt -l -w $$d/*.go || ret=$$? ; \
	 done ; exit $$ret
	@echo $(DONE) "Fmt\n"

## test: Tests code coverage
.PHONY: test
test:
	@echo "$(M)  👀 testing code...\n"
	@mkdir -pv $(REPORT_DIR)
	@touch $(REPORT_DIR)/test.out
	go test -count=1 -v ./... $(go list ./... | grep -v $IGNORE_DIRS) &> $(REPORT_DIR)/test.out
	@echo $(DONE) "Test\n"

## testwithcoverge: Tests code coverage
.PHONY: testwithcoverage
testwithcoverage:
	@echo "$(M)  👀 testing code with coverage...\n"
	@mkdir -pv $(REPORT_DIR)
	-go test -v ./pkg/... -coverprofile=$(REPORT_DIR)/$(COVERAGE_RPT)
	@echo $(DONE) "Test with Coverage\n"

## missing: Displays lines of code missing from coverage. Puts report in ./build/coverage.out
.PHONY: missing
missing: testwithcoverage
	@echo "$(M)  👀 missing coverage...\n"
	@mkdir -pv $(REPORT_DIR)
	go tool cover -func=$(REPORT_DIR)/$(COVERAGE_RPT) -o $(REPORT_DIR)/missing.out
	@echo $(DONE) "Missing\n"

## vet: Run go vet.  Puts report in ./build/vet.out
.PHONY: vet
vet:
	@echo "  $(M) 🏃 go vet..."
	@mkdir -pv $(REPORT_DIR)
	go vet -v ./... 2>&1 | tee $(REPORT_DIR)/vet.out
	@echo $(DONE) "Vet\n"

## reports: Runs vet, coverage, and missing reports
.PHONY: reports
reports: vet missing
	@echo $(DONE) "Reports\n"

## clean: Removes build, dist and report dirs
.PHONY: clean
clean:
	@echo "$(M)  🧹 Cleaning build ..."
	go clean ./... || true
	rm -rf $(REPORT_DIR)
	@echo $(DONE) "Clean\n"

.PHONY: dialects
dialects:
	$(eval export CGO_ENABLED = 0)
	go run ./cmd/gen-mavlink-dialects
	find ./pkg/dialects -type f -name '*.go' | xargs gofmt -l -w

## debug: Print make env information
.PHONY: debug
debug:
	$(info PATH=$(PATH))
	$(info GOPATH=$(GOPATH))
	$(info GOBIN=$(GOBIN))
	$(info CWD=$(CWD))
	$(info PKG=$(PKG))
	$(info APP_NAME=$(APP_NAME))
	$(info MD=$(MD))
	$(info WD=$(WD))
	$(info PKG_DIR=$(PKG_DIR))
	$(info CMD_DIR=$(CMD_DIR))
	$(info BUILD_DIR=$(BUILD_DIR))
	$(info DIST_DIR=$(DIST_DIR))
	$(info LOG_DIR=$(LOG_DIR))
	$(info REPORT_DIR=$(REPORT_DIR))
	$(info VET_RPT=$(VET_RPT))
	$(info COVERAGE_RPT=$(COVERAGE_RPT))
	$(info VERSION=$(VERSION))
	$(info GIT_COMMIT=$(GIT_COMMIT))
	$(info GIT_TAG=$(GIT_TAG))
	$(info GIT_BRANCH=$(GIT_BRANCH))
	$(info ARCHES=$(ARCHES))
	$(info OSES=$(OSES))
	$(info LDFLAGS=$(LDFLAGS))
	@echo -e $(DONE) "Debug\n"

.PHONY: help
help: Makefile
	@echo "Choose a command run in "$(PROJECTNAME)":"
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'