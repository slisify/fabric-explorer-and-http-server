#!/bin/bash
#
# Copyright SecureKey Technologies Inc. All Rights Reserved.
#
# SPDX-License-Identifier: Apache-2.0
#
# Environment variables that affect this script:
# GO_TESTFLAGS: Flags are added to the go test command.
# GO_LDFLAGS: Flags are added to the go test command (example: -s).
# TEST_RACE_CONDITIONS: Boolean on whether to test for race conditions.
# FABRIC_SDKGO_CODELEVEL_TAG: Go tag that represents the fabric code target
# FABRIC_SDKGO_TESTRUN_ID: An identifier for the current run of tests.
# FABRIC_FIXTURE_VERSION: Version of fabric fixtures
# FABRIC_CRYPTOCONFIG_VERSION: Version of cryptoconfig fixture to use
# CONFIG_FILE: config file to use

set -e

GO_CMD="${GO_CMD:-go}"
FABRIC_SDKGO_CODELEVEL_TAG="${FABRIC_SDKGO_CODELEVEL_TAG:-stable}"
FABRIC_SDKGO_TESTRUN_ID="${FABRIC_SDKGO_TESTRUN_ID:-${RANDOM}}"
FABRIC_CRYPTOCONFIG_VERSION="${FABRIC_CRYPTOCONFIG_VERSION:-unknown}"
FABRIC_FIXTURE_VERSION="${FABRIC_FIXTURE_VERSION:-unknown}"
CONFIG_FILE="${CONFIG_FILE:-config_test.yaml}"
#CONFIG_FILE="${CONFIG_FILE:-config_e2e.yaml}"
TEST_LOCAL="${TEST_LOCAL:-false}"
TEST_RACE_CONDITIONS="${TEST_RACE_CONDITIONS:-true}"
SCRIPT_DIR="$(dirname "$0")"
CC_MODE="${CC_MODE:-lifecycle}"
# TODO: better default handling for FABRIC_CRYPTOCONFIG_VERSION

GOMOD_PATH=$(cd ${SCRIPT_DIR} && ${GO_CMD} env GOMOD)
PROJECT_DIR=$(dirname ${GOMOD_PATH})
PKG_FILE_NAME="main.go"

source ${SCRIPT_DIR}/lib/find_packages.sh
source ${SCRIPT_DIR}/lib/docker.sh

# Temporary fix for Fabric base image
unset GOCACHE

echo "Running" $(basename "$0")
echo "PROJECT_DIR is ${PROJECT_DIR}"

waitForCoreVMUp

echo "Code level ${FABRIC_SDKGO_CODELEVEL_TAG} (Fabric ${FABRIC_FIXTURE_VERSION})"
echo "Running integration tests ..."
echo "${PROJECT_DIR}"

export FABRIC_SDK_GO_PROJECT_PATH="${PROJECT_DIR}/"
echo "FABRIC_SDK_GO_PROJECT_PATH is ${FABRIC_SDK_GO_PROJECT_PATH}"

cd "${PROJECT_DIR}"

#${GO_CMD} run ${RACEFLAG} -tags "${GO_TAGS}" ${PKG_FILE_NAME} -p 1 -timeout=120m configFile=${CONFIG_FILE}

${GO_CMD} run ${PKG_FILE_NAME}

cd ${PWD_ORIG}
