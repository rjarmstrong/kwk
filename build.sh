#!/usr/bin/env bash

set -ef -o pipefail

KWK_VERSION=v1.0.8
BUILD=${BUILDKITE_BUILD_NUMBER}

# TESTING
go test ./app


# PREP OUTPUT
output=/builds/${KWK_VERSION}
tmp=/builds/temp

if [[ -d "${tmp}" ]]; then
    rm -fr ${tmp}
else
    mkdir ${tmp}
fi

if [[ ! -d "${output}" ]]; then
    mkdir ${output}
fi


# COMPILING
arch=amd64

function compile(){
  os=$1
  file="kwk-${os}-${arch}"

  # COMPILE
  input=${tmp}/${file}
  env GOOS=${os} GOARCH=${arch} go build -ldflags "-X main.version=${KWK_VERSION} -X main.build=${BUILD}" -x -o ${input}

  # ZIP
  zipped=${output}/${file}.tar.gz
  tar cvzf ${zipped} ${input}

  # CHECKSUM
  sha1sum ${zipped} > ${zipped}.sha1
}

sed -i -- "s/RELEASE_VERSION/${KWK_VERSION}/" ./main.go
compile linux
compile darwin
compile windows

# CREATE NPM
rm -fr /builds/npm
mkdir /builds/npm

npm_dir=/builds/npm/${KWK_VERSION}
mkdir ${npm_dir}
echo $PWD
cp -R dist/npm/. ${npm_dir}/

bin_dir=${npm_dir}/bin
cp -R ${tmp}/. ${bin_dir}/
sed -i -- "s/RELEASE_VERSION/${KWK_VERSION}/" ${npm_dir}/package.json

tree ${npm_dir}

# CLEAN-UP
rm -fr ${tmp}