#!/usr/bin/env bash

# Setting up GOPATH

# create gopath
export GOPATH="$(pwd)/.gopath"
rm -rf $GOPATH
mkdir -p $GOPATH/src/rsprd.com/spread

# Ensure vendoring is enabled (for 1.5)
export GO15VENDOREXPERIMENT=1

# link source to GOPATH
cp -r * $GOPATH/src/rsprd.com/spread/

go build -v rsprd.com/spread/cmd/spread
