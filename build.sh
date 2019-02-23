#!/bin/sh
export GO111MODULE=on
GOOS=linux GOARCH=amd64 go build -v -a -o builds/qsub_linux
GOOS=darwin GOARCH=amd64 go build -v -a -o builds/qsub_darwin

