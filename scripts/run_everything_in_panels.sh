#!/bin/zsh

setopt err_exit
setopt no_unset
setopt pipe_fail

SCRIPT_DIR="${0:A:h}"
SESS_NAME="logwild_demo"

function die() {
    echo "i am dead"
    exit 222
}

function cleanup() {
    for window_addr in "1.1" "1.2" "1.3"; do
        printf "sending ctrlc to window_addr: %s\n" "${window_addr}"
        tmux send-keys -t ${SESS_NAME}:${window_addr} C-c
    done
    printf "successfully cleaned up\n"
}

trap 'cleanup' EXIT

if tmux has-session -t "${SESS_NAME}" 2>/dev/null; then

    printf "killing existing session %s\n" "${SESS_NAME}"
    tmux kill-session -t "${SESS_NAME}"
fi

# create a new tmux session
tmux new-session -d -s ${SESS_NAME}
# run openobserve backend in first pane in first window
tmux send-keys -t ${SESS_NAME}:1.1 'just run-observe-backend' C-m
# split first pane horizontally
tmux split-window -t ${SESS_NAME}:1.1 -h
# start otelcol in new pane we just created
tmux send-keys -t ${SESS_NAME}:1.2 'just run-otel-collector' C-m
# split the new pane vertically
tmux split-window -t ${SESS_NAME}:1.2 -v
# run logwild server in the new vertical split
tmux send-keys -t ${SESS_NAME}:1.3 'just run-app-container' C-m

read "REPLY?program running in other window waiting"
