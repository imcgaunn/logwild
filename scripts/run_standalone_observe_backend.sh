#!/bin/zsh

setopt err_exit
setopt no_unset
setopt pipe_fail

declare -g CONFIG_DIR="${0:A:h}/cfg"
declare -g LOCAL_BIN_DIR="$HOME/.local/bin"
declare -g LOCAL_ZINC_DATA_PATH="${0:A:h}/data"
declare -g OPENOBSERVE_CONTAINER_NAME="logwild-openobserve"

function die() {
    local msg="$1"
    printf "%s\n" "${msg}" >&2
    exit 222
}

function stop_openobserve() {
    printf "stopping openobserve container [name=%s]\n" "${OPENOBSERVE_CONTAINER_NAME}"
    docker kill "${OPENOBSERVE_CONTAINER_NAME}" >/dev/null 2>&1 || true
    docker rm "${OPENOBSERVE_CONTAINER_NAME}" >/dev/null 2>&1 || true
    printf "stopped openobserve container [name=%s]\n" "${OPENOBSERVE_CONTAINER_NAME}"
}

function start_openobserve() {
    # create data directory for storing zinc data, if it doesn't exist
    mkdir -p "${LOCAL_ZINC_DATA_PATH}"
    printf "starting openobserve container\n"
    docker create -i -t \
        --name ${OPENOBSERVE_CONTAINER_NAME} \
        --mount type=bind,src=${LOCAL_ZINC_DATA_PATH},dst=/data \
        --mount type=bind,src=/etc/localtime,dst=/etc/localtime,ro \
        --tmpfs=/tmp:rw,noexec,nosuid,size=128m \
        -p 5080:5080 \
        --network observe \
        -e ZO_ROOT_USER_EMAIL="root@example.com" \
        -e ZO_ROOT_USER_PASSWORD="Complexpass#123" \
        public.ecr.aws/zinclabs/openobserve:latest
    printf "started openobserve container [name=%s]\n" "${OPENOBSERVE_CONTAINER_NAME}"
    # attach the container's in/out file descriptors
    docker start -ia ${OPENOBSERVE_CONTAINER_NAME}
}

function cleanup() {
    # do any necessary cleanup here to make sure program doesn't
    # leave something running after exit.
    printf "cleanup called\n"
    stop_openobserve
}

# call cleanup function on exit to remove anything
# we may have created
trap 'cleanup' EXIT
start_openobserve || die "uhoh, couldn't start openobserve"
