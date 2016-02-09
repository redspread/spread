#!/usr/bin/env bash

# Setting up GOPATH

# create gopath
export GOPATH="$(pwd)/.gopath"
rm -rf $GOPATH
mkdir -p $GOPATH/src/rsprd.com

# link source to GOPATH
ln -sf $(pwd) $GOPATH/src/rsprd.com

# Copy in dependencies to get around Kube import check
cp -r ./vendor/* $GOPATH/src

go build rsprd.com/spread/cmd/spread
