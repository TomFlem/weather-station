#!/bin/bash
set -e -x
export MNAME=`uname -m`
export VER=`cat ./VERSION`

case $MNAME in
  x86_64)
    ARCH=x86_64
    ;;
  aarch64)
    ARCH=arm64
    ;;
  armv7l)
    ARCH=arm
    ;;
  *)
    ARCH=$MNAME
    ;;
esac
export ARCH