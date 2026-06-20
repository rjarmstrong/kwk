#!/usr/bin/env bash

set -ef -o pipefail

# BUILD
REVISION=$(git log --oneline -1)
VERSION=1.0.7
COMMIT=$(git log --pretty=format:'%h' -n 1)
DOCKER_IMAGE="kwk/fe"
REMOTE_TAG_PREFIX="477988974763.dkr.ecr.us-west-1.amazonaws.com/${DOCKER_IMAGE}"

echo ${REVISION}
echo "**** BUILDING ${DOCKER_IMAGE}+${VERSION}+${COMMIT} ****"

#echo "TESTING"
#echo "no tests"
#sed -i -- "s@\[git_revision\]@${REVISION} - ${DATE}@g" ./client/index.html

echo "NPM (outside container because of private npm repos)"
npm install
npm run dev

echo "BUILDING IMAGE"
docker build --rm=true --disable-content-trust -t ${DOCKER_IMAGE} .
docker tag ${DOCKER_IMAGE} ${REMOTE_TAG_PREFIX}:latest
docker tag ${DOCKER_IMAGE} ${REMOTE_TAG_PREFIX}:${VERSION}

#echo "PUSHING"
$(aws ecr get-login --region us-west-1 --profile kwk --no-include-email)
docker push ${REMOTE_TAG_PREFIX}:${VERSION}
echo "**** COMPLETED ${DOCKER_IMAGE}+${VERSION} ****"