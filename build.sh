#!/bin/bash

[ $(which go) ] || { echo "could not find go. please make sure it is installed and included in \$PATH"; exit 1; }

IMAGE_NAME=${1:-"ben/dynamic-pod-init"}

CGO_ENABLED=0 GOOS=linux go build -o main && \
chmod +x main

docker build -t $IMAGE_NAME .
