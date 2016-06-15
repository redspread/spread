#!/bin/sh

set -ex

VENDORED_PATH=vendor/libgit2

cd $VENDORED_PATH &&
mkdir -p _install &&
mkdir -p build &&
cd build &&
cmake -DTHREADSAFE=ON \
      -DBUILD_CLAR=OFF \
      .. &&
cmake --build .
