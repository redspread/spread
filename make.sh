#!/usr/bin/env bash

# Setting up GOPATH

# create gopath
export GOPATH="$(pwd)/.gopath"
rm -rf $GOPATH
mkdir -p $GOPATH/src/rsprd.com

# link source to GOPATH
cp -r pkg,cmd,cli $GOPATH/src/rsprd.com/spread/

# Copy in dependencies to get around Kube import check
cp -r ./vendor/* $GOPATH/src

go build rsprd.com/spread/cmd/spread