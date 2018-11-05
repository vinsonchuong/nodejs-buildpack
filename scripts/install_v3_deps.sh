#!/bin/bash
set -euo pipefail

cd "$( dirname "${BASH_SOURCE[0]}" )/.."
source .envrc

if [[ ! -f .bin/v3lifecycle/detector ]]; then
  (cd src/*/vendor/github.com/buildpack/lifecycle/cmd/detector && go install)
fi
if [[ ! -f .bin/v3lifecycle/builder ]]; then
  (cd src/*/vendor/github.com/buildpack/lifecycle/cmd/builder && go install)
fi
