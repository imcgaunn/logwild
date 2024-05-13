#!/bin/zsh

setopt err_exit
setopt no_unset
setopt pipe_fail

declare -g OBSERVE_NETWORK_NAME="${OBSERVE_NETWORK_NAME:-observe}"
declare -g OBSERVE_NETWORK_CIDR="10.0.0.0/24" # try not to overlap with anything important

docker network create observe --subnet "${OBSERVE_NETWORK_CIDR}"
