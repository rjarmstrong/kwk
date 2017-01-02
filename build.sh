#!/usr/bin/env bash

set -ef -o pipefail

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
#compile darwin
#compile windows

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

# CLEAN-UP
rm -fr ${tmp}