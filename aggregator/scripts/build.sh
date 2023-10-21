#!/bin/sh
. ./scripts/env.sh

if [ "$1" != "" ]; then
  BUILD_TAG=$VER.$1
else
  BUILD_TAG=$VER
fi

docker build --force-rm --no-cache -t tomflemiotech/weather-aggregator:$BUILD_TAG -f ./scripts/dockerfiles/Dockerfile .