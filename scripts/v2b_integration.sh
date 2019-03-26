#!/usr/bin/env bash

package_bp() {
    bp_dir=$(~/workspace/$1/scripts/package.sh | grep "Buildpack .tar into: " | awk -F " " '{print $NF}')
}


update_yaml () {
    shasum=$(shasum -a 256 "$1" | cut -f 1 -d " ")
    echo "$shasum"
    script=$(cat <<EOF
require 'YAML'
file = '/Users/pivotal/workspace/nodejs-buildpack/manifest.yml'
m = YAML.load_file(file)
m['dependencies'][$2]['sha256'] = '$shasum'
m['dependencies'][$2]['uri'] = 'file://' + "$bp_dir"
File.open(file, 'w') {|f| f.write m.to_yaml }
EOF
)
ruby -e "$script"
}

tarbp() {
    tar -czvf "$1".tgz "$1"
    # clean up
    rm -rf $bp_dir
    rm -rf "/Users/pivotal/.buildpack-packager/cache"
}


set -euo pipefail

# clear packager cache
rm -rf "/Users/pivotal/.buildpack-packager/cache"

cd "$( dirname "${BASH_SOURCE[0]}" )/.."
source .envrc
./scripts/install_tools.sh

GINKGO_NODES=${GINKGO_NODES:-3}
GINKGO_ATTEMPTS=${GINKGO_ATTEMPTS:-1}
export CF_STACK=${CF_STACK:-cflinuxfs3}

bp_dir=""
BPCOUNTER=0
for bp in nodejs-cnb npm-cnb yarn-cnb nodejs-compat-cnb; do
    echo "packaging bp"
    package_bp "$bp"
    echo "created $bp_dir"
    #tarbp "$bp_dir"
    update_yaml "$bp_dir" $BPCOUNTER
    BPCOUNTER=$(expr $BPCOUNTER + 1)
done


pushd v2b_integration
    echo "Run Uncached Shim Buildpack For V2B specs"
    ginkgo -r --flakeAttempts=$GINKGO_ATTEMPTS -nodes $GINKGO_NODES --slowSpecThreshold=60 -- --cached=true
popd
