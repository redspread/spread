#!/bin/sh

set -ex

VENDORED_PATH=vendor/libgit2

cd $VENDORED_PATH &&
mkdir -p install/lib &&
mkdir -p build &&
cd build &&
cmake . &&
make &&
sudo make install
