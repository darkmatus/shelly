#!/bin/bash
#
# This file is the base of all the automation scripts, it should be sourced at
# the top of the file as follow:
#
#   . ./common/base.sh

#
# Common variables
#

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

# Common Backend variables
#

COMPOSER_VERSION=${COMPOSER_VERSION:-"latest-stable"}

#
# Functions
#


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

# Login to private AWS ECR registry using AWS credentials from environment
# variable
function docker_login() {
  info "Login to private Docker registry using AWS credentials..."

  if [[ -z $AWS_REGION ]] || [[ -z $AWS_ACCESS_KEY_ID ]] || [[ -z $AWS_SECRET_ACCESS_KEY ]]
  then
    error "AWS credentials are not set to log into Docker registry, assuming local testing"
  else
    info "Using AWS credentials to login to Docker"
    require_binary "aws"

    info "Getting Docker registry credentials"
    docker_login=$(aws ecr get-login --no-include-email --region=${AWS_REGION})

    info "Login Docker to registry"
    eval "${docker_login}"
  fi
}

# Download Composer binary which is used to install PHP dependencies
# Usage:
#   download_composer
# You can also specify a specific version of Composer:
#   download_composer 1.5.2
function download_composer() {
  require_binary "wget"

  local composer_version=${1:-$COMPOSER_VERSION}

  if [ ! -f composer.phar ]; then
    wget -O composer.phar https://getcomposer.org/download/${composer_version}/composer.phar
  fi
  chmod +x composer.phar
}


# If build docker images using Docker from Minikube, we mount files from the
# /hosthome directory instead of /home directory.
# Reference: https://kubernetes.io/docs/getting-started-guides/minikube/#mounted-host-folders
function get_mount_path() {
  local mount_path=$(pwd)
  if [[ "${DOCKER_CERT_PATH}" =~ ^.*minikube.*$ ]] && [[ `uname -s` == "Linux" ]]
  then
    info "Building from Linux using Docker Minikube" >&2
    info "Setting mounting path to /hosthome/ directory instead of /home/" >&2
    local mount_path="${mount_path/\/home\///hosthome/}"
  fi
  echo $mount_path
}


# Install PHP dependencies using Composer. To speedup installation AWS S3 bucket
# is used to store the vendor directory.
# Usage:
#   composer_install
#
# You can also specify a non-standard build image:
#   composer_install custom-php-base-image
function composer_install() {
  require_binary "docker tar"

  local build_image=${1:-"php-base-build"}
  local user_id=${2:-""}

  if [ ! -z "$user_id" ]
  then
    local user_id="-u ${user_id}"
  fi

  local mount_path=$(get_mount_path)

  local composer_file="composer.json"
  if [ -e "composer.lock" ]; then
    composer_file="composer.lock"
  fi
#  echo $(pwd)
  ls -la $mount_path
#  docker run --rm -v $(pwd):/app composer install

  docker run --rm \
    -e COMPOSER_HOME=/composer \
    -e COMPOSER_ALLOW_SUPERUSER=1 \
    -v ~/.composer:/composer \
    $user_id \
    -v $mount_path:/var/www \
    -w /var/www $build_image php composer.phar install
}

# Generate documentation based on the Aglio Docker image
# Ref: https://github.com/danielgtaylor/aglio
# Usage:
#   generate_docs_aglio "-i var/docs/api/api.apib -o web/docs/index.html --theme-variables flatly --theme-full-width"
function generate_docs_aglio() {
  require_binary "docker"

  local mount_path=$(get_mount_path)

  docker run --rm \
    -v $mount_path:/var/www \
    -w /var/www \
    christianbladescb/aglio $1
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
