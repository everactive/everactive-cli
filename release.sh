#!/bin/bash
set -e

PACKAGE_VERSION=0.1.0
BINARY_VERSION=0.1
CI_PROJECT_ID=41739213
CI_JOB_TOKEN="${GITLAB_PERSONAL_TOKEN_EVERACTIVE_CLI}"
CI_API_V4_URL=https://gitlab.com/api/v4
PACKAGE_REGISTRY_URL="${CI_API_V4_URL}/projects/${CI_PROJECT_ID}/packages/generic/everactive-cli-${BINARY_VERSION}/${PACKAGE_VERSION}"
echo "PACKAGE_REGISTRY_URL=${PACKAGE_REGISTRY_URL}"

LINUXAMD64_BINFILENAME="everactive-cli-linux-amd64-${BINARY_VERSION}.tar.gz"
DARWINAMD64_BINFILENAME="everactive-cli-darwin-amd64-${BINARY_VERSION}.tar.gz"
WINDOWSAMD64_BINFILENAME="everactive-cli-windows-amd64-${BINARY_VERSION}.tar.gz"
LINUXARM_BINFILENAME="everactive-cli-linux-arm-${BINARY_VERSION}.tar.gz"
set -x
curl "${PACKAGE_REGISTRY_URL}/${LINUXAMD64_BINFILENAME}" --header "Private-Token: ${CI_JOB_TOKEN}" --upload-file "release/${LINUXAMD64_BINFILENAME}"
curl "${PACKAGE_REGISTRY_URL}/${LINUXARM_BINFILENAME}" --header "Private-Token: ${CI_JOB_TOKEN}" --upload-file "release/${LINUXARM_BINFILENAME}"
curl "${PACKAGE_REGISTRY_URL}/${DARWINAMD64_BINFILENAME}" --header "Private-Token: ${CI_JOB_TOKEN}" --upload-file "release/${DARWINAMD64_BINFILENAME}"
curl "${PACKAGE_REGISTRY_URL}/${WINDOWSAMD64_BINFILENAME}" --header "Private-Token: ${CI_JOB_TOKEN}" --upload-file "release/${WINDOWSAMD64_BINFILENAME}"
set +x