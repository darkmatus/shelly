#!/bin/bash
#
# This file is the base of all the automation scripts, it should be sourced at
# the top of the file as follow:
#
# source ./common/base.sh

#
# Common variables
#

# Colors
GREEN="\033[2;32m"
RED="\033[0;31m"
RESET="\033[0;0m"
YELLOW="\033[2;33m"


#
# Functions
#

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
