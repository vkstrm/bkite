#!/usr/bin/env bash

set -euo pipefail

REPO="vkstrm-co/stuff"
IMAGE="packages.buildkite.com/${REPO}/bkite:latest"

# Propagate env variables
env | grep BUILDKITE > env.list

buildkite-agent oidc request-token --audience "https://packages.buildkite.com/${REPO}" --lifetime 300 | docker login packages.buildkite.com/${REPO} --username buildkite --password-stdin
docker pull $IMAGE
docker run --rm -v $(pwd):$(pwd) -w $(pwd) --env-file env.list $IMAGE go run pipeline/cmd/pipeline.go > custom-pipe.yaml

cat custom-pipe.yaml

buildkite-agent pipeline upload custom-pipe.yaml
