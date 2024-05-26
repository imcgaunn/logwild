#!/bin/zsh

setopt err_exit
setopt no_unset
setopt pipe_fail

declare -g OTEL_COLLECTOR_CONTAINER_NAME="${OTEL_COLLECTOR_CONTAINER_NAME:-logwild-otel-coll}"
declare -g OTEL_COLLECTOR_GRPC_PORT="${OTEL_COLLECTOR_GRPC_PORT:-4317}"

declare -g LOGWILD_CONTAINER_NAME="logwild-app"
declare -g LOGWILD_HTTP_PORT="8888"
declare -g LOGWILD_OTEL_SVC_NAME="logwild-imcg"
declare -g LOGWILD_DEBUG="true"
declare -g LOGWILD_OUT_FILE="/tmp/logwild.log"
declare -g LOGWILD_OTLP_ENDPOINT="http://${OTEL_COLLECTOR_CONTAINER_NAME}:${OTEL_COLLECTOR_GRPC_PORT}"
declare -g LOGWILD_IMAGE_REPOSITORY="mcgaunn.com/logwild"
declare -g LOGWILD_IMAGE_TAG="latest"

function die() {
    local msg="$1"
    printf "%s\n" "${msg}" >&2
    exit 222
}

function stop_app() {
    printf "stopping logwild container %s\n" "${LOGWILD_CONTAINER_NAME}"
    docker kill "${LOGWILD_CONTAINER_NAME}" >/dev/null 2>&1 || true
    docker rm "${LOGWILD_CONTAINER_NAME}" >/dev/null 2>&1 || true
    printf "stopped logwild container %s\n" "${LOGWILD_CONTAINER_NAME}"
}

function start_app() {
    local app_http_port="$1"
    printf "starting app on port %s\n" "${app_http_port}"
    # run the container with all of my wonderful settings
    # and detach from it, while keeping stdin open and allocating a tty
    docker run -d -i -t \
        --name="${LOGWILD_CONTAINER_NAME}" \
        -p ${app_http_port}:${LOGWILD_HTTP_PORT} \
        --network observe \
        --mount type=bind,src=/etc/localtime,dst=/etc/localtime,ro \
        --tmpfs=/tmp:rw,noexec,nosuid,size=128m \
        "${LOGWILD_IMAGE_REPOSITORY}:${LOGWILD_IMAGE_TAG}" \
        --debug \
        --port-metrics 8889 \
        --port ${LOGWILD_HTTP_PORT} \
        --otel-service-name ${LOGWILD_OTEL_SVC_NAME} \
        --log-rate 5000 \
        --log-size 48 \
        --log-burst-duration 3 \
        --log-out-file "${LOGWILD_OUT_FILE}" \
        run  # the actual command
    printf "started container [name=%s]\n" "${LOGWILD_CONTAINER_NAME}"
    # attach the container's in/out file descriptors
    docker attach ${LOGWILD_CONTAINER_NAME}
}

function cleanup() {
    printf "cleanup called\n"
    stop_app
}


# call cleanup on exit to remove anything we may have created
trap 'cleanup' EXIT
start_app "${LOGWILD_HTTP_PORT}" || die "failed to start logwild container"
