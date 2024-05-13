#!/bin/zsh

setopt err_exit
setopt no_unset
setopt pipe_fail

declare -g CONFIG_DIR="${0:A:h}/cfg"
declare -g OTEL_COLLECTOR_GRPC_PORT="${OTEL_COLLECTOR_GRPC_PORT:-4317}"
declare -g OTEL_COLLECTOR_CONTAINER_NAME="${OTEL_COLLECTOR_CONTAINER_NAME:-ianpod-otel-coll}"

function die() {
  local msg="$1"
  printf "%s\n" "${msg}" >&2
  exit 222
}

function stop_collector() {
  # stop collector program started earlier
  printf "stopping collector container %s\n" "${OTEL_COLLECTOR_CONTAINER_NAME}"
  docker kill "${OTEL_COLLECTOR_CONTAINER_NAME}" > /dev/null 2>&1 || true
  docker rm "${OTEL_COLLECTOR_CONTAINER_NAME}" > /dev/null 2>&1 || true
  printf "stopped collector container ${OTEL_COLLECTOR_CONTAINER_NAME}"
}

function start_collector() {
  local collector_port="$1"
  # start the collector program listening on the desired grpc port
  printf "starting collector on port %s\n" "${collector_port}" 
  # run the container with all of my wonderful settings
  # and detach from it, while keeping stdin open and allocating a tty
  docker run -d -i -t \
    --name=${OTEL_COLLECTOR_CONTAINER_NAME} \
    -p 44317:${OTEL_COLLECTOR_GRPC_PORT} \
    --network observe \
    --mount type=bind,src=./myapplog,dst=/var/log/syslog,ro \
    --mount type=bind,src=/etc/timezone,dst=/etc/timezone,ro \
    --mount type=bind,src=/etc/localtime,dst=/etc/localtime,ro \
    --mount type=bind,src=${CONFIG_DIR}/otelcol.yaml,dst=/etc/otelcol-contrib/config.yaml \
    --tmpfs=/tmp:rw,noexec,nosuid,size=128m \
    docker.io/otel/opentelemetry-collector-contrib:0.95.0 || die "could not create container"
  printf "started container [name=%s]\n" "${OTEL_COLLECTOR_CONTAINER_NAME}"
  # attach the container's in/out file descriptors 
  docker attach ${OTEL_COLLECTOR_CONTAINER_NAME}
}

function cleanup() {
  # do any necessary cleanup here to make sure program doesn't
  # leave something running after exit.
  printf "cleanup called\n"
  stop_collector
}

# call cleanup function on exit to remove anything
# we may have created
trap 'cleanup' EXIT
start_collector "${OTEL_COLLECTOR_GRPC_PORT}" || die "failed to start collector"
