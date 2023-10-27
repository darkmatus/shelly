#!/bin/bash

# Colors
GREEN="\033[2;32m"
RED="\033[0;31m"
RESET="\033[0;0m"
YELLOW="\033[2;33m"

# Tag used for Docker image versioning
tag=$(git describe --tag --always --abbrev=0 | tr + _)
major_version=$(echo $tag | cut -d '.' -f1)

# if tag doesn't container 1.0.0_1 tag format, use commit number instead
# default major version is 0
if [ "$tag" == "" ]
then
  major_version=0
fi

export TAG=${TAG:-$tag}

# Check presence of a given list of binaries.
# Usage:
#   require_binary "kubectl docker git"
function require_binary() {
  for cmd in $1
  do
    info "Checking if binary ${cmd} is present"
    command -v $cmd >/dev/null 2>&1 || \
    {
      error >&2 "Binary ${cmd} is required to execute this script. Aborting.";
      exit 1;
    }
  done
}

# Logging functions
function error() {
  echo -e "[ERROR] ${RED}${1}${RESET}"
}

function info() {
  echo -e "[INFO] ${YELLOW}${1}${RESET}"
}

function success() {
  echo -e "[SUCCESS] ${GREEN}${1}${RESET}"
}

# Exponentially retry 5 times a given command
function retry() {
  local cmd=$1
  local next_wait_time=0
  set +e
  until $cmd || [ $next_wait_time -eq 4 ]; do
    ((NEXT_WAIT_TIME++))
    local sleep_time=$((next_wait_time**next_wait_time))

    error "Failed executing command, ${next_wait_time} attempt(s):"
    error "  > ${1}"
    error "Retrying in ${sleep_time} seconds"
    sleep $sleep_time
  done
  set -e
}

# Wait for an endpoint to return status code > 200 and < 400
# Usage:
#   wait_for_endpoint localhost:8080/ready
#
# You can specify a timeout, default 300s
#   wait_for_endpoint localhost:8080/ready 300
function wait_for_endpoint() {
  require_binary "curl"

  local endpoint=${1:-"localhost"}
  local timeout=${2:-"300"}
  local status_code="500"

  while [ "$status_code" -lt 200 -o "$status_code" -ge 400 ]
  do
      info "Waiting for the application to be ready (checking $endpoint)"

      local http_response=$(curl --silent --write-out "HTTPSTATUS:%{http_code}" $endpoint)
      local http_body=$(echo $http_response | sed -e 's/HTTPSTATUS\:.*//g')
      local status_code=$(echo $http_response | tr -d '\n' | sed -e 's/.*HTTPSTATUS://')

      local timeout=$[$timeout-1]
      if [ "$timeout" = 0 ]
      then
        echo "App seems to be broken as it don't come up. Aborting..."
        echo ""
        echo "Last health-probe response is:"
        echo "==================================================================="
        echo "HTTP Response Body: $http_body"
        echo "==================================================================="
        echo "Status code: $status_code"
        echo "==================================================================="
        echo "The Docker-logs are:"
        echo "==================================================================="
        docker-compose logs --tail=10
        echo "==================================================================="
        exit 1
      fi

      sleep 1
  done
}

# Graceful shutdown of docker-compose in case of error or at the end of the
# execution. The trap command should be use to detect error of exit of
# the script.
# Usage:
#   trap 'gracefull_shutdown_docker_compose' ERR EXIT

# Usage with custom docker-compose flag, for single file only:
#   trap 'gracefull_shutdown_docker_compose "-f docker-compose.yml"' ERR EXIT
function gracefull_shutdown_docker_compose() {
  exit_code=$?

  require_binary "docker-compose"

  local flags=${1}

  trap '' EXIT

  info "Stopping and removing containers"
  docker-compose $flags down

  if [ "$exit_code" -ne "0" ]
  then
    error "Failed"
    exit 1
  fi

  success "Done"
}

#
# Common variables
#
export IMAGE="darkmatus/shelly"

build_image=go-base-build