
# Tool commands (overridable)
GO_CMD             ?= go
DOCKER_CMD         ?= docker
DOCKER_COMPOSE_CMD ?= docker-compose

MAKEFILE_THIS      := $(lastword $(MAKEFILE_LIST))
THIS_PATH          := $(patsubst %/,%,$(dir $(abspath $(MAKEFILE_THIS))))

# Fabric versions used in the Makefile
FABRIC_STABLE_VERSION           := 2.2.0
FABRIC_STABLE_VERSION_MINOR     := 2.2

FABRIC_ARCH := amd64

# Fabric docker registries (overridable)
FABRIC_RELEASE_REGISTRY     ?=

# Docker-compose
BASE_DOCKER_COMPOSE_FILES 		   := -f ./docker-compose.yaml
#BASE_DOCKER_COMPOSE_FILES 		   := -f ./docker-compose-test-net.yaml
BASE_DOCKER_COMPOSE_FILES          := -f ./docker-compose-std.yaml $(BASE_DOCKER_COMPOSE_FILES)
DOCKER_COMPOSE_UP_FLAGS            := --remove-orphans --force-recreate
DOCKER_COMPOSE_UP_FLAGS       := $(DOCKER_COMPOSE_UP_FLAGS) --abort-on-container-exit


FABRIC_TOOLS_STABLE_TAG     = $(FABRIC_ARCH)-$(FABRIC_STABLE_VERSION)

# Fixture paths
FIXTURE_SCRIPTS_PATH      := $(THIS_PATH)/scripts
FIXTURE_DOCKERENV_PATH    := $(THIS_PATH)/fixtures/dockerenv
SCRIPTS_PATH      	  	  := scripts

# Local variables used by makefile
PROJECT_NAME           	  := jxwaChain-SDK-operation

# TODO: Make SDK compatible with latest cryptogen (https://jira.hyperledger.org/browse/FABG-977).
#FABRIC_CRYPTOCONFIG_VER         ?= v$(FABRIC_PREV_VERSION_MAJOR)
FABRIC_CRYPTOCONFIG_VER         ?= v1

# Fixture paths
FIXTURE_SCRIPTS_PATH      := $(THIS_PATH)/scripts
FIXTURE_DOCKERENV_PATH    := $(THIS_PATH)/fixtures/dockerenv
FIXTURE_CRYPTOCONFIG_PATH := $(THIS_PATH)/fixtures/fabric/$(FABRIC_CRYPTOCONFIG_VER)/crypto-config

# Fabric tool versions (overridable)
FABRIC_TOOLS_VERSION ?= $(FABRIC_STABLE_VERSION)

# Fabric tools docker image (overridable)
FABRIC_TOOLS_IMAGE ?= hyperledger/fabric-tools
FABRIC_TOOLS_TAG   ?= $(FABRIC_ARCH)-$(FABRIC_TOOLS_VERSION)

ifneq ($(FABRIC_RELEASE_REGISTRY),)
FABRIC_RELEASE_REGISTRY := $(FABRIC_RELEASE_REGISTRY)/
endif

#fabric-network-start 命令用于部署网络环境，包括对应的peer以及order节点等
.PHONY:fabric-network-start
fabric-network-start: clean depend-noforce populate-noforce
	@. $(FIXTURE_CRYPTOCONFIG_PATH)/env.sh && \
	    cd $(FIXTURE_DOCKERENV_PATH) && \
		$(DOCKER_COMPOSE_CMD) $(BASE_DOCKER_COMPOSE_FILES) up $(DOCKER_COMPOSE_UP_FLAGS)

#该命令用于把对链码的操作放在镜像中来执行
.PHONY: sdk-operation
sdk-operation:
	@cd $(FIXTURE_DOCKERENV_PATH) && $(DOCKER_COMPOSE_CMD) -f docker-compose-nopkcs11.yaml up $(DOCKER_COMPOSE_UP_FLAGS)

.PHONY: clean
clean: clean-temp clean-build

.PHONY: clean-temp
clean-temp:
	-rm -Rf /tmp/enroll_user /tmp/msp /tmp/keyvaluestore /tmp/hfc-kvs /tmp/state /tmp/state-store
	-rm -f integration-report.xml report.xml

.PHONY: clean-build
clean-build:
	@echo "now at clean-build!!"
	-$(GO_CMD) clean
	-FIXTURE_PROJECT_NAME=$(FIXTURE_PROJECT_NAME) DOCKER_REMOVE_FORCE=$(FIXTURE_DOCKER_REMOVE_FORCE) $(SCRIPTS_PATH)/clean_network.sh

.PHONY: populate-noforce
populate-noforce: populate-fixtures-stable-noforce

.PHONY: populate-fixtures-stable-noforce
populate-fixtures-stable-noforce:
	@FABRIC_CRYPTOCONFIG_VERSION=$(FABRIC_CRYPTOCONFIG_VER) \
	FABRIC_FIXTURE_VERSION=v$(FABRIC_STABLE_VERSION_MINOR) \
	$(SCRIPTS_PATH)/populate-fixtures.sh

.PHONY: populate-chaincode-vendor
populate-chaincode-vendor:
	#@cd fixtures/testdata/go/src/github.com/example_cc && go mod vendor
	#@cd fixtures/testdata/go/src/github.com/example_pvt_cc && go mod vendor
	#@cd fixtures/testdata/go/src/github.com/fire_alarm_cc && go mod vendor

.PHONY: depend-noforce
depend-noforce:	
	@echo "now at depend-noforce!!!"
	@$(SCRIPTS_PATH)/dependencies.sh

.PHONY: crypto-gen
crypto-gen:
	@echo "Generating crypto directory ..."
	@$(DOCKER_CMD) run -i -v /$(abspath .):/opt/workspace/$(PROJECT_NAME) -u $(shell id -u):$(shell id -g) \
		$(FABRIC_TOOLS_IMAGE):$(FABRIC_TOOLS_TAG) \
		//bin/bash -c "FABRIC_VERSION_DIR=fabric/$(FABRIC_CRYPTOCONFIG_VER) /opt/workspace/${PROJECT_NAME}/scripts/generate_crypto.sh"

.PHONY: channel-config-stable-gen
channel-config-stable-gen:
	@echo "Generating channel configuration transactions and blocks (code level stable) ..."
	@$(DOCKER_CMD) run -i \
		-v /$(abspath .):/opt/workspace/$(PROJECT_NAME) -u $(shell id -u):$(shell id -g) \
		$(FABRIC_TOOLS_IMAGE):$(FABRIC_TOOLS_STABLE_TAG) \
		//bin/bash -c "FABRIC_VERSION_DIR=fabric/v${FABRIC_STABLE_VERSION_MINOR}/ /opt/workspace/${PROJECT_NAME}/scripts/generate_channeltx.sh"
