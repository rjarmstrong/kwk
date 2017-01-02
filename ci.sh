#!/usr/bin/env bash

KWK_VERSION=v1.0.9
docker build . --rm -t kwkcli
removeContainer

docker run --name kwkcli kwkcli ls

# UPLOAD
export BUILDKITE_S3_ACCESS_KEY_ID=AKIAJRJBQNMZWLG653WA
export BUILDKITE_S3_SECRET_ACCESS_KEY=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
export BUILDKITE_S3_DEFAULT_REGION="us-east-1"
export BUILDKITE_ARTIFACT_UPLOAD_DESTINATION="s3://kwk-cli/${BUILDKITE_JOB_ID}"

buildkite-agent artifact upload /builds/${KWK_VERSION} s3://kwk-cli/${KWK_VERSION}

removeContainer

function removeContainer(){
    lines=$(docker ps -a -q -f name=kwkcli | wc -l)
    if [[ ${lines} > 0 ]]; then
        docker rm -f kwkcli
    fi
}