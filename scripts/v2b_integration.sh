#!/usr/bin/env bash
set -euo pipefail

cd "$( dirname "${BASH_SOURCE[0]}" )/.."
source .envrc
./scripts/install_tools.sh

GINKGO_NODES=${GINKGO_NODES:-3}
GINKGO_ATTEMPTS=${GINKGO_ATTEMPTS:-1}
export CF_STACK=${CF_STACK:-cflinuxfs3}

~/workspace/nodejs-compat-cnb/scripts/package.sh
~/workspace/nodejs-cnb/scripts/package.sh
~/workspace/npm-cnb/scripts/package.sh
~/workspace/yarn-cnb/scripts/package.sh

pushd v2b_integration
    echo "Run Uncached Shim Buildpack For V2B specs"
    ginkgo -r --flakeAttempts=$GINKGO_ATTEMPTS -nodes $GINKGO_NODES --slowSpecThreshold=60 -- --cached=true
popd
