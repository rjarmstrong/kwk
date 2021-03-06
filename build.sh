#!/usr/bin/env bash

set -ef -o pipefail
KWK_VERSION=v0.8.19
BUILD_NUMBER=$1
RELEASE_TIME=$(date +%s)
RELEASE_NOTES="- ....\n"

echo -e "\n\n\n**** kwk-cli build ${KWK_VERSION}+${BUILD_NUMBER} *****\n\n\n"

ARCH=amd64

# TESTING
./src/test.sh


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
  cd ./src/
  env GOOS=${os} GOARCH=${ARCH} CGO_ENABLED=0 go build -ldflags "-s -w -X main.version=${KWK_VERSION} -X main.build=${BUILD_NUMBER} -X main.releaseTime=${RELEASE_TIME}" -x -o ${binary}
  cd ..

  # ZIP
  zipped=/${binPath}/${file}.tar.gz
  tar cvzf ${zipped} -C ${tmp}/bin ${file}

  # CHECKSUM
  sha1sum ${zipped} > ${zipped}.sha1
}

sed -i -- "s/RELEASE_VERSION/${KWK_VERSION}/" ./src/main.go
compile linux
compile darwin
# compile windows

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
npmTar=${npmPath}/kwk-cli-npm.tar.gz
tar cvzf ${npmTar} -C ${npmTemp} .
sha1sum ${npmTar} > ${npmTar}.sha1

# CLEAN-UP
rm -fr ${tmp}

aws s3 cp /builds/${KWK_VERSION} s3://kwk-cli/${KWK_VERSION} --recursive --acl public-read
aws s3 cp s3://kwk-cli/${KWK_VERSION} s3://kwk-cli/latest --recursive --acl public-read

echo "{
\"version\":\"${KWK_VERSION}\",
\"build\":\"${BUILD_NUMBER}\",
\"time\":${RELEASE_TIME},
\"notes\":\"${RELEASE_NOTES}\"
}" > /builds/release-info.json

aws s3 cp /builds/release-info.json  s3://kwk-cli/release-info.json --acl public-read