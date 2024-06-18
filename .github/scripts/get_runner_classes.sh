#!/usr/bin/env bash
# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

#
# This script generates tag-sets that can be used as runs-on: values to select runners.

set -euo pipefail

case "$GITHUB_REPOSITORY" in
*-enterprise)
	# shellcheck disable=SC2129
	echo "compute-small=['ubuntu-latest']" >>"$GITHUB_OUTPUT"
	echo "compute-medium=['ubuntu-latest']" >>"$GITHUB_OUTPUT"
	echo "compute-large=['ubuntu-latest']" >>"$GITHUB_OUTPUT"
	# m5d.8xlarge is equivalent to our xl custom runner in CE
	echo "compute-xl=[ubuntu-latest']" >>"$GITHUB_OUTPUT"
	;;
*)
	# shellcheck disable=SC2129
	echo "compute-small=['ubuntu-latest']" >>"$GITHUB_OUTPUT"
	echo "compute-medium=['ubuntu-latest']" >>"$GITHUB_OUTPUT"
	echo "compute-large=['ubuntu-latest']" >>"$GITHUB_OUTPUT"
	echo "compute-xl=['ubuntu-latest']" >>"$GITHUB_OUTPUT"
	;;
esac
