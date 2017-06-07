#!/usr/bin/env bash

go test -v ./app/handlers/ -coverprofile cover.out; go tool cover -html=cover.out