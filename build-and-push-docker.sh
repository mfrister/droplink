#!/usr/bin/env bash

set -e

VERSION=$(date +%Y-%m-%d)
IMAGE=droplink

docker build --rm -t localhost:5000/$IMAGE:$VERSION .
docker tag localhost:5000/$IMAGE:$VERSION localhost:5000/$IMAGE:latest
docker push localhost:5000/$IMAGE:latest
