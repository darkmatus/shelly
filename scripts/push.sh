#!/bin/bash
#
# Push Docker container to given registry using AWS credentials
#

set -e

source $(dirname $0)/base.sh

if [ -z "${CI_REGISTRY}" ]
then
  error "Do not push images from local machine!"
  exit 1
fi

info "Pushing images"
docker push --all-tags ${REGISTRY_IMAGE}