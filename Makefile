PROJ_NAME = bhlquest

VERSION = $(shell git describe --tags)
VER = $(shell git describe --tags --abbrev=0)
DATE = $(shell TZ=UTC date +'%Y-%m-%d_%H:%M:%ST%Z')

NO_C = CGO_ENABLED=0
FLAGS_LD = -ldflags "-X codeberg.org/dimus/$(PROJ_NAME)/pkg.Build=$(DATE) \
                     -X codeberg.org/dimus/$(PROJ_NAME)/pkg.Version=$(VERSION)"
FLAGS_REL = -trimpath -ldflags "-s -w -X codeberg.org/dimus/$(PROJ_NAME)/pkg.Build=$(DATE)"
FLAGS_SHARED = $(NO_C) GOARCH=amd64

RELEASE_DIR = /tmp
TEST_OPTS = -v -shuffle=on -race -coverprofile=coverage.txt -covermode=atomic

GOCMD = go
GOTEST = $(GOCMD) test
GOVET = $(GOCMD) vet
GOBUILD = $(GOCMD) build $(FLAGS_LD)
GORELEASE = $(GOCMD) build $(FLAGS_REL)
GOCLEAN = $(GOCMD) clean
GOINSTALL = $(GOCMD) install $(FLAGS_LD)

GREEN  := $(shell tput -Txterm setaf 2)
YELLOW := $(shell tput -Txterm setaf 3)
WHITE  := $(shell tput -Txterm setaf 7)
RESET  := $(shell tput -Txterm sgr0)

.PHONY: all build clean

all: install

## Dependencies
deps: ## Download dependencies
	$(GOCMD) mod download;

## Tools
tools: deps ## Install tools
	@cat tools.go | grep _ | awk -F'"' '{print $$2}' | xargs -tI % go install %

## Build:
build: ## Build binary
	$(NO_C) $(GOCMD) build \
		-o $(PROJ_NAME) \
		$(FLAGS_LD) \
		.

## Build Release
buildrel: ## Build binary without debug info and with hardcoded version
	$(NO_C) $(GOCMD) build \
		-o $(PROJ_NAME) \
		$(FLAGS_REL) \
		.

## Install:
install: ## Build and install binary
	$(NO_C) $(GOINSTALL)

## Release
release: buildrel ## Build and package binaries for a release
	$(GOCLEAN); \
	$(FLAGS_SHARED) GOOS=linux $(GORELEASE); \
	tar zcvf $(RELEASE_DIR)/$(PROJ_NAME)-$(VER)-linux.tar.gz $(PROJ_NAME); \
	$(GOCLEAN);

## Clean
clean: ## Clean all the files and binaries generated by the Makefile
	rm -rf ./out

## Test:
test: ## Run the tests of the project
	$(GOTEST) $(TEST_OPTS) ./...

## Coverage
coverage: ## Run the tests of the project and export the coverage
	$(GOTEST) -cover -covermode=count -coverprofile=profile.cov ./...
	$(GOCMD) tool cover -func profile.cov

## Help:
help: ## Show this help
	@echo ''
	@echo 'Usage:'
	@echo '  $(YELLOW)make${RESET} $(GREEN)<target>$(RESET)'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} { \
		if (/^[0-9a-zA-Z_-]+:.*?##.*$$/) \
		  {printf "    $(YELLOW)%-20s$(GREEN)%s$(RESET)\n", $$1, $$2} \
		else if (/^## .*$$/) {printf "  $(CYAN)%s$(RESET)\n", substr($$1,4)} \
		}' $(MAKEFILE_LIST)

## Version
version: ## Display current version
	@echo $(VERSION)
