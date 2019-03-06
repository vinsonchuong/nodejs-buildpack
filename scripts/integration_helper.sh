#!/usr/bin/env bash

bp_dir=$1

cd bp_dir
base_name="$(basename $bp_dir)"
tar_name="$(ls -t | grep \"${bp_dir}.*.tgz\" | head -n 1)"

shasum="$(shasum -a 256 ${tar_name} | cut -d ' ' -f1)"

rm -rf "/Users/pivotal/.buildpack-packager/cache"

script=$(cat <<EOF
require 'YAML'
file = '/Users/pivotal/workspace/nodejs-buildpack/manifest.yml'
m = YAML.load_file(file)
m['dependencies'][0]['sha256'] = '$shasum'
File.open(file, 'w') {|f| f.write m.to_yaml }
EOF
)
ruby -e "$script"