#!/bin/bash

. VERSION

docker buildx build -t xxxsen/file_server:${VERSION} -t xxxsen/file_server:latest --platform=linux/amd64,linux/arm64 . --push