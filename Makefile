# Copyright 2017 The OpenPitrix Authors. All rights reserved.
# Use of this source code is governed by a Apache license
# that can be found in the LICENSE file.

TARG.Name:=runtime-provider-kubernetes
TRAG.Gopkg:=openpitrix.io/runtime-provider-kubernetes
TRAG.Version:=$(TRAG.Gopkg)/pkg/version

DOCKER_TAGS=latest
BUILDER_IMAGE=openpitrix/openpitrix-builder:release-v0.4.2
RUN_IN_DOCKER:=docker run -it -v `pwd`:/go/src/$(TRAG.Gopkg) -v `pwd`/tmp/mod:/go/pkg/mod -v `pwd`/tmp/cache:/root/.cache/go-build  -w /go/src/$(TRAG.Gopkg) -e GOBIN=/go/src/$(TRAG.Gopkg)/tmp/bin -e USER_ID=`id -u` -e GROUP_ID=`id -g` $(BUILDER_IMAGE)
GO_FMT:=goimports -l -w -e -local=openpitrix -srcdir=/go/src/$(TRAG.Gopkg)
GO_RACE:=go build -race
GO_VET:=go vet
GO_FILES:=./cmd ./pkg
GO_PATH_FILES:=./cmd/... ./pkg/...
define get_diff_files
    $(eval DIFF_FILES=$(shell git diff --name-only --diff-filter=ad | grep -E "^(cmd|pkg)/.+\.go"))
endef
define get_build_flags
    $(eval SHORT_VERSION=$(shell git describe --tags --always --dirty="-dev"))
    $(eval SHA1_VERSION=$(shell git show --quiet --pretty=format:%H))
	$(eval DATE=$(shell date +'%Y-%m-%dT%H:%M:%S'))
	$(eval BUILD_FLAG= -X $(TRAG.Version).ShortVersion="$(SHORT_VERSION)" \
		-X $(TRAG.Version).GitSha1Version="$(SHA1_VERSION)" \
		-X $(TRAG.Version).BuildDate="$(DATE)")
endef

CMD?=...
comma:= ,
empty:=
space:= $(empty) $(empty)
CMDS=$(subst $(comma),$(space),$(CMD))

.PHONY: help
help: ## This help
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z0-9_%-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

.PHONY: all
all: build ## Run generate and build

.PHONY: init-vendor
init-vendor: ## Initialize vendor and add dependence
	@if [ ! -f "$$(which govendor)" ]; then \
		go get -u github.com/kardianos/govendor; \
	fi
	govendor init
	govendor add +external
	@echo "init-vendor done"

.PHONY: update-vendor
update-vendor: ## Update dependence
	@if [ ! -f "$$(which govendor)" ]; then \
		go get -u github.com/kardianos/govendor; \
	fi
	govendor update +external
	govendor list
	@echo "update-vendor done"

.PHONY: fmt-all
fmt-all: ## Format all code
	$(RUN_IN_DOCKER) $(GO_FMT) $(GO_FILES)
	@echo "fmt done"

ARCH=$(shell uname -m)
.PHONY: check
check: ## go vet and race
	$(GO_RACE) $(GO_PATH_FILES)
	$(GO_VET) $(GO_PATH_FILES)

.PHONY: fmt
fmt: ## Format changed files
	$(call get_diff_files)
	$(if $(DIFF_FILES), \
		$(RUN_IN_DOCKER) $(GO_FMT) ${DIFF_FILES}, \
		$(info cannot find modified files from git) \
	)
	@echo "fmt done"

.PHONY: fmt-check
fmt-check: fmt-all ## Check whether all files be formatted
	$(call get_diff_files)
	$(if $(DIFF_FILES), \
		exit 2 \
	)

.PHONY: build
build: fmt ## Build all runtime-provider-kubernetes images
	mkdir -p ./tmp/bin
	$(call get_build_flags)
	$(RUN_IN_DOCKER) time go install -tags netgo -v -ldflags '$(BUILD_FLAG)' $(foreach cmd,$(CMDS),$(TRAG.Gopkg)/cmd/$(cmd))
	docker image prune -f 1>/dev/null 2>&1
	@echo "build done"

.PHONY: compose-update
compose-update: build compose-up ## Update service in docker compose
	@echo "compose-update done"

.PHONY: compose-up
compose-up: ## Launch runtime-provider-kubernetes in docker compose
	docker-compose up -d
	@echo "compose-up done"

.PHONY: compose-down
compose-down: ## Shutdown docker compose
	docker-compose down
	@echo "compose-down done"

BUILDX_BUILD_PUSH=docker buildx build --progress plain --platform linux/amd64,linux/arm64 --output=type=registry --push

build-push-image-%: ## build docker image
	$(BUILDX_BUILD_PUSH) -t openpitrix/runtime-provider-kubernetes:$* .
