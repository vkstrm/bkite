#!/usr/bin/env bash

set -euo pipefail

REPO="vkstrm-co/stuff"
IMAGE="packages.buildkite.com/${REPO}/bkite:latest"

buildkite-agent oidc request-token --audience "https://packages.buildkite.com/${REPO}" --lifetime 300 | docker login packages.buildkite.com/${REPO} --username buildkite --password-stdin
docker pull $IMAGE
docker run --rm -v $(pwd):$(pwd) -w $(pwd) $IMAGE go run pipeline/cmd/pipeline.go > pipeline.yaml

# go run pipeline/cmd/pipeline.go > pipeline.yaml
buildkite-agent pipeline upload pipeline.yaml
