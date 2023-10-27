#!/bin/bash
#
# Build
#

set -e

source $(dirname $0)/base.sh

# put current version into file
echo ${TAG} > current_version

# for the cluster part we don't need to build the mysql container
  docker build -t ${IMAGE}:shelly -f scripts/docker/go/Dockerfile .

success "ALL DONE!"
