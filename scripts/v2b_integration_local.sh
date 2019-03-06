#!/usr/bin/env bash
set -euo pipefail

cd "$( dirname "${BASH_SOURCE[0]}" )/.."
source .envrc
./scripts/install_tools.sh

# package up script
#~/workspace/nodejs-compat-cnb/scripts/package.sh
# parse output to get path of packaged cnb,
# calc sha
# update manifest?
# revert manifest changes?
nodejs_path="~/workspace/nodejs-cnb"
npm_path="~/workspace/npm-cnb"
yarn_path="~/workspace/yarn-cnb"

# Generate tars for each CNB buildpack
${nodejs_path}/scripts/package.sh
${npm_path}/scripts/package.sh
${yarn_path}/scripts/package.sh

# Calculate SHAs for each of the CNB tars and insert into Nodejs buildpack's manifest.yml
./scripts/integration_helper.sh ${nodejs_path}
./scripts/integration_helper.sh ${npm_path}
./scripts/integration_helper.sh ${yarn_path}

./scripts/v2b_integration.sh
