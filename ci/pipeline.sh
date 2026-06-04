#!/usr/bin/env bash

set -euo pipefail

REPO="vkstrm-co/stuff"
IMAGE="packages.buildkite.com/${REPO}/bkite:latest"

# Propagate env variables
env | grep BUILDKITE > env.list
CHANGED_PATHS=$(git --no-pager diff --name-only HEAD~1 | tr '\n' ',')

buildkite-agent oidc request-token --audience "https://packages.buildkite.com/${REPO}" --lifetime 300 | docker login packages.buildkite.com/${REPO} --username buildkite --password-stdin
docker pull $IMAGE
docker run --rm -v $(pwd):$(pwd) -w $(pwd) --env CHANGED_PATHS="${CHANGED_PATHS}" --env-file env.list $IMAGE go run ci/pipeline/cmd/pipeline.go > custom-pipe.yaml

if [[ "$(cat custom-pipe.yaml)" == "{}" ]]; then
  buildkite-agent stop
fi

buildkite-agent pipeline upload custom-pipe.yaml
