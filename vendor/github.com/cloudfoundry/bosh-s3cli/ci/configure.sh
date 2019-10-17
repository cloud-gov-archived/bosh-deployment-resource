#!/usr/bin/env bash

fly -t director sp -p s3cli -c ${PROJECT_DIR}/ci/pipeline.yml \
  -l <(lpass show --notes "s3cli concourse secrets") \
  -l <(lpass show --notes "pivotal-tracker-resource-keys") \
  -l <(lpass show --note "bosh:docker-images concourse secrets")
