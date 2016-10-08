#!/usr/bin/env bash

cd ./libs/app/
$GOPATH/bin/goconvey -port 3456 -launchBrowser false -depth 3 -cover false