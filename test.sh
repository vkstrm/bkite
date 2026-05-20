#!/usr/bin/env bash

set -euo pipefail

REPO="vkstrm-co/stuff"
IMAGE="packages.buildkite.com/${REPO}/bkite:latest"

if [[ $# == 0 ]]; then
  echo "Requires arguments"
  exit 0
fi

buildkite-agent oidc request-token --audience "https://packages.buildkite.com/${REPO}" --lifetime 300 | docker login packages.buildkite.com/${REPO} --username buildkite --password-stdin
docker pull $IMAGE
docker run --rm -v $(pwd):$(pwd) -w $(pwd) $IMAGE "$@"

curl \
  -X POST \
  --fail-with-body \
  -H "Authorization: Token token=\"$BUILDKITE_ANALYTICS_TOKEN\"" \
  -F "data=@junit.xml" \
  -F "format=junit" \
  -F "run_env[CI]=buildkite" \
  -F "run_env[key]=$BUILDKITE_BUILD_ID" \
  -F "run_env[number]=$BUILDKITE_BUILD_NUMBER" \
  -F "run_env[job_id]=$BUILDKITE_JOB_ID" \
  -F "run_env[branch]=$BUILDKITE_BRANCH" \
  -F "run_env[commit_sha]=$BUILDKITE_COMMIT" \
  -F "run_env[message]=$BUILDKITE_MESSAGE" \
  -F "run_env[url]=$BUILDKITE_BUILD_URL" \
  https://analytics-api.buildkite.com/v1/uploads

