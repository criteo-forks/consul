#!/usr/bin/env bash
#
# This script generates tag-sets that can be used as runs-on: values to select runners.

set -euo pipefail

# shellcheck disable=SC2129
echo "compute-small=['ubuntu-latest']" >> "$GITHUB_OUTPUT"
echo "compute-medium=['ubuntu-latest']" >> "$GITHUB_OUTPUT"
echo "compute-large=['ubuntu-latest']" >> "$GITHUB_OUTPUT"
echo "compute-xl=['ubuntu-latest']" >> "$GITHUB_OUTPUT"

