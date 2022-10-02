#!/bin/bash

VERSION=v0.0.1

docker buildx build -t xxxsen/file_server:${VERSION} -t xxxsen/file_server:latest --platform=linux/amd64,linux/arm64 . --push