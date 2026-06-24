#!/usr/bin/env bash

set -euo pipefail

REPO="vkstrm-co/stuff"
IMAGE="packages.buildkite.com/${REPO}/bkite:latest"

# Prepare a file with all BUILDKITE env variables
env | grep BUILDKITE > env.list
cat env.list

# Collect all paths in the repo that changed
CHANGED_PATHS=$(git --no-pager diff --name-only HEAD~1 | tr '\n' ',')

# Log into the Package Registry
buildkite-agent oidc request-token --audience "https://packages.buildkite.com/${REPO}" --lifetime 300 | docker login packages.buildkite.com/${REPO} --username buildkite --password-stdin

# Pull the custom image containing what the pipe needs to run the pipeline generator
docker pull $IMAGE

# Run the pipeline generator in the docker image
# Pass along the env variables in env.list file, as well as the changed paths
docker run --rm -v $(pwd):$(pwd) -w $(pwd) --env CHANGED_PATHS="${CHANGED_PATHS}" --env-file env.list $IMAGE go run ci/pipeline/cmd/main.go

# Print for troubleshooting or fun
cat pipe.yaml

# Tell the agent to use the dynamically generated pipe
buildkite-agent pipeline upload pipe.yaml
