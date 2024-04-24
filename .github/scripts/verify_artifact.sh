#!/bin/bash
# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1


set -euo pipefail

# verify_artifact.sh is the top-level script that implements the logic to decide
# which individual verification script to invoke. It decides which verification
# script to use based on artifact name it is given. By putting the logic in here,
# it keeps the workflow file simpler and easier to manage. It also doubles as a means
# to run verifications locally when necessary.

# set this so we can locate and execute the individual verification scripts.
SCRIPT_DIR="$( cd -- "$(dirname "$0")" >/dev/null 2>&1 ; pwd -P )"

function usage {
  echo "verify_artifact.sh <artifact_path> <expect_version>"
}

# Arguments:
#   $1 - artifact path (eg. /artifacts/consul-1.13.0~dev-1.i386.rpm)
#   $2 - expected version to match against (eg. v1.13.0-dev)
function main {
  local artifact_path="${1:-}"
  local expect_version="${2:-}"

  if [[ -z "${artifact_path}" ]]; then
    echo "ERROR: artifact path argument is required"
    usage
    exit 1
  fi

  if [[ -z "${expect_version}" ]]; then
    echo "ERROR: expected version argument is required"
    usage
    exit 1
  fi

  if [[ ! -e "${artifact_path}" ]]; then
    echo "ERROR: ${artifact_path} does not exist"
    usage
    exit 1
  fi

  # match against the various artifact names:
  # zip packages: consul_${version}_${os}_${arch}.zip
  case "${artifact_path}" in
    *.zip) verify_zip "${artifact_path}" "${expect_version}";;
    *)
      echo "${artifact_path} did not match known patterns"
      exit 1
      ;;
  esac
}


# Arguments:
#   $1 - path to zip (eg. consul_1.13.0-dev_linux_amd64.zip)
#   $2 - expected version to match against (eg. v1.13.0-dev)
function verify_zip {
  local artifact_path="${1:-}"
  local expect_version="${2:-}"
  local machine_os=$(uname -s)
  local machine_arch=$(uname -m)

  unzip "${artifact_path}"

  if [[ ! -e ./consul ]]; then
    echo "ERROR: ${artifact_path} did not contain a consul binary"
    exit 1
  fi

  case "${artifact_path}" in

    *_linux_386.zip | *_linux_amd64.zip)
      if [[ "${machine_os}" = 'Linux' && "${machine_arch}" = "x86_64" ]]; then
        # run the binary directly on the host when it's x86_64 Linux
        ${SCRIPT_DIR}/verify_bin.sh ./consul ${expect_version}
      else
        # otherwise, use Docker/QEMU
        docker run \
          --platform=linux/amd64 \
          -v $(pwd):/workdir \
          -v ${SCRIPT_DIR}:/scripts \
          -w /workdir  \
        amd64/debian \
        /scripts/verify_bin.sh \
        ./consul \
        "${expect_version}"
      fi
      ;;

    *_linux_arm.zip)
      if [[ "${machine_os}" = 'Linux' && "${machine_arch}" = arm* ]]; then
        # run the binary directly on the host when it's x86_64 Linux
        ${SCRIPT_DIR}/verify_bin.sh ./consul ${expect_version}
      else
        # otherwise, use Docker/QEMU
        docker run \
          --platform=linux/arm/v7 \
          -v $(pwd):/workdir \
          -v ${SCRIPT_DIR}:/scripts \
          -w /workdir  \
        arm32v7/debian \
        /scripts/verify_bin.sh \
        ./consul \
        "${expect_version}"
      fi
      ;;

    *_linux_arm64.zip)
      if [[ "${machine_os}" = 'Linux' && "${machine_arch}" = arm* ]]; then
        # run the binary directly on the host when it's x86_64 Linux
        ${SCRIPT_DIR}/verify_bin.sh ./consul ${expect_version}
      else
        # otherwise, use Docker/QEMU
        docker run \
          --platform=linux/arm64 \
          -v $(pwd):/workdir \
          -v ${SCRIPT_DIR}:/scripts \
          -w /workdir  \
        arm64v8/debian \
        /scripts/verify_bin.sh \
        ./consul \
        "${expect_version}"
      fi
      ;;

    *)
      echo "${artifact_path} did not match known patterns for zips"
      exit 1
      ;;
  esac
}

main "$@"
