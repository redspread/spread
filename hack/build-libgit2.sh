#!/bin/sh

set -ex

VENDORED_PATH=vendor/libgit2

cd $VENDORED_PATH &&
mkdir -p install/lib &&
mkdir -p build &&
cmake -DTHREADSAFE=ON \
      -DBUILD_CLAR=OFF \
      . &&
make &&
sudo make install
