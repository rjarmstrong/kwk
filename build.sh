#!/usr/bin/env bash

set -ef -o pipefail

KWK_VERSION=v1.0.9
ARCH=amd64
BUILD=${BUILDKITE_BUILD_NUMBER}

# TESTING
go test ./app


# PREP OUTPUT
releasePath=/builds/${KWK_VERSION}
tmp=/builds/temp
binPath=${releasePath}/bin
npmPath=${releasePath}/npm

if [[ -d "${tmp}" ]]; then
    rm -fr ${tmp}
else
    mkdir ${tmp}
fi

if [[ ! -d "${releasePath}" ]]; then
    mkdir ${releasePath}
fi

if [[ ! -d "${binPath}" ]]; then
    mkdir ${binPath}
fi

if [[ ! -d "${npmPath}" ]]; then
    mkdir ${npmPath}
fi

# COMPILING
function compile(){
  os=$1
  file="kwk-${os}-${ARCH}"

  # COMPILE
  binary=${tmp}/bin/${file}
  env GOOS=${os} GOARCH=${ARCH} go build -ldflags "-X main.version=${KWK_VERSION} -X main.build=${BUILD}" -x -o ${binary}

  # ZIP
  zipped=${binPath}/${file}.tar.gz
  tar cvzf ${zipped} -C ${tmp}/bin ${file}

  # CHECKSUM
  sha1sum ${zipped} > ${zipped}.sha1
}

sed -i -- "s/RELEASE_VERSION/${KWK_VERSION}/" ./main.go
compile linux
compile darwin
compile windows

# CREATE NPM
npmTemp=${tmp}/npm
rm -fr ${npmTemp}
mkdir ${npmTemp}

cp -R dist/npm/. ${npmTemp}
cp -R ${tmp}/bin/. ${npmTemp}/bin
sed -i -- "s/RELEASE_VERSION/${KWK_VERSION}/" ${npmTemp}/package.json
tree ${npmTemp}

# CREATE NPM TAR
echo "CREATING NPM ARCHIVE"
tar cvzf ${npmPath}/kwk-cli-npm.tar.gz -C ${npmTemp} .

# UPLOAD
export BUILDKITE_S3_ACCESS_KEY_ID=AKIAJRJBQNMZWLG653WA
export BUILDKITE_S3_SECRET_ACCESS_KEY=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
export BUILDKITE_S3_DEFAULT_REGION="us-east-1"
export BUILDKITE_ARTIFACT_UPLOAD_DESTINATION="s3://kwk-cli/${BUILDKITE_JOB_ID}"

buildkite-agent artifact upload ${releasePath} s3://kwk-cli/${KWK_VERSION}

# CLEAN-UP
rm -fr ${tmp}