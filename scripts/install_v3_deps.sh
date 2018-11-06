#!/bin/bash
set -euo pipefail

cd "$( dirname "${BASH_SOURCE[0]}" )/.."
source .envrc

if [[ ! -f .bin/detector ]]; then
  (cd src/*/vendor/github.com/buildpack/lifecycle/cmd/detector && go install)
fi
if [[ ! -f .bin/builder ]]; then
  (cd src/*/vendor/github.com/buildpack/lifecycle/cmd/builder && go install)
fi
